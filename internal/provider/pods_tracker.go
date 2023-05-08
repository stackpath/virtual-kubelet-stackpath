package provider

import (
	"context"
	"net/http"
	"time"

	"github.com/virtual-kubelet/virtual-kubelet/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1listers "k8s.io/client-go/listers/core/v1"
)

// Define the intervals for pod status updates and stale pod cleanup
var podStatusUpdateInterval = 5 * time.Second
var stalePodCleanupInterval = 5 * time.Minute

// PodsTracker manages the tracking of pod statuses and updates within a Kubernetes cluster
type PodsTracker struct {
	podLister      corev1listers.PodLister
	updateCallback func(*v1.Pod)
	handler        PodsTrackerHandler
}

type PodsTrackerHandler interface {
	GetPods(ctx context.Context) ([]*v1.Pod, error)
	GetPodStatus(ctx context.Context, ns, name string) (*v1.PodStatus, error)
	DeletePod(ctx context.Context, pod *v1.Pod) error
}

// BeginPodTracking initializes and manages background tracking for created pods
func (pt *PodsTracker) BeginPodTracking(ctx context.Context) {

	// Set up timers for periodic status updates and stale pods cleanup
	statusUpdatesTimer := time.NewTimer(podStatusUpdateInterval)
	cleanupTimer := time.NewTimer(stalePodCleanupInterval)
	defer statusUpdatesTimer.Stop()
	defer cleanupTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.G(ctx).WithError(ctx.Err()).Debug("Pod status update loop exiting")
			return
		case <-statusUpdatesTimer.C:
			pt.updatePods(ctx)
			statusUpdatesTimer.Reset(podStatusUpdateInterval)
		case <-cleanupTimer.C:
			pt.removeStalePods(ctx)
			cleanupTimer.Reset(stalePodCleanupInterval)
		}
	}
}

// updatePods synchronizes a list of pods in the indexer with their current status in the Kubernetes cluster.
// It iterates through all existing pods in the indexer, fetching their status from the provider via an API call.
// Once the status is retrieved, a callback function is invoked to handle any necessary updates.
func (pt *PodsTracker) updatePods(ctx context.Context) {
	k8sPods, err := pt.podLister.List(labels.Everything())
	if err != nil {
		log.G(ctx).WithError(err).Errorf("failed to retrieve pods list")
		return
	}
	for _, pod := range k8sPods {
		updatedPod := pod.DeepCopy()
		ok := pt.handlePodUpdates(ctx, updatedPod)
		if ok {
			pt.updateCallback(updatedPod)
		}
	}
}

// removeStalePods identifies and removes any pods in the PodsTracker that are no longer present in the Kubernetes cluster.
func (pt *PodsTracker) removeStalePods(ctx context.Context) {
	log.G(ctx).Debug("remove stale Pods")

	// Getting a list of pods indexed in k8s cluster
	clusterPods, err := pt.podLister.List(labels.Everything())
	if err != nil {
		log.G(ctx).WithError(err).Errorf("failed to retrieve pods list")
		return
	}
	// getting a list of pod identifiers of the pods running on the provider.
	activePods, err := pt.handler.GetPods(ctx)
	if err != nil {
		log.G(ctx).WithError(err).Errorf("failed to retrieve active workloads list")
		return
	}

	// Loop through all pods that are running on the provider
	for i := range activePods {
		if pod := getPodFromList(clusterPods, activePods[i].Namespace, activePods[i].Name); pod != nil {
			continue
		}

		log.G(ctx).Errorf("removing stale pod %v", activePods[i].Name)
		err := pt.handler.DeletePod(ctx, activePods[i])
		if err != nil {
			log.G(ctx).WithError(err).Errorf("failed to remove stale pod %v", activePods[i].Name)
		}
	}
}

// handlePodUpdates processes updates for a given pod in the PodsTracker, based on the current status of the pod within the Kubernetes cluster.
// The function returns a boolean indicating whether the update was successful (true) or not (false).
func (pt *PodsTracker) handlePodUpdates(ctx context.Context, pod *v1.Pod) bool {
	log.G(ctx).Debug("process Pod Updates")

	if pt.isPodStatusUpdateRequired(pod) {
		log.G(ctx).Infof("pod %s will skip pod status update", pod.Name)
		return false
	}

	newStatus, err := pt.handler.GetPodStatus(ctx, pod.Namespace, pod.Name)
	if err == nil && newStatus != nil {
		newStatus.DeepCopyInto(&pod.Status)
		return true
	}
	if err != nil {
		if apiError, ok := err.(*APIError); ok {
			if pod.Status.Phase == v1.PodRunning && apiError.statusCode == http.StatusNotFound {
				// Not found on the Edge side, probably was deleted by a user.
				// In that case, changing the pod's phase to 'Failed'
				// Set the pod to failed, this makes sure if the underlying container implementation is gone that a new pod will be created.
				pod.Status.Phase = v1.PodFailed
				pod.Status.Reason = "NotFoundOnProvider"
				pod.Status.Message = "the workload has been deleted from StackPath Edge Compute"
				now := metav1.NewTime(time.Now())
				for i := range pod.Status.ContainerStatuses {
					if pod.Status.ContainerStatuses[i].State.Running == nil {
						continue
					}

					pod.Status.ContainerStatuses[i].State.Terminated = &v1.ContainerStateTerminated{
						ExitCode:    137,
						Reason:      "NotFoundOnProvider",
						Message:     "the workload has been deleted from StackPath Edge Compute",
						FinishedAt:  now,
						StartedAt:   pod.Status.ContainerStatuses[i].State.Running.StartedAt,
						ContainerID: pod.Status.ContainerStatuses[i].ContainerID,
					}
					pod.Status.ContainerStatuses[i].State.Running = nil
				}
				return true
			}
		}
		log.G(ctx).WithError(err).Errorf("failed to retrieve pod %v status from provider", pod.Name)
	}
	return false
}

// isPodStatusUpdateRequired determines whether a given pod requires a status update within the PodsTracker.
// The function returns false if the pod has completed its execution (PodSucceeded), has failed (PodFailed),
// or is in the process of being terminated (DeletionTimestamp is set).
// Otherwise, it returns true, indicating that the pod's status should be updated.
func (pt *PodsTracker) isPodStatusUpdateRequired(pod *v1.Pod) bool {
	return pod.Status.Phase == v1.PodSucceeded || // Pod completed its execution
		pod.Status.Phase == v1.PodFailed ||
		pod.Status.Reason == "ProviderFailed" || // in case if provider failed to create/register the pod
		pod.DeletionTimestamp != nil // Terminating
}

func getPodFromList(list []*v1.Pod, ns, name string) *v1.Pod {
	for _, pod := range list {
		if pod.Namespace == ns && pod.Name == name {
			return pod
		}
	}

	return nil
}

package provider

import (
	"context"
	"time"

	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *StackpathProvider) getK8SPodStatusFrom(ctx context.Context, instance *workload_models.Workloadv1Instance) *v1.PodStatus {
	containerStatuses := make([]v1.ContainerStatus, 0, len(instance.ContainerStatuses))
	nameImageMap := make(map[string]string)

	for k, v := range instance.Containers {
		nameImageMap[k] = v.Image
	}

	isAllReady := true

	for i := range instance.ContainerStatuses {
		s := v1.ContainerStatus{
			Name:         instance.ContainerStatuses[i].Name,
			State:        getContainerState(instance.ContainerStatuses[i]),
			Ready:        instance.ContainerStatuses[i].Ready,
			RestartCount: instance.ContainerStatuses[i].RestartCount,
			Image:        nameImageMap[instance.ContainerStatuses[i].Name],
			ImageID:      "",
			ContainerID:  instance.ContainerStatuses[i].ContainerID,
		}
		containerStatuses = append(containerStatuses, s)
		if !s.Ready {
			isAllReady = false
		}
	}

	podIPs := make([]v1.PodIP, 0)
	if instance.IPAddress != "" {
		podIPs = append(podIPs, v1.PodIP{IP: instance.IPAddress})
	}
	if instance.IPV6Address != "" {
		podIPs = append(podIPs, v1.PodIP{IP: instance.IPV6Address})
	}

	ps := v1.PodStatus{
		Phase:             getPodPhaseFromInstancePhase(string(*instance.Phase)),
		Conditions:        getPodCondtions(instance, isAllReady),
		Message:           instance.Message,
		Reason:            instance.Reason,
		HostIP:            p.internalIP,
		PodIP:             instance.IPAddress,
		PodIPs:            podIPs,
		StartTime:         &metav1.Time{Time: time.Time(instance.StartedAt)},
		ContainerStatuses: containerStatuses,
		QOSClass:          v1.PodQOSBestEffort,
	}

	return &ps
}

func getPodCondtions(instance *workload_models.Workloadv1Instance, isAllReady bool) []v1.PodCondition {
	podPhase := getPodPhaseFromInstancePhase(string(*instance.Phase))

	switch podPhase {
	case v1.PodRunning, v1.PodSucceeded:

		readyConditionStatus := v1.ConditionFalse
		readyConditionTime := instance.CreatedAt
		if isAllReady {
			readyConditionStatus = v1.ConditionTrue
			readyConditionTime = instance.StartedAt
		}

		return []v1.PodCondition{
			{
				Type:               v1.PodReady,
				Status:             readyConditionStatus,
				LastTransitionTime: metav1.Time{Time: time.Time(readyConditionTime)},
			}, {
				Type:               v1.PodInitialized,
				Status:             v1.ConditionTrue,
				LastTransitionTime: metav1.Time{Time: time.Time(instance.CreatedAt)},
			}, {
				Type:               v1.PodScheduled,
				Status:             v1.ConditionTrue,
				LastTransitionTime: metav1.Time{Time: time.Time(instance.CreatedAt)},
			},
		}
	}

	return []v1.PodCondition{}
}

func getContainerState(s *workload_models.V1ContainerStatus) v1.ContainerState {
	if s.Running != nil {
		return v1.ContainerState{
			Running: &v1.ContainerStateRunning{
				StartedAt: metav1.Time{Time: time.Time(s.Running.StartedAt)},
			},
		}
	} else if s.Waiting != nil {
		return v1.ContainerState{
			Waiting: &v1.ContainerStateWaiting{
				Reason:  s.Waiting.Reason,
				Message: s.Waiting.Message,
			},
		}
	} else if s.Terminated != nil {
		return v1.ContainerState{
			Terminated: &v1.ContainerStateTerminated{
				ExitCode:   s.Terminated.ExitCode,
				FinishedAt: metav1.Time{Time: time.Time(s.Terminated.FinishedAt)},
				Reason:     s.Terminated.Reason,
				Message:    s.Terminated.Message,
				StartedAt:  metav1.Time{Time: time.Time(s.Terminated.StartedAt)},
			},
		}
	}

	// The container state is waiting by default
	return v1.ContainerState{Waiting: &v1.ContainerStateWaiting{}}
}

func getPodPhaseFromInstancePhase(state string) v1.PodPhase {
	switch state {
	case "RUNNING":
		return v1.PodRunning
	case "COMPLETED":
		return v1.PodSucceeded
	case "STOPPED":
		return v1.PodSucceeded
	case "FAILED":
		return v1.PodFailed
	case "SCHEDULING":
		return v1.PodPending
	case "STARTING":
		return v1.PodPending
	case "DELETING":
		return v1.PodPending
	}

	return v1.PodUnknown
}

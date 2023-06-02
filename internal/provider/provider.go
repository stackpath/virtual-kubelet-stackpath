// Package provider implements the stackpath virtual kubelet provider
package provider

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_client"
	"github.com/stackpath/vk-stackpath-provider/internal/config"
	"github.com/virtual-kubelet/virtual-kubelet/log"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	stats "github.com/virtual-kubelet/virtual-kubelet/node/api/statsv1alpha1"
	"github.com/virtual-kubelet/virtual-kubelet/node/nodeutil"
	"golang.org/x/crypto/ssh"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
)

const (
	// targetName is a constant string representing the name shared by all targets
	// in the workloads. In this case, it is set to the city code.
	targetName = "city-code"

	// targetOrdinal is a constant string representing the ordinal number for
	// the target instance. It is set to "0" to indicate the first instance in
	// a given workload.
	targetOrdinal = "0"

	nodeNameLabelKey = "vk-node-name"

	podNameLabelKey = "vk-pod-name"

	podNamespaceLabelKey = "vk-pod-namespace"

	containerConsoleHost = "container-console.edgeengine.io"

	containerConsolePort = 9600
)

// StackpathProvider is a struct that implements the virtual-kubelet provider interface
type StackpathProvider struct {
	secretLister    corev1listers.SecretLister
	configMapLister corev1listers.ConfigMapLister
	podLister       corev1listers.PodLister

	stackpathClient    *workload_client.EdgeCompute
	apiConfig          *config.Config
	cpu                string
	memory             string
	pods               string
	storage            string
	operatingSystem    string
	nodeName           string
	startTime          time.Time
	internalIP         string
	daemonEndpointPort int32

	podsTracker *PodsTracker

	logger log.Logger
}

// NewStackpathProvider creates a stackpath virtual kubelet provider
func NewStackpathProvider(ctx context.Context, stackpathClient *workload_client.EdgeCompute, apiConfig *config.Config, providerConfig nodeutil.ProviderConfig, internalIP string, daemonEndpointPort int32) (*StackpathProvider, error) {
	log.G(ctx).Debug("creating a new StackPath provider")
	var provider StackpathProvider
	provider.configMapLister = providerConfig.ConfigMaps
	provider.secretLister = providerConfig.Secrets
	provider.podLister = providerConfig.Pods
	provider.stackpathClient = stackpathClient
	provider.nodeName = providerConfig.Node.Name
	provider.startTime = time.Now()
	provider.apiConfig = apiConfig
	provider.internalIP = internalIP
	provider.daemonEndpointPort = daemonEndpointPort
	provider.setNodeCapacity()
	provider.logger = log.G(ctx)

	return &provider, nil
}

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *StackpathProvider) CreatePod(ctx context.Context, pod *v1.Pod) error {

	w, err := p.getWorkloadFrom(pod)
	if err != nil {
		return err
	}

	err = p.createWorkload(ctx, w)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *StackpathProvider) UpdatePod(ctx context.Context, pod *v1.Pod) error {
	// NOP. Not Implemented in this version
	return nil
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *StackpathProvider) DeletePod(ctx context.Context, pod *v1.Pod) error {
	log.G(ctx).Debugf("deleting the pod %s", pod.Name)

	err := p.deleteWorkload(ctx, pod.Namespace, pod.Name)
	if err != nil {
		return err
	}
	return nil
}

// GetPod retrieves a pod by its name and namespace running on StackPath Edge Compute
// and updates the cluster's pod status to match the provider's.
// If a pod with the specified name is not found, it returns nil and an error.
func (p *StackpathProvider) GetPod(ctx context.Context, namespace, name string) (*v1.Pod, error) {
	p.logger.Debugf("getting the pod (namespace: %s, name: %s)", namespace, name)

	_, err := p.getWorkload(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	instance, err := p.getWorkloadInstance(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	pod, err := p.podLister.Pods(namespace).Get(name)
	if err != nil {
		return nil, err
	}

	updatedPod := pod.DeepCopy()

	podStatus := p.getK8SPodStatusFrom(ctx, instance)
	updatedPod.Status = *podStatus

	return updatedPod, nil
}

// GetPodStatus retrieves the status of a pod by name and namespace from the provider.
func (p *StackpathProvider) GetPodStatus(ctx context.Context, namespace, name string) (*v1.PodStatus, error) {
	log.G(ctx).Debugf("getting the pod's status (namespace: %s, name: %s)", namespace, name)

	instance, err := p.getWorkloadInstance(ctx, namespace, name)

	if err != nil {
		return nil, err
	}
	return p.getK8SPodStatusFrom(ctx, instance), nil
}

// GetPods retrieves a list of all pods running on the provider.
func (p *StackpathProvider) GetPods(ctx context.Context) ([]*v1.Pod, error) {
	log.G(ctx).Info("getting a list of workloads")

	workloads, err := p.getWorkloads(ctx)
	if err != nil {
		return nil, err
	}

	if workloads == nil {
		log.G(ctx).Info("no workloads found")
		return nil, nil
	}

	log.G(ctx).Infof("%s workload(s) found", strconv.Itoa(len(workloads)))

	pods := make([]*v1.Pod, 0, len(workloads))

	for _, workload := range workloads {
		// Only use workloads that were created by this provider
		if workload.Metadata.Labels[nodeNameLabelKey] != p.nodeName {
			continue
		}

		podNamespace := workload.Metadata.Labels[podNamespaceLabelKey]
		podName := workload.Metadata.Labels[podNameLabelKey]
		instance, err := p.getWorkloadInstance(ctx, podNamespace, podName)
		if err != nil {
			log.G(ctx).WithFields(log.Fields{
				"id": workload.ID,
			}).WithError(err).Errorf("couldn't get an instance of the workload. The workload will be deleted.")

			// In case if the workload exists but doesn't have any instances associated with it,
			// the workload has to be deleted. The following lines needed to add the info about the pod
			// that then will be used for removing any staled workload like that.
			pods = append(pods, &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: podNamespace,
				},
			})
			continue
		}

		pod, err := p.getPodFromListerByInstance(ctx, instance, &podNamespace, &podName)
		if err != nil {
			log.G(ctx).WithFields(log.Fields{
				"id":            workload.ID,
				"instance-name": instance.Name,
			}).WithError(err).Errorf("couldn't translate the instance to a pod")
			continue
		}

		pods = append(pods, pod)
	}

	return pods, nil
}

// GetContainerLogs returns the logs of a pod by name that is running as a StackPath workload
func (p *StackpathProvider) GetContainerLogs(ctx context.Context, namespace, podName, containerName string, opts api.ContainerLogOpts) (io.ReadCloser, error) {
	// NOP. Not Implemented in this version
	return nil, nil
}

// RunInContainer executes a command in a container in the pod, copying data
// between in/out/err and the container's stdin/stdout/stderr.
func (p *StackpathProvider) RunInContainer(ctx context.Context, namespace, name, container string, cmd []string, attach api.AttachIO) error {

	conf := &ssh.ClientConfig{
		User:            p.getSSHUsername(namespace, name, container),
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		Auth: []ssh.AuthMethod{
			ssh.Password(p.apiConfig.ClientSecret),
		},
	}

	var conn *ssh.Client

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", containerConsoleHost, containerConsolePort), conf)
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if attach.TTY() {
		// Set up terminal modes
		modes := ssh.TerminalModes{}
		err = session.RequestPty("Xterm", 120, 60, modes)
		if err != nil {
			return err
		}
	}

	sessionStdinPipe, err := session.StdinPipe()
	if err != nil {
		return err
	}
	defer sessionStdinPipe.Close()

	sessionStdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	sessionStderrPipe, err := session.StderrPipe()
	if err != nil {
		return err
	}

	// goroutine that is responsible for listening 'resize' channel and update the session window measurements accordingly
	go func() {
		for {
			select {
			case size := <-attach.Resize():
				err := session.WindowChange(int(size.Height), int(size.Width))
				if err != nil {
					p.logger.Info(err)
				}
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	// Channel that is used for sending errors happened during stdout pipe read.
	c := make(chan error, 1)

	aout := attach.Stdout()
	if aout != nil {
		defer aout.Close()
	}
	go func() {
		_, err := io.Copy(aout, sessionStdoutPipe)
		if err != nil {
			// io.EOF or an error
			c <- err
			return
		}
	}()
	go func() { io.Copy(aout, sessionStderrPipe) }()

	ain := attach.Stdin()
	if ain != nil {
		go func() { io.Copy(sessionStdinPipe, ain) }()
	}

	// sending the command
	go func() {
		if err := session.Run(strings.Join(cmd, " ") + "\n"); err != nil {
			c <- err
			return
		}
	}()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case e := <-c:
			p.logger.Debug(e)
			break loop
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}

	return ctx.Err()
}

func (p *StackpathProvider) AttachToContainer(ctx context.Context, namespace, podName, containerName string, attach api.AttachIO) error {
	// NOP. Not Implemented in this version
	return nil
}

// GetMetricsResource gets the metrics for the node, including running pods
func (p *StackpathProvider) GetMetricsResource(context.Context) ([]*dto.MetricFamily, error) {
	// NOP. Not Implemented in this version
	return nil, nil
}

// NotifyPods instructs the notifier to call the passed in function when the pod status changes.
// The provided pointer to a Pod is guaranteed to be used in a read-only fashion.
func (p *StackpathProvider) NotifyPods(ctx context.Context, notifierCallback func(*v1.Pod)) {
	p.podsTracker = &PodsTracker{
		podLister:      p.podLister,
		updateCallback: notifierCallback,
		handler:        p,
	}

	go p.podsTracker.BeginPodTracking(ctx)
}

// GetStatsSummary gets the stats for the node, including running pods
func (p *StackpathProvider) GetStatsSummary(ctx context.Context) (*stats.Summary, error) {
	// NOP. Not implemented in this version
	return nil, nil
}

func (p *StackpathProvider) getSSHUsername(namespace string, name string, containerName string) string {
	return p.apiConfig.StackID + "." + p.getInstanceName(namespace, name) + "." + containerName + "." + p.apiConfig.ClientID
}

package provider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_client/instance"
	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_client/instance_logs"
	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_client/workload"
	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
	"github.com/virtual-kubelet/virtual-kubelet/errdefs"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *StackpathProvider) getWorkloads(ctx context.Context) ([]*workload_models.V1Workload, error) {
	pageSize := "99999"

	params := &workload.GetWorkloadsParams{
		StackID:          p.apiConfig.StackID,
		Context:          ctx,
		PageRequestFirst: &pageSize,
	}

	response, err := p.stackpathClient.Workload.GetWorkloads(params, nil)
	if err != nil {
		return nil, NewStackPathError(err)
	}

	return response.Payload.Results, nil
}

func (p *StackpathProvider) getWorkload(ctx context.Context, namespace string, name string) (*workload_models.V1Workload, error) {
	getWorkloadParams := workload.GetWorkloadParams{
		Context:    ctx,
		StackID:    p.apiConfig.StackID,
		WorkloadID: p.getWorkloadSlug(namespace, name),
	}

	workloadResult, err := p.stackpathClient.Workload.GetWorkload(&getWorkloadParams, nil)
	if err != nil {
		return nil, NewStackPathError(err)
	}
	return workloadResult.Payload.Workload, nil
}

func (p *StackpathProvider) getWorkloadInstance(ctx context.Context, namespace string, name string) (*workload_models.Workloadv1Instance, error) {
	getInstanceParams := instance.GetWorkloadInstanceParams{
		Context:      ctx,
		StackID:      p.apiConfig.StackID,
		WorkloadID:   p.getWorkloadSlug(namespace, name),
		InstanceName: p.getInstanceName(namespace, name),
	}
	instanceResult, err := p.stackpathClient.Instance.GetWorkloadInstance(&getInstanceParams, nil)
	if err != nil {
		return nil, NewStackPathError(err)
	}
	return instanceResult.Payload.Instance, nil
}

func (p *StackpathProvider) getPodFromListerByInstance(ctx context.Context, instance *workload_models.Workloadv1Instance, namespace, name *string) (*v1.Pod, error) {
	pod, err := p.podLister.Pods(*namespace).Get(*name)
	// in case pod got deleted, we want to continue the workflow to kick off remove stale pods from the provider
	if errdefs.IsNotFound(err) || pod == nil {
		return &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      *name,
				Namespace: *namespace,
			},
		}, nil
	}
	if err != nil {
		return nil, err
	}

	updatedPod := pod.DeepCopy()

	podState := p.getK8SPodStatusFrom(ctx, instance)

	updatedPod.Status = *podState

	return updatedPod, nil
}

func (p *StackpathProvider) createWorkload(ctx context.Context, w *workload_models.V1Workload) error {
	params := workload.CreateWorkloadParams{
		Body:    &workload_models.V1CreateWorkloadRequest{Workload: w},
		StackID: p.apiConfig.StackID,
		Context: ctx,
	}

	_, err := p.stackpathClient.Workload.CreateWorkload(&params, nil)
	if err != nil {
		return NewStackPathError(err)
	}
	return nil
}

func (p *StackpathProvider) deleteWorkload(ctx context.Context, podNamespace, podName string) error {
	params := workload.DeleteWorkloadParams{
		StackID:    p.apiConfig.StackID,
		WorkloadID: p.getWorkloadSlug(podNamespace, podName),
		Context:    ctx,
	}

	_, err := p.stackpathClient.Workload.DeleteWorkload(&params, nil)
	if err != nil {
		return NewStackPathError(err)
	}

	return nil
}

func (p *StackpathProvider) getInstanceLogsReader(ctx context.Context, podNamespace, podName, containerName string, opts api.ContainerLogOpts) io.ReadCloser {
	params := instance_logs.GetLogsParams{
		Context:       ctx,
		StackID:       p.apiConfig.StackID,
		WorkloadID:    p.getWorkloadSlug(podNamespace, podName),
		ContainerName: &containerName,
		InstanceName:  p.getInstanceName(podNamespace, podName),
		Follow:        &opts.Follow,
		Previous:      &opts.Previous,
		Timestamps:    &opts.Timestamps,
	}

	sinceTime := strfmt.DateTime(opts.SinceTime)
	params.SinceTime = &sinceTime

	if opts.SinceSeconds > 0 {
		sinceSeconds := strconv.Itoa(opts.SinceSeconds)
		params.SinceSeconds = &sinceSeconds
	}
	if opts.LimitBytes > 0 {
		limitBytes := strconv.Itoa(opts.LimitBytes)
		params.LimitBytes = &limitBytes
	}
	if opts.Tail > 0 {
		tailLines := strconv.Itoa(opts.Tail)
		params.TailLines = &tailLines
	}

	reader, writer := io.Pipe()

	go func() {
		_, err := p.stackpathClient.InstanceLogs.GetLogs(&params, nil, writer)
		if err != nil && err != io.EOF {
			if err == io.ErrUnexpectedEOF {
				err = errors.New("the container logs retrieval process has been interrupted due to 60 seconds of inactivity")
			}
			writer.Write([]byte(err.Error()))
		}
		writer.Close()
	}()
	return reader
}

// getInstanceName returns the name of the first instance running in a workload.
//
// The instance name is formatted as follows:
// <workload-slug>-<target-name>-<deployment-city-code>-<ordinal>
func (p *StackpathProvider) getInstanceName(namespace, name string) string {
	return strings.ToLower(
		fmt.Sprintf("%s-%s-%s-%s", p.getWorkloadSlug(namespace, name), targetName, p.apiConfig.CityCode, targetOrdinal),
	)
}

// getWorkloadSlug returns the name of the workload associated with the pod.
//
// The workload name is formatted as follows:
// <pod-namespace>-<pod-name>
func (p *StackpathProvider) getWorkloadSlug(namespace, name string) string {
	return strings.ToLower(
		fmt.Sprintf("%s-%s", namespace, name),
	)
}

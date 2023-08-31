package provider

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetK8SPodStatusFrom(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	currentTime := time.Now()

	var testCases = []struct {
		description             string
		instance                *workload_models.Workloadv1Instance
		expectedStatus          v1.PodStatus
		expectedContainerStatus v1.ContainerState
	}{
		{
			description:    "successfully gets pod status from an instance with 'starting' phase",
			expectedStatus: v1.PodStatus{Phase: v1.PodPending},
			instance: createTestInstance(
				"test",
				workload_models.Workloadv1InstanceInstancePhaseSTARTING.Pointer(),
				&workload_models.V1ContainerStatus{
					Waiting: &workload_models.ContainerStatusWaiting{
						Reason:  "waiting",
						Message: "waiting message"},
				},
			),
			expectedContainerStatus: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "waiting", Message: "waiting message"}},
		},
		{
			description:    "successfully gets pod status from an instance with 'running' phase",
			expectedStatus: v1.PodStatus{Phase: v1.PodRunning},
			instance: createTestInstance(
				"test",
				workload_models.Workloadv1InstanceInstancePhaseRUNNING.Pointer(),
				&workload_models.V1ContainerStatus{
					Running: &workload_models.ContainerStatusRunning{
						StartedAt: strfmt.DateTime(currentTime),
					},
				},
			),
			expectedContainerStatus: v1.ContainerState{
				Running: &v1.ContainerStateRunning{StartedAt: metav1.Time{Time: time.Time(currentTime)}},
			},
		},
		{
			description:    "successfully gets pod status from an instance with 'failed' phase",
			expectedStatus: v1.PodStatus{Phase: v1.PodFailed},
			instance: createTestInstance(
				"test",
				workload_models.Workloadv1InstanceInstancePhaseFAILED.Pointer(),
				&workload_models.V1ContainerStatus{
					Terminated: &workload_models.ContainerStatusTerminated{
						FinishedAt: strfmt.DateTime(currentTime.AddDate(0, 0, -1)),
					},
				},
			),
			expectedContainerStatus: v1.ContainerState{
				Terminated: &v1.ContainerStateTerminated{FinishedAt: metav1.Time{Time: currentTime.AddDate(0, 0, -1)}},
			},
		},
		{
			description:    "successfully gets pod status from an instance with 'completed' phase",
			expectedStatus: v1.PodStatus{Phase: v1.PodSucceeded},
			instance: createTestInstance(
				"test",
				workload_models.Workloadv1InstanceInstancePhaseCOMPLETED.Pointer(),
				&workload_models.V1ContainerStatus{
					Terminated: &workload_models.ContainerStatusTerminated{
						FinishedAt: strfmt.DateTime(currentTime.AddDate(0, 0, -1)),
					},
				},
			),
			expectedContainerStatus: v1.ContainerState{
				Terminated: &v1.ContainerStateTerminated{FinishedAt: metav1.Time{Time: currentTime.AddDate(0, 0, -1)}},
			},
		},
		{
			description:    "successfully gets pod status from an instance with 'scheduled' phase",
			expectedStatus: v1.PodStatus{Phase: v1.PodPending},
			instance: createTestInstance(
				"test",
				workload_models.Workloadv1InstanceInstancePhaseSCHEDULING.Pointer(),
				&workload_models.V1ContainerStatus{
					Waiting: &workload_models.ContainerStatusWaiting{
						Reason:  "waiting",
						Message: "waiting message"},
				},
			),
			expectedContainerStatus: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "waiting", Message: "waiting message"}},
		},
		{
			description:    "successfully gets pod status from an instance with 'stopped' phase",
			expectedStatus: v1.PodStatus{Phase: v1.PodSucceeded},
			instance: createTestInstance(
				"test",
				workload_models.Workloadv1InstanceInstancePhaseSTOPPED.Pointer(),
				&workload_models.V1ContainerStatus{
					Terminated: &workload_models.ContainerStatusTerminated{
						FinishedAt: strfmt.DateTime(currentTime.AddDate(0, 0, -1)),
					},
				},
			),
			expectedContainerStatus: v1.ContainerState{
				Terminated: &v1.ContainerStateTerminated{FinishedAt: metav1.Time{Time: currentTime.AddDate(0, 0, -1)}},
			},
		},
		{
			description:    "successfully gets pod status from an instance with 'scheduling' phase and the container state not set yet",
			expectedStatus: v1.PodStatus{Phase: v1.PodPending},
			instance: createTestInstance(
				"test",
				workload_models.Workloadv1InstanceInstancePhaseSCHEDULING.Pointer(),
				&workload_models.V1ContainerStatus{},
			),
			expectedContainerStatus: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{}},
		},
		{
			description:    "successfully gets pod status from an instance that has a phase that doesn't supported by k8s",
			expectedStatus: v1.PodStatus{Phase: v1.PodUnknown},
			instance: createTestInstance(
				"test",
				workload_models.NewWorkloadv1InstanceInstancePhase("wrong-phase"),
				&workload_models.V1ContainerStatus{},
			),
			expectedContainerStatus: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{}},
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			status := provider.getK8SPodStatusFrom(ctx, c.instance)
			assert.Equal(t, c.expectedStatus.Phase, status.Phase)
			assert.Equal(t, c.expectedContainerStatus, status.ContainerStatuses[0].State)
		})
	}
}

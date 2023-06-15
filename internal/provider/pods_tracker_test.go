package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_client"
	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_client/instance"
	workloads "github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_client/workloads"
	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_models"
	mocks "github.com/stackpath/vk-stackpath-provider/internal/mocks"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRemoveStalePods(t *testing.T) {
	podName := fmt.Sprintf("test-pod-%s", uuid.New().String())
	stalePodName := fmt.Sprintf("test-pod-%s", uuid.New().String())
	podNamespace := fmt.Sprintf("test-ns-%s", uuid.New().String())
	nodeName := "test-node"

	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := context.Background()

	k8sPods := []*v1.Pod{createTestPod(podName, podNamespace)}

	wsc := mocks.NewWorkloadsClientService(mockController)
	isc := mocks.NewInstanceClientService(mockController)
	stackPathClientMock := workload_client.EdgeCompute{Workloads: wsc, Instance: isc}

	activePodsLister := mocks.NewMockPodLister(mockController)
	k8sPodsLister := mocks.NewMockPodLister(mockController)
	mockPodsNamespaceLister := mocks.NewMockPodNamespaceLister(mockController)

	provider, err := createTestProvider(ctx, mocks.NewMockConfigMapLister(mockController), mocks.NewMockSecretLister(mockController), activePodsLister, &stackPathClientMock)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}
	provider.nodeName = nodeName

	podsTracker := &PodsTracker{
		podLister: k8sPodsLister,
		updateCallback: func(updatedPod *v1.Pod) {
		},
		handler: provider,
	}

	testCases := []struct {
		description     string
		initMockedCalls func()
	}{
		{
			description: "successfully removes a stale pod (workload)",
			initMockedCalls: func() {
				// Mocks getting a list of pods indexed in k8s cluster
				k8sPodsLister.EXPECT().List(gomock.Any()).Return(k8sPods, nil).Times(1)

				// Mocks GetWorkloads() that returns two workloads, one is active,
				// second one is considered to be stale (no info in k8s cluster about it)
				wsc.EXPECT().GetWorkloads(gomock.Any(), gomock.Any()).Return(&workloads.GetWorkloadsOK{
					Payload: &workload_models.V1GetWorkloadsResponse{
						Results: []*workload_models.V1Workload{
							{
								Name: podName,
								Metadata: &workload_models.V1Metadata{
									Labels: workload_models.V1StringMapEntry{
										nodeNameLabelKey:     nodeName,
										podNamespaceLabelKey: podNamespace,
										podNameLabelKey:      podName,
									},
								},
							},
							{
								Name: stalePodName,
								Metadata: &workload_models.V1Metadata{
									Labels: workload_models.V1StringMapEntry{
										nodeNameLabelKey:     nodeName,
										podNamespaceLabelKey: podNamespace,
										podNameLabelKey:      stalePodName,
									},
								},
							},
							{
								Name: "MustBeIgnored",
								Metadata: &workload_models.V1Metadata{
									Labels: workload_models.V1StringMapEntry{
										nodeNameLabelKey:     "NotVirtualKubeletNode",
										podNamespaceLabelKey: podNamespace,
										podNameLabelKey:      stalePodName,
									},
								},
							},
						},
					},
				}, nil).Times(1)

				// Mocks GetWorkloadInstance() that returns an info about the instance (not staled)
				isc.EXPECT().GetWorkloadInstance(gomock.Any(), gomock.Any()).Return(&instance.GetWorkloadInstanceOK{
					Payload: &workload_models.V1GetWorkloadInstanceResponse{
						Instance: &workload_models.Workloadv1Instance{
							Name:  podName,
							Phase: workload_models.Workloadv1InstanceInstancePhaseRUNNING.Pointer(),
						},
					},
				}, nil).Times(1)

				// Mocks GetWorkloadInstance() that returns an info about the instance (staled)
				isc.EXPECT().GetWorkloadInstance(gomock.Any(), gomock.Any()).Return(&instance.GetWorkloadInstanceOK{
					Payload: &workload_models.V1GetWorkloadInstanceResponse{
						Instance: &workload_models.Workloadv1Instance{
							Name:  stalePodName,
							Phase: workload_models.Workloadv1InstanceInstancePhaseRUNNING.Pointer(),
						},
					},
				}, nil).Times(1)

				// Next two calls mock pod lister calls for getting info from k8s that happens during provider.getWorkloadInstance() call
				activePodsLister.EXPECT().Pods(podNamespace).Return(mockPodsNamespaceLister).Times(2)
				mockPodsNamespaceLister.EXPECT().Get(podName).Return(createTestPod(podName, podNamespace), nil).Times(1)
				mockPodsNamespaceLister.EXPECT().Get(stalePodName).Return(createTestPod(stalePodName, podNamespace), nil).Times(1)

				// Mocks workload deletion that deletes the staled pod
				wsc.EXPECT().DeleteWorkload(
					&workloads.DeleteWorkloadParams{
						StackID:    provider.apiConfig.StackID,
						WorkloadID: provider.getWorkloadSlug(podNamespace, stalePodName),
						Context:    ctx,
					}, gomock.Any()).Return(nil, nil).Times(1)
			},
		}, {
			description: "successfully removes a stale pod (workload) even if there was an API error during getting workload's instance information",
			initMockedCalls: func() {
				// Mocks getting a list of pods indexed in k8s cluster
				k8sPodsLister.EXPECT().List(gomock.Any()).Return(k8sPods, nil).Times(1)

				// Mocks GetWorkloads() that returns two workloads, one is active,
				// second one is considered to be stale (no info in k8s cluster about it)
				wsc.EXPECT().GetWorkloads(gomock.Any(), gomock.Any()).Return(&workloads.GetWorkloadsOK{
					Payload: &workload_models.V1GetWorkloadsResponse{
						Results: []*workload_models.V1Workload{
							{
								Name: podName,
								Metadata: &workload_models.V1Metadata{
									Labels: workload_models.V1StringMapEntry{
										nodeNameLabelKey:     nodeName,
										podNamespaceLabelKey: podNamespace,
										podNameLabelKey:      podName,
									},
								},
							},
							{
								Name: stalePodName,
								Metadata: &workload_models.V1Metadata{
									Labels: workload_models.V1StringMapEntry{
										nodeNameLabelKey:     nodeName,
										podNamespaceLabelKey: podNamespace,
										podNameLabelKey:      stalePodName,
									},
								},
							},
						},
					},
				}, nil).Times(1)

				// Mocks GetWorkloadInstance() that returns an info about the instance (not staled)
				isc.EXPECT().GetWorkloadInstance(gomock.Any(), gomock.Any()).Return(&instance.GetWorkloadInstanceOK{
					Payload: &workload_models.V1GetWorkloadInstanceResponse{
						Instance: &workload_models.Workloadv1Instance{
							Name:  podName,
							Phase: workload_models.Workloadv1InstanceInstancePhaseRUNNING.Pointer(),
						},
					},
				}, nil).Times(1)

				// Mocking a case where it fails to get the Workload's instance. In this case it is expected the DeleteWorkload to be called anyways
				isc.EXPECT().GetWorkloadInstance(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("a 404 error was returned from StackPath: \"record not found\""))

				// Next three calls mock pod lister calls for getting info from k8s that happens during provider.getWorkloadInstance() call
				activePodsLister.EXPECT().Pods(podNamespace).Return(mockPodsNamespaceLister).Times(1)
				mockPodsNamespaceLister.EXPECT().Get(podName).Return(createTestPod(podName, podNamespace), nil).Times(1)

				// Mocks workload deletion that deletes staled pod
				wsc.EXPECT().DeleteWorkload(
					&workloads.DeleteWorkloadParams{
						StackID:    provider.apiConfig.StackID,
						WorkloadID: provider.getWorkloadSlug(podNamespace, stalePodName),
						Context:    ctx,
					}, gomock.Any()).Return(nil, nil).Times(1)
			},
		}, {
			description: "fail to remove stale pod (workload) due to an error happened on getting a list of pods running in a cluster",
			initMockedCalls: func() {
				// Mocks getting a list of pods indexed in k8s cluster
				k8sPodsLister.EXPECT().List(gomock.Any()).Return(nil, fmt.Errorf("Error")).Times(1)
			},
		}, {
			description: "fail to remove stale pod (workload) due to an error happened on getting a list of workloads",
			initMockedCalls: func() {
				// Mocks getting a list of pods indexed in k8s cluster
				k8sPodsLister.EXPECT().List(gomock.Any()).Return(k8sPods, nil).Times(1)

				// Mocks GetWorkloads() that returns two workloads, one is active,
				// second one is considered to be stale (no info in k8s cluster about it)
				wsc.EXPECT().GetWorkloads(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error")).Times(1)
			},
		}, {
			description: "No stale pods were removed because there's no any on the provider side.",
			initMockedCalls: func() {
				// Mocks getting a list of pods indexed in k8s cluster
				k8sPodsLister.EXPECT().List(gomock.Any()).Return(k8sPods, nil).Times(1)

				// Mocks GetWorkloads() that returns zero workloads
				wsc.EXPECT().GetWorkloads(gomock.Any(), gomock.Any()).Return(
					&workloads.GetWorkloadsOK{
						Payload: &workload_models.V1GetWorkloadsResponse{
							Results: []*workload_models.V1Workload{},
						},
					}, nil).Times(1)
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			c.initMockedCalls()
			provider.podsTracker = podsTracker
			podsTracker.removeStalePods(context.Background())
			assert.Len(t, k8sPods, 1)
		})
	}
}

func TestHandlePodUpdates(t *testing.T) {
	podName := fmt.Sprintf("test-pod-%s", uuid.New().String())
	podNamespace := fmt.Sprintf("test-ns-%s", uuid.New().String())
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := context.Background()

	isc := mocks.NewInstanceClientService(mockController)
	stackPathClientMock := workload_client.EdgeCompute{Instance: isc}

	provider, err := createTestProvider(ctx, mocks.NewMockConfigMapLister(mockController), mocks.NewMockSecretLister(mockController), mocks.NewMockPodLister(mockController), &stackPathClientMock)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	podLister := mocks.NewMockPodLister(mockController)

	podsTracker := &PodsTracker{
		podLister:      podLister,
		updateCallback: func(p *v1.Pod) {},
		handler:        provider,
	}

	pod := createTestPod(podName, podNamespace)

	testCases := []struct {
		description               string
		initialPodPhase           v1.PodPhase
		expectedPodPhase          v1.PodPhase
		workloadInstancePhase     workload_models.Workloadv1InstanceInstancePhase
		isPodStatusUpdateRequired bool
	}{
		{
			description:               "Successfully updates pod phase to running",
			initialPodPhase:           v1.PodPending,
			expectedPodPhase:          v1.PodRunning,
			workloadInstancePhase:     workload_models.Workloadv1InstanceInstancePhaseRUNNING,
			isPodStatusUpdateRequired: true,
		},

		{
			description:               "Successfully updates pod phase to Failed",
			initialPodPhase:           v1.PodRunning,
			expectedPodPhase:          v1.PodFailed,
			workloadInstancePhase:     workload_models.Workloadv1InstanceInstancePhaseFAILED,
			isPodStatusUpdateRequired: true,
		},
		{
			description:               "Successfully updates pod phase to Succeeded",
			initialPodPhase:           v1.PodRunning,
			expectedPodPhase:          v1.PodSucceeded,
			workloadInstancePhase:     workload_models.Workloadv1InstanceInstancePhaseCOMPLETED,
			isPodStatusUpdateRequired: true,
		},
		{
			description:               "Ignores the status update, the pod is failed already",
			initialPodPhase:           v1.PodFailed,
			expectedPodPhase:          v1.PodFailed,
			workloadInstancePhase:     workload_models.Workloadv1InstanceInstancePhaseRUNNING,
			isPodStatusUpdateRequired: false,
		},
		{
			description:               "Ignores the status update, the pod is succeeded already",
			initialPodPhase:           v1.PodSucceeded,
			expectedPodPhase:          v1.PodSucceeded,
			workloadInstancePhase:     workload_models.Workloadv1InstanceInstancePhaseRUNNING,
			isPodStatusUpdateRequired: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			pod.Status.Phase = c.initialPodPhase

			i := createTestInstance(
				provider.getInstanceName(podNamespace, podName),
				workload_models.NewWorkloadv1InstanceInstancePhase(c.workloadInstancePhase),
				&workload_models.V1ContainerStatus{
					Waiting: &workload_models.ContainerStatusWaiting{
						Reason:  "waiting",
						Message: "waiting message"},
				},
			)

			params := instance.GetWorkloadInstanceParams{
				Context:      ctx,
				StackID:      provider.apiConfig.StackID,
				WorkloadID:   provider.getWorkloadSlug(podNamespace, podName),
				InstanceName: provider.getInstanceName(podNamespace, podName),
			}

			if c.isPodStatusUpdateRequired {
				isc.EXPECT().GetWorkloadInstance(&params, nil).Return(
					&instance.GetWorkloadInstanceOK{
						Payload: &workload_models.V1GetWorkloadInstanceResponse{
							Instance: i,
						},
					},
					nil,
				).Times(1)
			}

			isPodUpdated := podsTracker.handlePodUpdates(context.Background(), pod)

			if c.isPodStatusUpdateRequired {
				assert.Equal(t, isPodUpdated, true, "pod must be updated")
				assert.Equalf(t, c.expectedPodPhase, pod.Status.Phase, "pod status must be updated to %s", c.expectedPodPhase)
				assert.NotNil(t, pod.Status.StartTime, "podStatus start time must be set")
			} else {
				assert.Equal(t, isPodUpdated, false, "pod should not be updated")
			}
		})
	}
}

func TestStatusSyncup(t *testing.T) {
	podName := fmt.Sprintf("test-pod-%s", uuid.New().String())
	podNamespace := fmt.Sprintf("test-ns-%s", uuid.New().String())
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := context.Background()

	isc := mocks.NewInstanceClientService(mockController)
	stackPathClientMock := workload_client.EdgeCompute{Instance: isc}

	provider, err := createTestProvider(ctx, mocks.NewMockConfigMapLister(mockController), mocks.NewMockSecretLister(mockController), mocks.NewMockPodLister(mockController), &stackPathClientMock)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	podLister := mocks.NewMockPodLister(mockController)

	podsTracker := &PodsTracker{
		podLister:      podLister,
		updateCallback: func(p *v1.Pod) {},
		handler:        provider,
	}

	now := metav1.NewTime(time.Now())

	testCases := []struct {
		description             string
		expectedUpdate          bool
		podPhase                v1.PodPhase
		expectedPodPhase        v1.PodPhase
		expectedStatusReason    string
		expectedStatusMessage   string
		expectedContainerStatus v1.ContainerStatus
		initMockedCalls         func()
	}{
		{
			description:           "successfully terminates a running pod if the workload associated with that pod was not found (API returned 404 error)",
			expectedUpdate:        true,
			podPhase:              v1.PodRunning,
			expectedPodPhase:      v1.PodFailed,
			expectedStatusReason:  "NotFoundOnProvider",
			expectedStatusMessage: "the workload has been deleted from StackPath Edge Compute",
			expectedContainerStatus: v1.ContainerStatus{State: v1.ContainerState{
				Terminated: &v1.ContainerStateTerminated{
					ExitCode:   137,
					Reason:     "NotFoundOnProvider",
					Message:    "the workload has been deleted from StackPath Edge Compute",
					FinishedAt: now,
				},
			}},
			initMockedCalls: func() {
				params := instance.GetWorkloadInstanceParams{
					Context:      ctx,
					StackID:      provider.apiConfig.StackID,
					WorkloadID:   provider.getWorkloadSlug(podNamespace, podName),
					InstanceName: provider.getInstanceName(podNamespace, podName),
				}

				isc.EXPECT().GetWorkloadInstance(&params, nil).Return(nil,
					NewStackPathError(&APIError{
						statusCode: 404,
						message:    "Not found",
						requestID:  "123"})).Times(1)
			},
		},
		{
			description:           "ignore pod phase update if API returned non-404 error",
			expectedUpdate:        false,
			podPhase:              v1.PodRunning,
			expectedPodPhase:      v1.PodRunning,
			expectedStatusReason:  "",
			expectedStatusMessage: "",
			expectedContainerStatus: v1.ContainerStatus{
				State: v1.ContainerState{
					Running: &v1.ContainerStateRunning{},
				},
			},
			initMockedCalls: func() {
				params := instance.GetWorkloadInstanceParams{
					Context:      ctx,
					StackID:      provider.apiConfig.StackID,
					WorkloadID:   provider.getWorkloadSlug(podNamespace, podName),
					InstanceName: provider.getInstanceName(podNamespace, podName),
				}

				isc.EXPECT().GetWorkloadInstance(&params, nil).Return(nil,
					NewStackPathError(&APIError{
						statusCode: 500,
						message:    "Internal Server Error",
						requestID:  "123"})).Times(1)
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			c.initMockedCalls()
			pod := createTestPod(podName, podNamespace)
			pod.Status.Phase = c.podPhase
			pod.Status.ContainerStatuses = []v1.ContainerStatus{
				{
					Name: "test-container-1",
					State: v1.ContainerState{
						Running: &v1.ContainerStateRunning{},
					},
				},
				{
					Name: "test-container-2",
					State: v1.ContainerState{
						Terminated: &v1.ContainerStateTerminated{},
					},
				},
				{
					Name: "test-container-3",
					State: v1.ContainerState{
						Waiting: &v1.ContainerStateWaiting{},
					},
				},
			}

			isPodUpdated := podsTracker.handlePodUpdates(context.Background(), pod)

			if pod.Status.ContainerStatuses[0].State.Terminated != nil {
				pod.Status.ContainerStatuses[0].State.Terminated.FinishedAt = now
			}

			assert.Equal(t, c.expectedUpdate, isPodUpdated)
			assert.Equalf(t, c.expectedPodPhase, pod.Status.Phase, "pod status must be updated to %s", c.expectedPodPhase)
			assert.Equal(t, c.expectedStatusReason, pod.Status.Reason, "the pod's status message is not correct")
			assert.Equal(t, c.expectedStatusMessage, pod.Status.Message, "the pod's status message is not correct")
			assert.Equal(t, c.expectedContainerStatus.State.Running, pod.Status.ContainerStatuses[0].State.Running, "the container should not be in running state")
			assert.Equal(t, c.expectedContainerStatus.State.Waiting, pod.Status.ContainerStatuses[0].State.Waiting, "the container should not be in waiting state")
			assert.Equal(t, c.expectedContainerStatus.State.Terminated, pod.Status.ContainerStatuses[0].State.Terminated)
			assert.NotNil(t, pod.Status.ContainerStatuses[1].State.Terminated)
			assert.Nil(t, pod.Status.ContainerStatuses[1].State.Running)
			assert.Nil(t, pod.Status.ContainerStatuses[1].State.Waiting)
			assert.Nil(t, pod.Status.ContainerStatuses[2].State.Terminated)
			assert.Nil(t, pod.Status.ContainerStatuses[2].State.Running)
			assert.NotNil(t, pod.Status.ContainerStatuses[2].State.Waiting)

		})
	}
}

func TestBeginPodTracking(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2*time.Second))
	defer cancel()

	wsc := mocks.NewWorkloadsClientService(mockController)
	isc := mocks.NewInstanceClientService(mockController)
	stackPathClientMock := workload_client.EdgeCompute{Instance: isc, Workloads: wsc}
	podLister := mocks.NewMockPodLister(mockController)

	provider, err := createTestProvider(ctx, mocks.NewMockConfigMapLister(mockController), mocks.NewMockSecretLister(mockController), podLister, &stackPathClientMock)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	podsTracker := &PodsTracker{
		podLister:      podLister,
		updateCallback: func(p *v1.Pod) {},
		handler:        provider,
	}
	podStatusUpdateInterval = 1
	stalePodCleanupInterval = 1

	podLister.EXPECT().List(gomock.Any()).Return(nil, nil).AnyTimes()
	wsc.EXPECT().GetWorkloads(gomock.Any(), gomock.Any()).Return(&workloads.GetWorkloadsOK{
		Payload: &workload_models.V1GetWorkloadsResponse{
			Results: []*workload_models.V1Workload{},
		},
	}, nil).AnyTimes()
	podsTracker.BeginPodTracking(ctx)
}

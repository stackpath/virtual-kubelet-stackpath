package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
	"github.com/stackpath/virtual-kubelet-stackpath/internal/mocks"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestHttpHeaders(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}
	var tests = []struct {
		k8sHTTPHeaders        []v1.HTTPHeader
		expectedSPHttpHeaders workload_models.V1StringMapEntry
	}{
		{
			k8sHTTPHeaders:        []v1.HTTPHeader{{Name: "test1", Value: "value1"}, {Name: "test2", Value: "value2"}},
			expectedSPHttpHeaders: workload_models.V1StringMapEntry{"test1": "value1", "test2": "value2"},
		},
	}

	for _, test := range tests {
		workloadHTTP := provider.getHTTPHeadersFrom(test.k8sHTTPHeaders)
		for key, value := range workloadHTTP {
			assert.Equal(t, test.expectedSPHttpHeaders[key], value)
		}
	}
}

func TestProbe(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}
	var tests = []struct {
		description     string
		k8sProbe        *v1.Probe
		expectedSPProbe *workload_models.V1Probe
		containerPorts  []v1.ContainerPort
		err             error
	}{
		{
			description: "http probe",
			k8sProbe: &v1.Probe{
				FailureThreshold:    int32(3),
				InitialDelaySeconds: int32(3),
				TimeoutSeconds:      int32(4),
				PeriodSeconds:       int32(4),
				SuccessThreshold:    int32(5),
				ProbeHandler: v1.ProbeHandler{
					HTTPGet: &v1.HTTPGetAction{
						Path:        "/some/path",
						Port:        intstr.IntOrString{Type: intstr.Int, IntVal: int32(8000)},
						Host:        "some-host",
						Scheme:      v1.URISchemeHTTP,
						HTTPHeaders: []v1.HTTPHeader{{Name: "test-name", Value: "test-value"}},
					},
				},
			},
			expectedSPProbe: &workload_models.V1Probe{
				FailureThreshold: int32(3),
				HTTPGet: &workload_models.V1HTTPGetAction{
					HTTPHeaders: workload_models.V1StringMapEntry{"test-name": "test-value"},
					Path:        "/some/path",
					Port:        int32(8000),
					Scheme:      "HTTP",
				},
				InitialDelaySeconds: int32(3),
				PeriodSeconds:       int32(4),
				SuccessThreshold:    int32(5),
				TCPSocket:           nil,
				TimeoutSeconds:      int32(4),
			},
		},
		{
			description: "tcp probe with int port",
			k8sProbe: &v1.Probe{
				ProbeHandler: v1.ProbeHandler{
					TCPSocket: &v1.TCPSocketAction{
						Port: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
					},
				},
			},
			expectedSPProbe: &workload_models.V1Probe{
				TCPSocket: &workload_models.V1TCPSocketAction{
					Port: 80,
				},
			},
		},
		{
			description: "tcp probe with string port that doesn't exist",
			k8sProbe: &v1.Probe{
				ProbeHandler: v1.ProbeHandler{
					TCPSocket: &v1.TCPSocketAction{
						Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
					},
				},
			},
			err: fmt.Errorf("unable to find named port: http"),
		},
		{
			description: "tcp probe with string port that exists",
			k8sProbe: &v1.Probe{
				ProbeHandler: v1.ProbeHandler{
					TCPSocket: &v1.TCPSocketAction{
						Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
					},
				},
			},
			containerPorts: []v1.ContainerPort{
				{
					Name:          "http",
					ContainerPort: 80,
				},
			},
			expectedSPProbe: &workload_models.V1Probe{
				TCPSocket: &workload_models.V1TCPSocketAction{
					Port: 80,
				},
			},
		},
		{
			description: "http probe with string port that doesn't exist",
			k8sProbe: &v1.Probe{
				ProbeHandler: v1.ProbeHandler{
					HTTPGet: &v1.HTTPGetAction{
						Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
					},
				},
			},
			err: fmt.Errorf("unable to find named port: http"),
		},
		{
			description: "http probe with string port that exists",
			k8sProbe: &v1.Probe{
				ProbeHandler: v1.ProbeHandler{
					HTTPGet: &v1.HTTPGetAction{
						Path:        "/some/path",
						Host:        "some-host",
						Scheme:      v1.URISchemeHTTP,
						HTTPHeaders: []v1.HTTPHeader{{Name: "test-name", Value: "test-value"}},
						Port:        intstr.IntOrString{Type: intstr.String, StrVal: "custom-port"},
					},
				},
			},
			containerPorts: []v1.ContainerPort{
				{
					Name:          "custom-port",
					ContainerPort: 123,
				},
			},
			expectedSPProbe: &workload_models.V1Probe{
				HTTPGet: &workload_models.V1HTTPGetAction{
					Port:        int32(123),
					HTTPHeaders: workload_models.V1StringMapEntry{"test-name": "test-value"},
					Path:        "/some/path",
					Scheme:      "HTTP",
				},
			},
		},
		{
			description:     "empty probe returns nil",
			k8sProbe:        &v1.Probe{},
			expectedSPProbe: nil,
		},
		{
			description:     "nil probe returns nil",
			k8sProbe:        (*v1.Probe)(nil),
			expectedSPProbe: nil,
		},
		{
			description: "unsupported exec probe returns nil",
			k8sProbe: &v1.Probe{
				ProbeHandler: v1.ProbeHandler{
					Exec: &v1.ExecAction{
						Command: []string{"some", "command"},
					},
				},
			},
			expectedSPProbe: nil,
		},
		{
			description: "unsupported grpc probe returns nil",
			k8sProbe: &v1.Probe{
				ProbeHandler: v1.ProbeHandler{
					GRPC: &v1.GRPCAction{
						Port: 80,
					},
				},
			},
			expectedSPProbe: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			workloadProbe, err := provider.getWorkloadContainerProbeFrom(test.k8sProbe, test.containerPorts)
			if err != nil {
				assert.Equal(t, test.err, err, test.description)
			} else {
				assert.Equal(t, test.expectedSPProbe, workloadProbe, test.description)
			}
		})
	}
}

func TestContainerLifecycle(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}
	var tests = []struct {
		description       string
		k8sContainer      *v1.Container
		expectedLifecycle *workload_models.V1ContainerLifecycle
		err               error
	}{
		{
			description:       "no lifecycle returned for a container that doesn't have lifecycle events",
			k8sContainer:      &v1.Container{},
			expectedLifecycle: nil,
			err:               nil,
		},
		{
			description: "no lifecycle returned for container with an empty lifecycle handler",
			k8sContainer: &v1.Container{Lifecycle: &v1.Lifecycle{
				PostStart: &v1.LifecycleHandler{},
				PreStop:   &v1.LifecycleHandler{},
			}},
			expectedLifecycle: nil,
			err:               nil,
		},
		{
			description: "lifecycle with TCP socket action that has a string port that doesn't exist",
			k8sContainer: &v1.Container{
				Lifecycle: &v1.Lifecycle{
					PostStart: &v1.LifecycleHandler{
						TCPSocket: &v1.TCPSocketAction{
							Port: intstr.IntOrString{
								Type: intstr.String, StrVal: "http",
							},
						},
					},
				},
			},
			expectedLifecycle: nil,
			err:               fmt.Errorf("unable to find named port: http"),
		},
		{
			description: "lifecycle with HTTP get action that has a string port that doesn't exist",
			k8sContainer: &v1.Container{
				Lifecycle: &v1.Lifecycle{
					PreStop: &v1.LifecycleHandler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.IntOrString{
								Type: intstr.String, StrVal: "http",
							},
						},
					},
				},
			},
			expectedLifecycle: nil,
			err:               fmt.Errorf("unable to find named port: http"),
		},
		{
			description: "exec lifecycle handler is not supported",
			k8sContainer: &v1.Container{
				Lifecycle: &v1.Lifecycle{
					PostStart: &v1.LifecycleHandler{
						Exec: &v1.ExecAction{
							Command: []string{"test"},
						},
					},
				},
			},
			expectedLifecycle: nil,
			err:               nil,
		},
		{
			description: "post start lifecycle",
			k8sContainer: &v1.Container{Lifecycle: &v1.Lifecycle{
				PostStart: &v1.LifecycleHandler{
					HTTPGet: &v1.HTTPGetAction{
						Path:        "path1",
						Port:        intstr.IntOrString{Type: intstr.Int, IntVal: int32(8000)},
						Scheme:      v1.URISchemeHTTP,
						HTTPHeaders: []v1.HTTPHeader{{Name: "test1", Value: "value1"}, {Name: "test2", Value: "value2"}},
					},
					TCPSocket: &v1.TCPSocketAction{
						Port: intstr.IntOrString{Type: intstr.Int, IntVal: int32(8001)},
					},
				},
				PreStop: &v1.LifecycleHandler{
					HTTPGet: &v1.HTTPGetAction{
						Path:        "path2",
						Port:        intstr.IntOrString{Type: intstr.Int, IntVal: int32(8000)},
						Scheme:      v1.URISchemeHTTP,
						HTTPHeaders: []v1.HTTPHeader{{Name: "test3", Value: "value3"}, {Name: "test4", Value: "value4"}},
					},
					TCPSocket: &v1.TCPSocketAction{
						Port: intstr.IntOrString{Type: intstr.Int, IntVal: int32(8002)},
					},
				},
			}},
			expectedLifecycle: &workload_models.V1ContainerLifecycle{
				PostStart: &workload_models.V1ContainerLifecycleHandler{
					HTTPGet: &workload_models.V1HTTPGetAction{
						Path:        "path1",
						Port:        8000,
						Scheme:      "HTTP",
						HTTPHeaders: workload_models.V1StringMapEntry{"test1": "value1", "test2": "value2"},
					},
					TCPSocket: &workload_models.V1TCPSocketAction{
						Port: 8001,
					},
				},
				PreStop: &workload_models.V1ContainerLifecycleHandler{
					HTTPGet: &workload_models.V1HTTPGetAction{
						Path:        "path2",
						Port:        8000,
						Scheme:      "HTTP",
						HTTPHeaders: workload_models.V1StringMapEntry{"test3": "value3", "test4": "value4"},
					},
					TCPSocket: &workload_models.V1TCPSocketAction{
						Port: 8002,
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			containerLifecycle, err := provider.getWorkloadContainerLifecycle(test.k8sContainer)
			if err != nil {
				assert.Equal(t, test.err, err, test.description)
			} else {
				assert.Equal(t, test.expectedLifecycle, containerLifecycle, test.description)
			}
		})
	}
}

func TestVolumeMounts(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}
	var tests = []struct {
		description            string
		k8sVolumeMounts        []v1.VolumeMount
		expectedSPVolumeMounts []*workload_models.V1InstanceVolumeMount
	}{
		{
			description:            "test volume mount",
			k8sVolumeMounts:        []v1.VolumeMount{{Name: "test1", MountPath: "/path/1"}, {Name: "test2", MountPath: "/path/2"}},
			expectedSPVolumeMounts: []*workload_models.V1InstanceVolumeMount{{MountPath: "/path/1", Slug: "test1"}, {MountPath: "/path/2", Slug: "test2"}},
		},
		{
			description:            "test volume mount for k8s default secret is skipped",
			k8sVolumeMounts:        []v1.VolumeMount{{Name: "test1", MountPath: "/path/1"}, {Name: "kube-api-access-vlkdv", MountPath: "/var/run/secrets/kubernetes.io/serviceaccount"}},
			expectedSPVolumeMounts: []*workload_models.V1InstanceVolumeMount{{MountPath: "/path/1", Slug: "test1"}},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			workloadVolumeMounts := provider.getWorkloadContainerVolumeMountsFrom(test.k8sVolumeMounts)
			for index, volume := range workloadVolumeMounts {
				assert.Equal(t, test.expectedSPVolumeMounts[index].MountPath, volume.MountPath)
				assert.Equal(t, test.expectedSPVolumeMounts[index].Slug, volume.Slug)
			}
		})
	}
}

func TestContainerEnv(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}
	var tests = []struct {
		description   string
		k8sEnv        []v1.EnvVar
		expectedSPEnv workload_models.V1EnvironmentVariableMapEntry
	}{
		{
			description:   "happy case env vars",
			k8sEnv:        []v1.EnvVar{{Name: "some-name", Value: "some-value"}, {Name: "some-name1", Value: "some-value1"}},
			expectedSPEnv: workload_models.V1EnvironmentVariableMapEntry{"some-name": workload_models.V1EnvironmentVariable{Value: "some-value"}, "some-name1": workload_models.V1EnvironmentVariable{Value: "some-value1"}},
		},
		{
			description:   "empty env vars results in empty env vars",
			k8sEnv:        []v1.EnvVar{},
			expectedSPEnv: workload_models.V1EnvironmentVariableMapEntry{},
		},
		{
			description: "ensure K8S default environment variables are dropped",
			k8sEnv: []v1.EnvVar{
				{
					Name:  "KUBERNETES_PORT",
					Value: "tcp://10.43.0.1:443",
				},
				{
					Name:  "KUBERNETES_PORT_443_TCP",
					Value: "tcp://10.43.0.1:443",
				},
				{
					Name:  "KUBERNETES_PORT_443_TCP_ADDR",
					Value: "10.43.0.1",
				},
				{
					Name:  "KUBERNETES_PORT_443_TCP_PORT",
					Value: "443",
				},
				{
					Name:  "KUBERNETES_PORT_443_TCP_PROTO",
					Value: "tcp",
				},
				{
					Name:  "KUBERNETES_SERVICE_HOST",
					Value: "10.43.0.1",
				},
				{
					Name:  "KUBERNETES_SERVICE_PORT",
					Value: "443",
				},
				{
					Name:  "KUBERNETES_SERVICE_PORT_HTTPS",
					Value: "443",
				},
			},
			expectedSPEnv: workload_models.V1EnvironmentVariableMapEntry{},
		},
		{
			description: "value from not supported is skipped",
			k8sEnv: []v1.EnvVar{
				{
					Name:      "VALUE_FROM_ENV",
					ValueFrom: &v1.EnvVarSource{ConfigMapKeyRef: &v1.ConfigMapKeySelector{Key: "reference-key"}},
				},
			},
			expectedSPEnv: workload_models.V1EnvironmentVariableMapEntry{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			workloadEnvs := provider.getWorkloadContainerEnvFrom(test.k8sEnv)
			for index, env := range workloadEnvs {
				assert.Equal(t, test.expectedSPEnv[index].Value, env.Value)
			}
		})
	}
}

func TestContainerPorts(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}
	var tests = []struct {
		k8sPort        []v1.ContainerPort
		expectedSPPort workload_models.V1InstancePortMapEntry
		err            error
	}{
		{
			k8sPort:        []v1.ContainerPort{{Name: "some-name", Protocol: v1.ProtocolTCP, ContainerPort: int32(8000)}, {Name: "", Protocol: v1.ProtocolUDP, ContainerPort: int32(7000)}},
			expectedSPPort: workload_models.V1InstancePortMapEntry{"some-name": workload_models.V1InstancePort{Port: int32(8000), Protocol: "TCP"}, "default": workload_models.V1InstancePort{Port: int32(7000), Protocol: "UDP"}},
			err:            nil,
		},
	}

	for _, test := range tests {
		workloadPorts := provider.getWorkloadContainerPortsFrom(test.k8sPort)
		for index, port := range workloadPorts {
			assert.Equal(t, test.expectedSPPort[index].Protocol, port.Protocol)
		}

	}
}

func TestInstanceSizes(t *testing.T) {
	var tests = []struct {
		description             string
		k8sResourceRequirements v1.ResourceRequirements
		expectedResources       workload_models.V1StringMapEntry
	}{
		{
			description:             "no resource requirements result in sp-1",
			k8sResourceRequirements: v1.ResourceRequirements{},
			expectedResources:       containerResourcesSP1,
		},
		{
			description:             "small cpu resource requirements without memory result in sp-1",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"cpu": *resource.NewMilliQuantity(100, resource.DecimalSI)}},
			expectedResources:       containerResourcesSP1,
		},
		{
			description:             "small cpu resource high memory result in sp-5",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"cpu": *resource.NewMilliQuantity(100, resource.DecimalSI), "memory": *resource.NewQuantity(24*1024*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP5,
		},
		{
			description:             "small cpu resource memory fit for sp-4 result in sp-4",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"cpu": *resource.NewMilliQuantity(100, resource.DecimalSI), "memory": *resource.NewQuantity(10*1024*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP4,
		},
		{
			description:             "small cpu resource memory fit for sp-3 result in sp-3",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"cpu": *resource.NewMilliQuantity(100, resource.DecimalSI), "memory": *resource.NewQuantity(8*1024*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP3,
		},
		{
			description:             "small cpu resource memory fit for sp-2 result in sp-2",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"cpu": *resource.NewMilliQuantity(100, resource.DecimalSI), "memory": *resource.NewQuantity(4*1024*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP2,
		},
		{
			description:             "small cpu resource memory fit for sp-1 result in sp-1",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"cpu": *resource.NewMilliQuantity(100, resource.DecimalSI), "memory": *resource.NewQuantity(2*1024*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP1,
		},
		{
			description:             "small memory resource requirements without cpu result in sp-1",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"memory": *resource.NewQuantity(100*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP1,
		},
		{
			description:             "small memory resource requirements with high cpu result in sp-5",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"cpu": *resource.NewQuantity(32, resource.DecimalSI), "memory": *resource.NewQuantity(100*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP5,
		},
		{
			description:             "small memory resource requirements with with cpu fit for sp-4 result in sp-4",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"cpu": *resource.NewQuantity(4, resource.DecimalSI), "memory": *resource.NewQuantity(100*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP4,
		},
		{
			description:             "small memory resource requirements with cpu fit for sp-2 and sp-3 result in sp-2",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"cpu": *resource.NewQuantity(2, resource.DecimalSI), "memory": *resource.NewQuantity(100*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP2,
		},
		{
			description:             "small memory resource requirements with cpu fit for sp-1 result in sp-1",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"cpu": *resource.NewQuantity(1, resource.DecimalSI), "memory": *resource.NewQuantity(100*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP1,
		},
		{
			description:             "cpu matching sp-2 and sp-3 with memory matching sp-3 results in sp-3",
			k8sResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{"cpu": *resource.NewQuantity(2, resource.DecimalSI), "memory": *resource.NewQuantity(5*1024*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP3,
		},
		{
			description:             "using limits instead of requests works as well",
			k8sResourceRequirements: v1.ResourceRequirements{Limits: v1.ResourceList{"cpu": *resource.NewQuantity(2, resource.DecimalSI), "memory": *resource.NewQuantity(5*1024*1024*1024, resource.BinarySI)}},
			expectedResources:       containerResourcesSP3,
		},
		{
			description: "requests and limits maximum calculated correctly",
			k8sResourceRequirements: v1.ResourceRequirements{
				Limits:   v1.ResourceList{"cpu": *resource.NewQuantity(8, resource.DecimalSI), "memory": *resource.NewQuantity(5*1024*1024*1024, resource.BinarySI)},
				Requests: v1.ResourceList{"cpu": *resource.NewQuantity(2, resource.DecimalSI), "memory": *resource.NewQuantity(5*1024*1024*1024, resource.BinarySI)},
			},
			expectedResources: containerResourcesSP5,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := toSPInstanceSize(test.k8sResourceRequirements)
			assert.Equal(t, test.expectedResources, result)
		})
	}
}

func TestWorkloadVolumes(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	var tests = []struct {
		description string
		pod         v1.Pod
		expected    []workload_models.V1VolumeClaim
		err         error
		len         int
	}{
		{
			description: "empty pod return ok without volumes",
			pod:         v1.Pod{},
			expected:    []workload_models.V1VolumeClaim{},
			err:         nil,
			len:         0,
		},
		{
			description: "pod with empty dir returns no volumes since it's not supported",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "empty-dir",
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{
									SizeLimit: resource.NewQuantity(1*1024*1024*1024, resource.DecimalSI),
									Medium:    v1.StorageMediumDefault,
								},
							},
						},
					},
				},
			},
			expected: []workload_models.V1VolumeClaim{},
			err:      nil,
			len:      0,
		},
		{
			description: "pod with invalid csi driver returns 0 volumes",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "valid-volume",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver:           "wrong-driver",
									VolumeAttributes: map[string]string{"size": "1Gi"},
								},
							},
						},
					},
				},
			},
			expected: []workload_models.V1VolumeClaim{},
			err:      nil,
			len:      0,
		},
		{
			description: "pod with valid volume returns 1 volume",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "valid-volume",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver:           stackpathVirtualKubeletCSIDriver,
									VolumeAttributes: map[string]string{"size": "1Gi"},
								},
							},
						},
					},
				},
			},
			expected: []workload_models.V1VolumeClaim{
				{
					Metadata: &workload_models.V1Metadata{},
					Name:     "valid-volume",
					Slug:     "valid-volume",
					Spec: &workload_models.V1VolumeClaimSpec{
						Resources: &workload_models.V1ResourceRequirements{Requests: workload_models.V1StringMapEntry{"storage": "1Gi"}, Limits: workload_models.V1StringMapEntry{"storage": "1Gi"}},
					},
				},
			},
			err: nil,
			len: 1,
		},
		{
			description: "pod with valid and invalid volumes only returns valid volumes",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "valid-volume-1",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver:           stackpathVirtualKubeletCSIDriver,
									VolumeAttributes: map[string]string{"size": "1Gi"},
								},
							},
						},
						{
							Name: "invalid-volume",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver:           "wrong-driver",
									VolumeAttributes: map[string]string{"size": "1Gi"},
								},
							},
						},
						{
							Name: "valid-volume-2",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver:           stackpathVirtualKubeletCSIDriver,
									VolumeAttributes: map[string]string{"size": "1Gi"},
								},
							},
						},
					},
				},
			},
			expected: []workload_models.V1VolumeClaim{
				{
					Metadata: &workload_models.V1Metadata{},
					Name:     "valid-volume-1",
					Slug:     "valid-volume-1",
					Spec: &workload_models.V1VolumeClaimSpec{
						Resources: &workload_models.V1ResourceRequirements{Requests: workload_models.V1StringMapEntry{"storage": "1Gi"}, Limits: workload_models.V1StringMapEntry{"storage": "1Gi"}},
					},
				},
				{
					Metadata: &workload_models.V1Metadata{},
					Name:     "valid-volume-2",
					Slug:     "valid-volume-2",
					Spec: &workload_models.V1VolumeClaimSpec{
						Resources: &workload_models.V1ResourceRequirements{Requests: workload_models.V1StringMapEntry{"storage": "1Gi"}, Limits: workload_models.V1StringMapEntry{"storage": "1Gi"}},
					},
				},
			},
			err: nil,
			len: 2,
		},
		{
			description: "pod with valid volume over max size returns volume with 1000Gi",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "valid-volume",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver:           stackpathVirtualKubeletCSIDriver,
									VolumeAttributes: map[string]string{"size": "2000Gi"},
								},
							},
						},
					},
				},
			},
			expected: []workload_models.V1VolumeClaim{
				{
					Metadata: &workload_models.V1Metadata{},
					Name:     "valid-volume",
					Slug:     "valid-volume",
					Spec: &workload_models.V1VolumeClaimSpec{
						Resources: &workload_models.V1ResourceRequirements{Requests: workload_models.V1StringMapEntry{"storage": "1000Gi"}, Limits: workload_models.V1StringMapEntry{"storage": "1000Gi"}},
					},
				},
			},
			err: nil,
			len: 1,
		},
		{
			description: "pod without attributes returns default 1Gi volume",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "valid-volume",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver: stackpathVirtualKubeletCSIDriver,
								},
							},
						},
					},
				},
			},
			expected: []workload_models.V1VolumeClaim{
				{
					Metadata: &workload_models.V1Metadata{},
					Name:     "valid-volume",
					Slug:     "valid-volume",
					Spec: &workload_models.V1VolumeClaimSpec{
						Resources: &workload_models.V1ResourceRequirements{Requests: workload_models.V1StringMapEntry{"storage": "1Gi"}, Limits: workload_models.V1StringMapEntry{"storage": "1Gi"}},
					},
				},
			},
			err: nil,
			len: 1,
		},
		{
			description: "pod with empty attributes returns default 1Gi volume",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "valid-volume",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver:           stackpathVirtualKubeletCSIDriver,
									VolumeAttributes: map[string]string{},
								},
							},
						},
					},
				},
			},
			expected: []workload_models.V1VolumeClaim{
				{
					Metadata: &workload_models.V1Metadata{},
					Name:     "valid-volume",
					Slug:     "valid-volume",
					Spec: &workload_models.V1VolumeClaimSpec{
						Resources: &workload_models.V1ResourceRequirements{Requests: workload_models.V1StringMapEntry{"storage": "1Gi"}, Limits: workload_models.V1StringMapEntry{"storage": "1Gi"}},
					},
				},
			},
			err: nil,
			len: 1,
		},
		{
			description: "pod with non relevant attributes returns default 1Gi volume",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "valid-volume",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver: stackpathVirtualKubeletCSIDriver,
									VolumeAttributes: map[string]string{
										"some": "variable",
									},
								},
							},
						},
					},
				},
			},
			expected: []workload_models.V1VolumeClaim{
				{
					Metadata: &workload_models.V1Metadata{},
					Name:     "valid-volume",
					Slug:     "valid-volume",
					Spec: &workload_models.V1VolumeClaimSpec{
						Resources: &workload_models.V1ResourceRequirements{Requests: workload_models.V1StringMapEntry{"storage": "1Gi"}, Limits: workload_models.V1StringMapEntry{"storage": "1Gi"}},
					},
				},
			},
			err: nil,
			len: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			volumes, err := provider.getWorkloadVolumesFrom(&test.pod)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Len(t, volumes, test.len)
				for index, volume := range volumes {
					assert.Equal(t, test.expected[index], *volume)
				}
			}
		})
	}
}

func TestGetWorkloadFromErrorsWithInvalidVolumeSize(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	pod := v1.Pod{
		Spec: v1.PodSpec{
			Volumes: []v1.Volume{
				{
					Name: "valid-volume",
					VolumeSource: v1.VolumeSource{
						CSI: &v1.CSIVolumeSource{
							Driver:           stackpathVirtualKubeletCSIDriver,
							VolumeAttributes: map[string]string{"size": "invalid"},
						},
					},
				},
			},
		},
	}

	_, err = provider.getWorkloadFrom(&pod)
	assert.ErrorContains(t, err, "quantities must match the regular expression")
}

func TestGetWorkloadFromErrorsWithInvalidLivenessProbe(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}
	testCases := []struct {
		description   string
		pod           *v1.Pod
		expectedError error
	}{
		{
			description: "fails to create get a workload form the pod due to bad liveness probe port set for the container",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "valid-volume",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver:           stackpathVirtualKubeletCSIDriver,
									VolumeAttributes: map[string]string{"size": "1Gi"},
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:  "test",
							Image: "nginx:latest",
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: int32(80),
								},
							},
							LivenessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									HTTPGet: &v1.HTTPGetAction{
										Port: intstr.IntOrString{Type: intstr.String, StrVal: "not-exists"},
									},
								},
							},
						},
					},
				},
			},
			expectedError: errors.New("unable to find named port: not-exists"),
		}, {
			description: "fails to create get a workload form the pod due to bad liveness probe port set for the init container",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "valid-volume",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver:           stackpathVirtualKubeletCSIDriver,
									VolumeAttributes: map[string]string{"size": "1Gi"},
								},
							},
						},
					},
					InitContainers: []v1.Container{
						{
							Name:  "test",
							Image: "nginx:latest",
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: int32(80),
								},
							},
							LivenessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									HTTPGet: &v1.HTTPGetAction{
										Port: intstr.IntOrString{Type: intstr.String, StrVal: "not-exists"},
									},
								},
							},
						},
					},
				},
			},
			expectedError: errors.New("unable to find named port: not-exists"),
		},
	}
	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			_, err = provider.getWorkloadFrom(c.pod)
			assert.Equal(t, err, c.expectedError)
		},
		)
	}
}

func TestGetWorkloadFromErrorsWithInvalidReadinessProbe(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}
	pod := v1.Pod{
		Spec: v1.PodSpec{
			Volumes: []v1.Volume{
				{
					Name: "valid-volume",
					VolumeSource: v1.VolumeSource{
						CSI: &v1.CSIVolumeSource{
							Driver:           stackpathVirtualKubeletCSIDriver,
							VolumeAttributes: map[string]string{"size": "1Gi"},
						},
					},
				},
			},
			Containers: []v1.Container{
				{
					Name:  "test",
					Image: "nginx:latest",
					Ports: []v1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: int32(80),
						},
					},
					ReadinessProbe: &v1.Probe{
						ProbeHandler: v1.ProbeHandler{
							HTTPGet: &v1.HTTPGetAction{
								Port: intstr.IntOrString{Type: intstr.String, StrVal: "not-exists"},
							},
						},
					},
				},
			},
		},
	}

	_, err = provider.getWorkloadFrom(&pod)
	assert.ErrorContains(t, err, "unable to find named port")
}

func TestGetWorkloadFromErrorsWithInvalidStartupProbe(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}
	pod := v1.Pod{
		Spec: v1.PodSpec{
			Volumes: []v1.Volume{
				{
					Name: "valid-volume",
					VolumeSource: v1.VolumeSource{
						CSI: &v1.CSIVolumeSource{
							Driver:           stackpathVirtualKubeletCSIDriver,
							VolumeAttributes: map[string]string{"size": "1Gi"},
						},
					},
				},
			},
			Containers: []v1.Container{
				{
					Name:  "test",
					Image: "nginx:latest",
					Ports: []v1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: int32(80),
						},
					},
					StartupProbe: &v1.Probe{
						ProbeHandler: v1.ProbeHandler{
							HTTPGet: &v1.HTTPGetAction{
								Port: intstr.IntOrString{Type: intstr.String, StrVal: "not-exists"},
							},
						},
					},
				},
			},
		},
	}

	_, err = provider.getWorkloadFrom(&pod)
	assert.ErrorContains(t, err, "unable to find named port")
}

func TestGetWorkloadFromErrorsWithImagePullSecret(t *testing.T) {
	ctx := context.Background()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	pod := v1.Pod{
		Spec: v1.PodSpec{
			Volumes: []v1.Volume{
				{
					Name: "valid-volume",
					VolumeSource: v1.VolumeSource{
						CSI: &v1.CSIVolumeSource{
							Driver:           stackpathVirtualKubeletCSIDriver,
							VolumeAttributes: map[string]string{"size": "1Gi"},
						},
					},
				},
			},
			Containers: []v1.Container{
				{
					Name:  "test",
					Image: "nginx:latest",
					Ports: []v1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: int32(80),
						},
					},
				},
			},
			ImagePullSecrets: []v1.LocalObjectReference{
				{
					Name: "image-pull-secret",
				},
			},
		},
	}
	pod.Namespace = "test"

	secretListerMock := mocks.NewMockSecretLister(mockController)
	secretNamespaceListerMock := mocks.NewMockSecretNamespaceLister(gomock.NewController(t))
	dockerConfig := dockerConfigJSON{
		Auths: dockerConfig{
			"server": dockerConfigEntry{
				Username: "user",
				Password: "password",
				Email:    "user@gmail.com",
				Auth:     "some-encoded-string",
			},
		},
	}
	secret, err := json.Marshal(dockerConfig)
	if err != nil {
		t.Errorf("failed to marshal test docker config %+v", dockerConfig)
	}

	secretNamespaceListerMock.EXPECT().Get("image-pull-secret").Return(
		&v1.Secret{
			Type: "kubernetes.io/dockercfg",
			Data: map[string][]byte{
				".dockercfg": secret,
			},
		},
		nil,
	)
	secretListerMock.EXPECT().Secrets("test").Return(secretNamespaceListerMock)
	provider, err := createTestProvider(ctx, nil, secretListerMock, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	_, err = provider.getWorkloadFrom(&pod)
	assert.ErrorContains(t, err, "legacy format kubernetes.io/dockercfg is not supported")
}

func TestContainerWithCommand(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	var tests = []struct {
		description string
		container   v1.Container
		expected    workload_models.V1ContainerSpec
		err         string
	}{

		{
			description: "command and args concatenated successfully",
			container: v1.Container{
				Command: []string{
					"test",
					"command",
				},
				Args: []string{
					"test",
					"args",
				},
			},
			expected: workload_models.V1ContainerSpec{
				Command: []string{
					"test",
					"command",
				},
				Args: []string{
					"test",
					"args",
				},
			},
		},
		{
			description: "command with no args result ok",
			container: v1.Container{
				Command: []string{
					"test",
					"command",
				},
			},
			expected: workload_models.V1ContainerSpec{
				Command: []string{
					"test",
					"command",
				},
			},
		},
		{
			description: "args and no command return error",
			container: v1.Container{
				Args: []string{
					"test",
					"command",
				},
			},
			expected: workload_models.V1ContainerSpec{
				Command: ([]string)(nil),
			},
			err: "args not supported without command",
		},
	}

	for _, test := range tests {
		containerSpec, err := provider.getWorkloadContainerSpecFrom(&test.container)
		if err != nil {
			assert.ErrorContains(t, err, test.err, test.description)
		} else {
			assert.Equal(t, test.expected.Command, containerSpec.Command, test.description)
			assert.Equal(t, test.expected.Args, containerSpec.Args, test.description)
		}
	}
}

func TestGetWorkloadFrom(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	var tests = []struct {
		description string
		pod         v1.Pod
		expected    *workload_models.V1Workload
		err         string
	}{
		{
			description: "successfully translates from v1.Pod to workload_model.V1Workload",
			pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
					Annotations: map[string]string{
						"workload.platform.stackpath.net/remote-management": "true",
						"anycast.platform.stackpath.net":                    "true",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Image:   "test-image",
							Command: []string{"/sh/bash"},
							Name:    "nginx",
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
								},
							},
							TerminationMessagePath:   "/test/path",
							TerminationMessagePolicy: v1.TerminationMessageReadFile,
						},
					},
					InitContainers: []v1.Container{
						{
							Image:   "init-image",
							Command: []string{"/sh/bash"},
							Name:    "init",
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
			expected: &workload_models.V1Workload{
				Name: "namespace-name",
				Slug: "namespace-name",
				Metadata: &workload_models.V1Metadata{
					Labels: map[string]string{
						podNameLabelKey:      "name",
						podNamespaceLabelKey: "namespace",
						nodeNameLabelKey:     "vk-mock",
					},
					Annotations: map[string]string{
						"workload.platform.stackpath.net/remote-management": "true",
						"anycast.platform.stackpath.net":                    "true",
					},
				},
				Spec: &workload_models.V1WorkloadSpec{
					Containers: workload_models.V1ContainerSpecMapEntry{
						"nginx": workload_models.V1ContainerSpec{
							Env:     workload_models.V1EnvironmentVariableMapEntry{},
							Image:   "test-image",
							Command: []string{"/sh/bash"},
							Ports: workload_models.V1InstancePortMapEntry{"http": workload_models.V1InstancePort{
								Port:                        8080,
								EnableImplicitNetworkPolicy: false,
								Protocol:                    "TCP",
							}},
							Resources: &workload_models.V1ResourceRequirements{
								Limits: workload_models.V1StringMapEntry{
									"cpu":    "1",
									"memory": "2Gi",
								},
								Requests: workload_models.V1StringMapEntry{
									"cpu":    "1",
									"memory": "2Gi",
								},
							},
							VolumeMounts:             []*workload_models.V1InstanceVolumeMount{},
							TerminationMessagePath:   "/test/path",
							TerminationMessagePolicy: workload_models.V1ContainerTerminationMessagePolicyFILE.Pointer(),
						},
					},
					InitContainers: workload_models.V1ContainerSpecMapEntry{
						"init": workload_models.V1ContainerSpec{
							Env:     workload_models.V1EnvironmentVariableMapEntry{},
							Image:   "init-image",
							Command: []string{"/sh/bash"},
							Ports: workload_models.V1InstancePortMapEntry{"http": workload_models.V1InstancePort{
								Port:                        8080,
								EnableImplicitNetworkPolicy: false,
								Protocol:                    "TCP",
							}},
							Resources: &workload_models.V1ResourceRequirements{
								Limits: workload_models.V1StringMapEntry{
									"cpu":    "1",
									"memory": "2Gi",
								},
								Requests: workload_models.V1StringMapEntry{
									"cpu":    "1",
									"memory": "2Gi",
								},
							},
							VolumeMounts: []*workload_models.V1InstanceVolumeMount{},
						},
					},
					ImagePullCredentials: workload_models.V1WrappedImagePullCredentials{},
					NetworkInterfaces: []*workload_models.V1NetworkInterface{
						{
							EnableOneToOneNat: true,
							IPFamilies:        []*workload_models.V1IPFamily{workload_models.NewV1IPFamily("IPv4")},
							IPV6Subnet:        "",
							Network:           "default",
							Subnet:            "",
						},
					},
					VolumeClaimTemplates: []*workload_models.V1VolumeClaim{},
					Runtime:              &workload_models.V1WorkloadInstanceRuntimeSettings{},
				},
				Targets: workload_models.V1TargetMapEntry{
					"city-code": workload_models.V1Target{
						Spec: &workload_models.V1TargetSpec{
							DeploymentScope: "cityCode",
							Deployments: &workload_models.V1DeploymentSpec{
								MaxReplicas: 1,
								MinReplicas: 1,
								Selectors: []*workload_models.V1MatchExpression{
									{
										Key:      "cityCode",
										Operator: "in",
										Values:   []string{testCityCode},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			description: "fails to translate due to an error in container's lifecycle",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "valid-volume",
							VolumeSource: v1.VolumeSource{
								CSI: &v1.CSIVolumeSource{
									Driver:           stackpathVirtualKubeletCSIDriver,
									VolumeAttributes: map[string]string{"size": "1Gi"},
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:  "test",
							Image: "nginx:latest",
							Lifecycle: &v1.Lifecycle{
								PostStart: &v1.LifecycleHandler{
									TCPSocket: &v1.TCPSocketAction{
										Port: intstr.IntOrString{
											Type: intstr.String, StrVal: "http",
										},
									},
								},
							},
						},
					},
				},
			},
			expected: nil,
			err:      "unable to find named port: http",
		},

		{
			description: "successfully translates from v1.Pod to workload_model.V1Workload with TerminationMessagePolicy FALLBACK_TO_LOGS_ON_ERROR",
			pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Image:   "test-image",
							Command: []string{"/sh/bash"},
							Name:    "nginx",
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
								},
							},
							TerminationMessagePath:   "/test/path",
							TerminationMessagePolicy: v1.TerminationMessageFallbackToLogsOnError,
						},
					},
				},
			},
			expected: &workload_models.V1Workload{
				Name: "namespace-name",
				Slug: "namespace-name",
				Metadata: &workload_models.V1Metadata{
					Labels: map[string]string{
						podNameLabelKey:      "name",
						podNamespaceLabelKey: "namespace",
						nodeNameLabelKey:     "vk-mock",
					},
				},
				Spec: &workload_models.V1WorkloadSpec{
					Containers: workload_models.V1ContainerSpecMapEntry{
						"nginx": workload_models.V1ContainerSpec{
							Env:     workload_models.V1EnvironmentVariableMapEntry{},
							Image:   "test-image",
							Command: []string{"/sh/bash"},
							Ports: workload_models.V1InstancePortMapEntry{"http": workload_models.V1InstancePort{
								Port:                        8080,
								EnableImplicitNetworkPolicy: false,
								Protocol:                    "TCP",
							}},
							Resources: &workload_models.V1ResourceRequirements{
								Limits: workload_models.V1StringMapEntry{
									"cpu":    "1",
									"memory": "2Gi",
								},
								Requests: workload_models.V1StringMapEntry{
									"cpu":    "1",
									"memory": "2Gi",
								},
							},
							VolumeMounts:             []*workload_models.V1InstanceVolumeMount{},
							TerminationMessagePath:   "/test/path",
							TerminationMessagePolicy: workload_models.V1ContainerTerminationMessagePolicyFALLBACKTOLOGSONERROR.Pointer(),
						},
					},

					ImagePullCredentials: workload_models.V1WrappedImagePullCredentials{},
					NetworkInterfaces: []*workload_models.V1NetworkInterface{
						{
							EnableOneToOneNat: true,
							IPFamilies:        []*workload_models.V1IPFamily{workload_models.NewV1IPFamily("IPv4")},
							IPV6Subnet:        "",
							Network:           "default",
							Subnet:            "",
						},
					},
					VolumeClaimTemplates: []*workload_models.V1VolumeClaim{},
					Runtime:              &workload_models.V1WorkloadInstanceRuntimeSettings{},
				},
				Targets: workload_models.V1TargetMapEntry{
					"city-code": workload_models.V1Target{
						Spec: &workload_models.V1TargetSpec{
							DeploymentScope: "cityCode",
							Deployments: &workload_models.V1DeploymentSpec{
								MaxReplicas: 1,
								MinReplicas: 1,
								Selectors: []*workload_models.V1MatchExpression{
									{
										Key:      "cityCode",
										Operator: "in",
										Values:   []string{testCityCode},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		workload, err := provider.getWorkloadFrom(&test.pod)
		if err != nil {
			assert.ErrorContains(t, err, test.err, test.description)
		} else {
			assert.Equal(t, test.expected, workload, test.description)
		}
	}
}

func TestWorkloadRuntimeSettings(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	var testInt64 int64 = 60
	var testString string = "value"
	var testBool = true

	var tests = []struct {
		description string
		pod         v1.Pod
		expected    *workload_models.V1WorkloadInstanceRuntimeSettings
		err         error
		len         int
	}{
		{
			description: "empty pod return ok without runtime settings",
			pod:         v1.Pod{},
			expected:    &workload_models.V1WorkloadInstanceRuntimeSettings{},
		},
		{
			description: "successfully translate all runtime settings",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					TerminationGracePeriodSeconds: &testInt64,
					HostAliases: []v1.HostAlias{{
						IP:        "127.0.0.1",
						Hostnames: []string{"Hostname1", "Hostname2"},
					}},
					DNSConfig: &v1.PodDNSConfig{
						Options: []v1.PodDNSConfigOption{{
							Name:  "Test",
							Value: &testString,
						}},
						Nameservers: []string{"Nameserver1", "Nameserver2"},
						Searches:    []string{"Search1", "Search2"},
					},
					ShareProcessNamespace: &testBool,
				},
			},
			expected: &workload_models.V1WorkloadInstanceRuntimeSettings{
				Containers: &workload_models.V1WorkloadInstanceContainerRuntimeSettings{
					TerminationGracePeriodSeconds: strconv.Itoa(int(testInt64)),
					HostAliases: []*workload_models.V1HostAlias{
						{
							IP:        "127.0.0.1",
							Hostnames: []string{"Hostname1", "Hostname2"},
						},
					},
					DNSConfig: &workload_models.V1DNSConfig{
						Nameservers: []string{"Nameserver1", "Nameserver2"},
						Searches:    []string{"Search1", "Search2"},
						Options: []*workload_models.V1DNSConfigOption{
							{
								Name:  "Test",
								Value: testString,
							},
						},
					},
					ShareProcessNamespace: testBool,
				},
			},
		},
		{
			description: "successfully translate share process name runtime setting ",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					ShareProcessNamespace: &testBool,
				},
			},
			expected: &workload_models.V1WorkloadInstanceRuntimeSettings{
				Containers: &workload_models.V1WorkloadInstanceContainerRuntimeSettings{
					ShareProcessNamespace: testBool,
				},
			},
		},
		{
			description: "successfully translate DNS config runtime setting ",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					DNSConfig: &v1.PodDNSConfig{
						Options: []v1.PodDNSConfigOption{{
							Name:  "Test",
							Value: &testString,
						}},
						Nameservers: []string{"Nameserver1", "Nameserver2"},
						Searches:    []string{"Search1", "Search2"},
					},
				},
			},
			expected: &workload_models.V1WorkloadInstanceRuntimeSettings{
				Containers: &workload_models.V1WorkloadInstanceContainerRuntimeSettings{
					DNSConfig: &workload_models.V1DNSConfig{
						Nameservers: []string{"Nameserver1", "Nameserver2"},
						Searches:    []string{"Search1", "Search2"},
						Options: []*workload_models.V1DNSConfigOption{
							{
								Name:  "Test",
								Value: testString,
							},
						},
					},
				},
			},
		},
		{
			description: "successfully translate host aliases runtime setting ",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					HostAliases: []v1.HostAlias{{
						IP:        "127.0.0.1",
						Hostnames: []string{"Hostname1", "Hostname2"},
					}},
				},
			},
			expected: &workload_models.V1WorkloadInstanceRuntimeSettings{
				Containers: &workload_models.V1WorkloadInstanceContainerRuntimeSettings{
					HostAliases: []*workload_models.V1HostAlias{
						{
							IP:        "127.0.0.1",
							Hostnames: []string{"Hostname1", "Hostname2"},
						},
					}},
			},
		},
		{
			description: "successfully translate termination grace period runtime setting",
			pod: v1.Pod{
				Spec: v1.PodSpec{
					TerminationGracePeriodSeconds: &testInt64,
				},
			},
			expected: &workload_models.V1WorkloadInstanceRuntimeSettings{
				Containers: &workload_models.V1WorkloadInstanceContainerRuntimeSettings{
					TerminationGracePeriodSeconds: strconv.Itoa(int(testInt64)),
				},
			},
		},
	}

	for _, c := range tests {
		t.Run(c.description, func(t *testing.T) {
			runtime := provider.getWorkloadRuntimeSettingsFrom(c.pod.Spec)
			assert.Equal(t, c.expected, runtime)
		})
	}
}

func TestGetWorkloadContainerImagePullPolicyFrom(t *testing.T) {
	ctx := context.Background()

	provider, err := createTestProvider(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	var tests = []struct {
		description string
		pullPolicy  v1.PullPolicy
		expected    *workload_models.V1ContainerImagePullPolicy
	}{
		{
			description: "successfully translate pull always policy",
			pullPolicy:  v1.PullAlways,
			expected:    workload_models.V1ContainerImagePullPolicyALWAYS.Pointer(),
		},
		{
			description: "successfully translate pull not presented policy",
			pullPolicy:  v1.PullIfNotPresent,
			expected:    workload_models.V1ContainerImagePullPolicyIFNOTPRESENT.Pointer(),
		},
		{
			description: "returns nil for not supported image pull policy",
			pullPolicy:  v1.PullNever,
			expected:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			policy := provider.getWorkloadContainerImagePullPolicyFrom(test.pullPolicy)
			assert.Equal(t, test.expected, policy)
		})
	}
}

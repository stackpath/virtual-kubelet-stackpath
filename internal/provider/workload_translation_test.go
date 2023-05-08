package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_models"
	"github.com/stackpath/vk-stackpath-provider/internal/mocks"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
		workloadHTTP := provider.getProbeHTTPHeadersFrom(test.k8sHTTPHeaders)
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
	}

	_, err = provider.getWorkloadFrom(&pod)
	assert.ErrorContains(t, err, "unable to find named port")
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
		}
	}
}

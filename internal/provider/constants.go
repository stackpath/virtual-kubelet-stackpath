// Package provider implements the stackpath virtual kubelet provider
package provider

import "github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_models"

var (
	containerResourcesSP1 = workload_models.V1StringMapEntry{"cpu": "1", "memory": "2Gi"}
	containerResourcesSP2 = workload_models.V1StringMapEntry{"cpu": "2", "memory": "4Gi"}
	containerResourcesSP3 = workload_models.V1StringMapEntry{"cpu": "2", "memory": "8Gi"}
	containerResourcesSP4 = workload_models.V1StringMapEntry{"cpu": "4", "memory": "16Gi"}
	containerResourcesSP5 = workload_models.V1StringMapEntry{"cpu": "8", "memory": "32Gi"}
)

const defaultK8sServiceAccountMountPath = "/var/run/secrets/kubernetes.io/serviceaccount"
const stackpathVirtualKubeletCSIDriver = "virtual-kubelet.storage.compute.edgeengine.io"

var k8sEnvVariablesToIgnore = []string{
	"KUBERNETES_PORT",
	"KUBERNETES_PORT_443_TCP",
	"KUBERNETES_PORT_443_TCP_ADDR",
	"KUBERNETES_PORT_443_TCP_PORT",
	"KUBERNETES_PORT_443_TCP_PROTO",
	"KUBERNETES_SERVICE_HOST",
	"KUBERNETES_SERVICE_PORT",
	"KUBERNETES_SERVICE_PORT_HTTPS",
}

package provider

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	defaultCPUCoresNumber  = "10000"
	defaultMemorySize      = "1Ti"
	defaultStorageSize     = "1Ti"
	defaultPodsLimit       = "1000"
	defaultOperatingSystem = "Linux"
)

func (p *StackpathProvider) ConfigureNode(ctx context.Context, node *v1.Node) {
	node.Status.Capacity = p.getNodeCapacity()
	node.Status.Allocatable = p.getNodeCapacity()
	node.Status.NodeInfo.OperatingSystem = p.operatingSystem
	// node.Status.Addresses = p.getNodeAddresses()
	// node.Status.DaemonEndpoints = p.getNodeDaemonEndpoints()
}

func (p *StackpathProvider) getNodeCapacity() v1.ResourceList {
	resourceList := v1.ResourceList{
		v1.ResourceCPU:     resource.MustParse(p.cpu),
		v1.ResourceMemory:  resource.MustParse(p.memory),
		v1.ResourcePods:    resource.MustParse(p.pods),
		v1.ResourceStorage: resource.MustParse(p.storage),
	}

	return resourceList
}

// getNodeAddresses returns a list of addresses for the node status
// within Kubernetes.
func (p *StackpathProvider) getNodeAddresses() []v1.NodeAddress {
	return []v1.NodeAddress{
		{
			Type:    "InternalIP",
			Address: p.internalIP,
		},
	}
}

// getNodeDaemonEndpoints returns NodeDaemonEndpoints for the node status
// within Kubernetes.
func (p *StackpathProvider) getNodeDaemonEndpoints() v1.NodeDaemonEndpoints {
	return v1.NodeDaemonEndpoints{
		KubeletEndpoint: v1.DaemonEndpoint{
			Port: p.daemonEndpointPort,
		},
	}
}

func (p *StackpathProvider) setNodeCapacity() {
	p.cpu = defaultCPUCoresNumber
	p.memory = defaultMemorySize
	p.pods = defaultPodsLimit
	p.storage = defaultStorageSize
	p.operatingSystem = defaultOperatingSystem
}

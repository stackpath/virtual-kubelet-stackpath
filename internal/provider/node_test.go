package provider

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestConfigureNode(t *testing.T) {
	ctx := context.TODO()

	p, err := createTestProvider(ctx, nil, nil, nil, nil)
	p.cpu = "1"
	p.memory = "1Gi"
	p.pods = "10"
	p.storage = "1Gi"
	p.operatingSystem = "Linux"

	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	node := v1.Node{}
	p.ConfigureNode(ctx, &node)
	assert.Equal(t, strconv.FormatInt(node.Status.Capacity.Cpu().Value(), 10), p.cpu)
	assert.Equal(t, node.Status.Capacity.Memory().String(), p.memory)
	assert.Equal(t, node.Status.Capacity.Storage().String(), p.storage)
	assert.Equal(t, strconv.FormatInt(node.Status.Capacity.Pods().Value(), 10), p.pods)
	assert.Equal(t, node.Status.NodeInfo.OperatingSystem, p.operatingSystem)
}

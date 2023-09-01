package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_client"
	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_client/workload"
	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
	"github.com/stackpath/virtual-kubelet-stackpath/internal/auth"
	"github.com/stackpath/virtual-kubelet-stackpath/internal/config"
	"github.com/virtual-kubelet/virtual-kubelet/log"
	"gotest.tools/assert"
)

func TestPodLifecycle(t *testing.T) {
	// Define constants for namespace, pod name, and container name.
	const (
		namespace        = "vk-test"
		podName          = "pod-test"
		containerName    = "webserver-test"
		timeoutSeconds   = 240
		expectedWorkload = `{"metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{\"workload.platform.stackpath.net/remote-management\":\"true\"},\"name\":\"pod-test\",\"namespace\":\"vk-test\"},\"spec\":{\"containers\":[{\"env\":[{\"name\":\"VAR\",\"value\":\"val\"}],\"image\":\"nginxinc/nginx-unprivileged:latest\",\"imagePullPolicy\":\"Always\",\"lifecycle\":{\"preStop\":{\"httpGet\":{\"path\":\"shutdown\",\"port\":8080}}},\"livenessProbe\":{\"exec\":{\"command\":[\"uname\"]},\"initialDelaySeconds\":5,\"periodSeconds\":10},\"name\":\"webserver-test\",\"ports\":[{\"containerPort\":80,\"name\":\"http\"},{\"containerPort\":443,\"name\":\"https\"}],\"readinessProbe\":{\"failureThreshold\":1,\"httpGet\":{\"httpHeaders\":[{\"name\":\"Custom-Header\",\"value\":\"Awesome\"}],\"path\":\"/\",\"port\":8080},\"initialDelaySeconds\":5,\"periodSeconds\":10,\"successThreshold\":2,\"timeoutSeconds\":10},\"resources\":{\"limits\":{\"cpu\":\"2\",\"memory\":\"4Gi\"},\"requests\":{\"cpu\":\"250m\",\"memory\":\"1Gi\"}},\"securityContext\":{\"allowPrivilegeEscalation\":true,\"capabilities\":{\"add\":[\"NET_ADMIN\"],\"drop\":[\"LINUX_IMMUTABLE\"]},\"runAsGroup\":3000,\"runAsNonRoot\":true,\"runAsUser\":2000},\"startupProbe\":{\"failureThreshold\":1,\"initialDelaySeconds\":1,\"periodSeconds\":5,\"successThreshold\":1,\"tcpSocket\":{\"port\":8080},\"timeoutSeconds\":1},\"terminationMessagePath\":\"/var/log/nginx/error.log\",\"terminationMessagePolicy\":\"File\",\"volumeMounts\":[{\"mountPath\":\"/disk-1\",\"name\":\"volume-1\"}],\"workingDir\":\"/\"}],\"dnsConfig\":{\"nameservers\":[\"192.0.2.1\"],\"options\":[{\"name\":\"ndots\",\"value\":\"2\"}],\"searches\":[\"ns1.svc.cluster-domain.example\",\"my.dns.search.suffix\"]},\"hostAliases\":[{\"hostnames\":[\"foo.local\",\"bar.local\"],\"ip\":\"127.0.0.1\"},{\"hostnames\":[\"foo.remote\",\"bar.remote\"],\"ip\":\"10.1.2.3\"}],\"initContainers\":[{\"args\":[\"echo hello there!;sleep 1;\"],\"command\":[\"sh\",\"-c\"],\"image\":\"busybox:1.28\",\"name\":\"init-container\"}],\"nodeSelector\":{\"kubernetes.io/role\":\"agent\",\"type\":\"virtual-kubelet\"},\"securityContext\":{\"runAsGroup\":3000,\"runAsNonRoot\":true,\"runAsUser\":2000,\"supplementalGroups\":[100,200],\"sysctls\":[{\"name\":\"net.core.somaxconn\",\"value\":\"1024\"}]},\"shareProcessNamespace\":true,\"tolerations\":[{\"effect\":\"NoSchedule\",\"key\":\"virtual-kubelet.io/provider\",\"operator\":\"Equal\",\"value\":\"stackpath\"}],\"volumes\":[{\"csi\":{\"driver\":\"virtual-kubelet.storage.compute.edgeengine.io\",\"volumeAttributes\":{\"size\":\"2Gi\"}},\"name\":\"volume-1\"}]}}\n","workload.platform.stackpath.net/remote-management":"true"},"labels":{"vk-node-name":"vk-sp-test-node-dfw","vk-pod-name":"pod-test","vk-pod-namespace":"vk-test"},"version":"1"},"name":"vk-test-pod-test","slug":"vk-test-pod-test","spec":{"containers":{"webserver-test":{"args":[],"command":[],"env":{"VAR":{"value":"val"}},"image":"nginxinc/nginx-unprivileged:latest","imagePullPolicy":"ALWAYS","lifecycle":{"preStop":{"httpGet":{"path":"shutdown","port":8080,"scheme":"HTTP"}}},"livenessProbe":{"exec":{"command":["uname"]},"failureThreshold":3,"initialDelaySeconds":5,"periodSeconds":10,"successThreshold":1,"timeoutSeconds":1},"ports":{"http":{"port":80,"protocol":"TCP"},"https":{"port":443,"protocol":"TCP"}},"readinessProbe":{"failureThreshold":1,"httpGet":{"httpHeaders":{"Custom-Header":"Awesome"},"path":"/","port":8080,"scheme":"HTTP"},"initialDelaySeconds":5,"periodSeconds":10,"successThreshold":2,"timeoutSeconds":10},"resources":{"limits":{"cpu":"2","memory":"4Gi"},"requests":{"cpu":"2","memory":"4Gi"}},"securityContext":{"allowPrivilegeEscalation":true,"capabilities":{"add":["NET_ADMIN"],"drop":["LINUX_IMMUTABLE"]},"runAsGroup":"3000","runAsNonRoot":true,"runAsUser":"2000"},"startupProbe":{"failureThreshold":1,"initialDelaySeconds":1,"periodSeconds":5,"successThreshold":1,"tcpSocket":{"port":8080},"timeoutSeconds":1},"terminationMessagePath":"/var/log/nginx/error.log","terminationMessagePolicy":"FILE","volumeMounts":[{"mountPath":"/disk-1","slug":"volume-1"}],"workingDir":"/"}},"initContainers":{"init-container":{"args":["echo hello there!;sleep 1;"],"command":["sh","-c"],"image":"busybox:1.28","imagePullPolicy":"IF_NOT_PRESENT","resources":{"limits":{"cpu":"1","memory":"2Gi"},"requests":{"cpu":"1","memory":"2Gi"}},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"FILE","volumeMounts":[]}},"networkInterfaces":[{"enableOneToOneNat":true,"ipFamilies":["IPv4"],"network":"default"}],"runtime":{"containers":{"dnsConfig":{"nameservers":["192.0.2.1"],"options":[{"name":"ndots","value":"2"}],"searches":["ns1.svc.cluster-domain.example","my.dns.search.suffix"]},"hostAliases":[{"hostnames":["foo.local","bar.local"],"ip":"127.0.0.1"},{"hostnames":["foo.remote","bar.remote"],"ip":"10.1.2.3"}],"securityContext":{"runAsGroup":"3000","runAsNonRoot":true,"runAsUser":"2000","supplementalGroups":["100","200"],"sysctls":[{"name":"net.core.somaxconn","value":"1024"}]},"shareProcessNamespace":true,"terminationGracePeriodSeconds":"30"}},"volumeClaimTemplates":[{"metadata":{},"name":"volume-1","phase":"PENDING","slug":"volume-1","spec":{"resources":{"limits":{"storage":"2Gi"},"requests":{"storage":"2Gi"}}}}]},"status":"ACTIVE","targets":{"city-code":{"spec":{"deploymentScope":"cityCode","deployments":{"maxReplicas":1,"minReplicas":1,"selectors":[{"key":"cityCode","operator":"in","values":["DFW"]}]}}}}}`
	)

	ctx := context.Background()

	apiConfig, err := config.NewConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}

	stackpathClient := createSPClient(ctx, t, apiConfig)

	// Clean up 'vk-test' namespace from previous tests if it exists.
	t.Log("Cleaning up pods from previous tests if it exist...")
	runKubectlCommand(t, "delete", "namespace", namespace, "--ignore-not-found")

	// Create 'vk-test' namespace for the test.
	t.Log("Creating 'vk-test' namespace for the test...")
	runKubectlCommand(t, "apply", "-f", "fixtures/namespace.yml")

	// Apply the podspec file to create the pod.
	t.Log("Applying the podspec file...")
	runKubectlCommand(t, "apply", "-f", "fixtures/test_pod_with_all_supported_features.yml", "--namespace="+namespace)

	// Wait for the pod to be created and ready.
	t.Log("Waiting for the pod to be created...")
	waitForPodReady(t, podName, namespace, timeoutSeconds)
	t.Log("Successfully created the pod")

	// Get the pod's status.
	t.Log("Getting pod status...")
	podStatus := getPodStatus(t, podName, namespace)
	t.Logf("Got the pod's status: %s", podStatus)

	// Check the container's status.
	t.Log("Getting container's status...")
	containerReady := isContainerReady(t, podName, namespace, containerName)
	t.Logf("Got the container's status: %v", containerReady)

	t.Log("Getting workload data by calling SP// API...")
	workloadInfo := getWorkloadInfo(ctx, t, apiConfig, stackpathClient, namespace, podName)
	t.Log("Got the workload data")

	// Cleaning up changing fields
	workloadInfo.ID = ""
	workloadInfo.StackID = ""
	workloadInfo.Metadata.CreatedAt = nil
	// Marshal the workload struct into JSON format
	jsonBytes, err := json.Marshal(workloadInfo)
	if err != nil {
		t.Fatalf("Error marshaling struct: %v", err)
	}
	// Convert the JSON bytes to a string
	workloadString := string(jsonBytes)
	assert.Equal(t, workloadString, expectedWorkload)

	// Get the container's logs.
	t.Log("Getting the container logs...")
	containerLogs := getContainerLogs(t, podName, namespace, containerName, 5)
	t.Logf("Container logs: %s", containerLogs)

	// Check container exec.
	t.Log("Checking execute commands on container...")
	execOutput := execOnContainer(t, podName, namespace, containerName, "--", "/bin/sh", "-c", "ls")
	if strings.Contains(execOutput, "index.php") {
		t.Fatal("Failed to execute on the container")
	}
	t.Log("Successfully executed commands on the container")

	// Clean up.
	t.Log("Cleaning up...")
	runKubectlCommand(t, "delete", "namespace", namespace, "--ignore-not-found")
}

func runKubectlCommand(t *testing.T, args ...string) {
	cmd := kubectl(args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Error running kubectl %v: %s", args, string(out))
	}
}

func createSPClient(ctx context.Context, t *testing.T, apiConfig *config.Config) *workload_client.EdgeCompute {
	runtime, err := auth.NewRuntime(ctx, apiConfig.ClientID, apiConfig.ClientSecret, apiConfig.ApiHost, "version-test")
	if err != nil {
		log.G(ctx).Fatal(err)
	}
	// Create StackPath client
	return workload_client.New(runtime, nil)
}

func getWorkloadInfo(ctx context.Context, t *testing.T, apiConfig *config.Config, client *workload_client.EdgeCompute, namespace string, podName string) *workload_models.V1Workload {

	getWorkloadParams := workload.GetWorkloadParams{
		Context: ctx,
		StackID: apiConfig.StackID,
		WorkloadID: strings.ToLower(
			fmt.Sprintf("%s-%s", namespace, podName),
		),
	}

	workloadResult, err := client.Workload.GetWorkload(&getWorkloadParams, nil)
	if err != nil {
		t.Fatal(err)
	}
	return workloadResult.Payload.Workload
}

func waitForPodReady(t *testing.T, podName, namespace string, timeoutSeconds int) {
	deadline, ok := t.Deadline()
	timeout := time.Until(deadline)
	if !ok {
		timeout = time.Duration(timeoutSeconds) * time.Second
	}
	cmd := kubectl("wait", "--for=condition=ready", "--timeout="+timeout.String(), "pod/"+podName, "--namespace="+namespace)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Error waiting for pod to be ready: %s", string(out))
	}
}

func getPodStatus(t *testing.T, podName, namespace string) string {
	cmd := kubectl("get", "pod", "--field-selector=status.phase=Running", "--namespace="+namespace, "--output=jsonpath={.items..metadata.name}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(out))
	}
	return string(out)
}

func isContainerReady(t *testing.T, podName, namespace, containerName string) bool {
	cmd := kubectl("get", "pod", podName, "--namespace="+namespace, "--output=jsonpath={.status.containerStatuses[0].ready}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(out))
	}
	return string(out) == "true"
}

func getContainerLogs(t *testing.T, podName, namespace, containerName string, tailLines int) string {
	cmd := kubectl("logs", "pod/"+podName, "-c", containerName, "--namespace="+namespace, fmt.Sprintf("--tail=%d", tailLines))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(out))
	}
	return string(out)
}

func execOnContainer(t *testing.T, podName, namespace, containerName string, command ...string) string {
	cmdArgs := append([]string{"exec", "pod/" + podName, "-c", containerName, "--namespace=" + namespace}, command...)
	cmd := kubectl(cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(out))
	}
	return string(out)
}

func kubectl(args ...string) *exec.Cmd {
	// Define the kubectl command.
	cmd := exec.Command("kubectl", args...)
	return cmd
}

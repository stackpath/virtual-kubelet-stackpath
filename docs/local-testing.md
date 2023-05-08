
# Testing Virtual Kubelet for StackPath Edge Compute platform

This guide will walk you through the steps to test the  Virtual Kubelet for StackPath Edge Compute Platform in a Kubernetes cluster.

## Prerequisites
 - [Docker](https://www.docker.com/) 
 - [Kubernetes](https://kubernetes.io/)
 - [kubectl](https://kubernetes.io/docs/tasks/tools/)
 - [Kubernetes Metrics Server](https://github.com/kubernetes-sigs/metrics-server)

## Building the Docker Image
To build the Docker image containing the Virtual Kubelet with StackPath Edge Compute provider, simply execute the following command: 
```
make build-image
```
This command builds a new `stackpath.com/virtual-kubelet` image, complete with the StackPath Edge Compute provider integration.

## Kubernetes Metrics Server
Metrics Server is a scalable, efficient source of container resource metrics for Kubernetes built-in autoscaling pipelines.

Metrics Server collects resource metrics from Kubelets and exposes them in Kubernetes apiserver through Metrics API for use by [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) and [Vertical Pod Autoscaler](https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler/). Metrics API can also be accessed by `kubectl top`, making it easier to debug autoscaling pipelines.

### Installing the Metrics Server
The Kubernetes Metrics Server is not installed by default in your Kubernetes (k8s) cluster. Follow these steps to install it and ensure it's running properly:
1. To install the Metrics Server, execute the following command:
```
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```
This command creates the necessary components for the k8s Metrics Server.

2. To verify that the Metrics Server is running, use this command:
```
kubectl get pods -n kube-system| grep metric
```
You should see something like: 
```
metrics-server-6c6cdff5d5-7hsf8          1/1     Running   0                13m
```
This indicates that the Metrics Server is running smoothly.

If you see output like:
```
v1beta1.metrics.k8s.io                 kube-system/metrics-server   False (MissingEndpoints)   90m
``` 
You may need to update the configuration for the service. Perform the following steps to fix potential issues:

1. Execute `kubectl edit deployments.apps -n kube-system metrics-server`. This will open an editor with the Metrics Server configuration.
2. Scroll down until you see the `spec` object. Add the following options:
```
--kubelet-insecure-tls=true
--requestheader-client-ca-file=
``` 
Save your edits. In the terminal, you should see the message:
On the terminal, you will see the message:
```
deployment.apps/metrics-server edited
```

With the Metrics Server installed and running, you can now monitor cluster resource usage. To view the resources of the nodes, run:
```
kubectl top nodes
```
To monitor the resources of the pods, run:
```
kubectl top pods
```

The `top` command provides real-time resource usage information. However, it does not display historical data, such as your resource usage from yesterday or a week ago.


## Deployment
TBD

## Automated e2e tests
TBD
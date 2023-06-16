# Virtual Kubelet Provider for StackPath Edge Compute

StackPath's Virtual Kubelet (VK) provider allows you to leverage the power of Kubernetes (K8s) to deploy and manage your applications across StackPath's expansive Edge Compute network, increasing scalability and reliability, while decreasing latency.

This feature enables you to use the Kubernetes control plane to create and manage pods as you normally would, without having to worry about managing your own hardware and infrastructure, as StackPath's Virtual Kubelet provider takes care of scheduling these pods for you on our Edge Compute nodes.

## Key Features

- **Volumes using `csi`**. Mount volumes in your pods using the `csi` volume type with the driver `virtual-kubelet.storage.compute.edgeengine.io`.
- **Environment variables**. Set environment variables for your pods using the Kubernetes `env` field in your pod specification.
- **Instance size selection**. Specify resource requirements for your pods using the Kubernetes `resources` field in your pod specification.
- **Liveness and readiness probes**. Configure `liveness` and `readiness` probes for your pods using - the Kubernetes livenessProbe and readinessProbe fields in your pod specification.
- **Private images using image pull secrets**. Use Kubernetes image pull secrets to securely pull private container images from a registry using the Kubernetes `imagePullSecrets` field in your pod specification.
- **Run in container**. Run commands inside your pods by calling `kubectl exec` and specifying the pod and command.

## Limitations

- **Limited instance types**. StackPath Edge Compute product currently supports five instance types: SP-1 through SP-5 ([SP// Containers](https://www.stackpath.com/products/containers/)). The provider will launch the smallest instance that provides the resources defined in the pod specification. If the pod specification requires more resources than what is available in the SP-5 instance, the provider will provision the SP-5 instance type.
Here are the specifications for each of the available instance types:

    | Subscription | Cores | RAM  |
    | ---  |---| ---  |
    | SP-1 | 1 | 2GB  |
    | SP-2 | 2 | 4GB  |
    | SP-3 | 2 | 8GB  |
    | SP-4 | 4 | 16GB |
    | SP-5 | 8 | 32GB |

- **Limited network control**. The provider currently does not support custom network settings for the StackPath workload. The workload will run with a public and private IP, and network policies must be created separately.
- **Limited probe support**. The provider currently only supports the `httpGet` and `tcpSocket` probes for liveness and readiness checks. Other probe types, such as `grpc` or `exec`, are not currently supported.
- **Limited Kubernetes features**. The provider only supports some of the Kubernetes pod specification as supported by the StackPath edge compute platform, there may be some advanced features that are not yet supported or that require additional configuration. **The provider will ignore any specification that aren't supported when creating the StackPath workload**.
In addition, the workloads created on the StackPath platform will not have network access to the Kubernetes API or any pods running in nodes that aren't in the virtual kubelet provider.
- **Pod name length**. The provider is subject to the limitations of StackPath's workload slugs, which are limited to 63 characters. The provider constructs the slug by concatenating the namespace with the pod name separated by a dash. It is important to ensure that this string does not exceed 63 characters, as exceeding this limit will prevent the pod from being created.

## Getting Started

To use the StackPath Edge Computing Virtual Kubelet provider, you'll need to have the following prerequisites:

- A Kubernetes cluster
- A StackPath account
- [Credentials](https://stackpath.dev/docs/stackpath-api-authentication#api-credentials) for your StackPath account (a client ID and a client secret)

## Installing the Provider

See [deployment details](./deployment/README.md)

## Enabling Remote Management for Pods

To enable remote management for pods, you can use the `workload.platform.stackpath.net/remote-management` annotation in the pod definition metadata. By setting this annotation to `true`, the remote management capabilities for the containers listed in the pod will be enabled. 

To enable remote management, add the following annotation to your pod definition metadata:

```yaml
annotations:
  workload.platform.stackpath.net/remote-management: "true"
```
By default, if this annotation is not provided or set to "false", remote management will be disabled.

For more information about Edge Compute Workload Metadata and other terms related to StackPath Edge Compute, refer to the [Learn Edge Compute Terms](https://support.stackpath.com/hc/en-us/articles/360059500391-Learn-Edge-Compute-Terms) page.

It is important to note that enabling remote management should be done with caution and only for trusted pods or in controlled environments where appropriate security measures are in place.

## Pod Spec File Examples

We've included some sample Pod spec files in the [tests/e2e](tests/e2e/) folder to help you get started with. These examples are intended to demonstrate how to configure the virtual kubelet to work with different types of StackPath workloads.

To use one of these examples, simply copy the YAML file and modify it to suit your needs. You can then use `kubectl` to create the Pod in your Kubernetes cluster.

Note that these examples are provided for demonstration purposes only and may not be suitable for production use. Be sure to review and modify the configuration settings to suit your specific needs.

## Conclusion

The Virtual Kubelet Provider for StackPath Edge Compute is a useful tool for deploying Kubernetes pods on StackPath's highly performant and distributed edge infrastructure. With this provider, you can leverage the power of Kubernetes to deploy and manage your applications, while taking advantage of StackPath's Edge Compute capabilities.

While the provider has some limitations, it is a great option for workloads that need to be close to end-users and require low latency or high bandwidth. By seamlessly extending your Kubernetes clusters to the edge with the Virtual Kubelet Provider for StackPath Edge Compute, you can improve the performance and reliability of your applications.

## Contributing

If you encounter any issues or have ideas for improvements, don't hesitate to contribute to the project on GitHub. The StackPath team welcomes and values community feedback and contributions.

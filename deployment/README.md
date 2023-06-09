# Virtual Kubelet Deployment README

This README file provides instructions on how to deploy a Kubernetes deployment for StackPath Virtual Kubelet Provider using Kustomize.

## Usage

1. Confirm that Kustomize is installed in your environment by running the `kustomize version` command. If you haven't already installed Kustomize, follow the installation instructions [here](https://kubectl.docs.kubernetes.io/installation/kustomize/).

To find the kustomize version embedded in recent versions of kubectl, run kubectl version:

```bash
> kubectl version --short --client
Client Version: v1.26.0
Kustomize Version: v4.5.7
```

2. Clone this repository to your local environment.

3. Navigate to the `base` directory, which contains the base Virtual Kubelet deployment.

```bash
cd deployment/kustomize/base
```

4. Follow this [guide](https://stackpath.dev/docs/stackpath-api-quick-start#api-credentials) to obtain StackPath API credentials and update the `config.properties` file with your StackPath account, stack, client, and secret IDs:

```txt
SP_STACK_ID=<your-stack-id>
SP_CLIENT_ID=<your-client-id>
SP_CLIENT_SECRET=<your-client-secret>
```

5. To deploy the Virtual Kubelet resources, run the following command:

```bash
kubectl apply -k .
```

This will create the Virtual Kubelet deployment in your Kubernetes cluster.

Note that a secret will be generated from the `config.properties` file specified in the `secretGenerator` section of the `kustomization.yaml` file. This secret contains the values of the environment variables specified in the `config.properties` file.

## Updating Resources

To customize the Virtual Kublet deployment, create an overlay directory (in this example `vk-deployment-updated`) within the `overlays` directory with a `kustomization.yaml` file that specifies the changes you want to make.

```txt
.
├── base
│   ├── cluster-role.yaml
│   ├── config.properties
│   ├── kustomization.yaml
│   ├── namespace.yaml
│   ├── service-account.yaml
│   └── vk-deployment.yaml
└── overlays
    └── vk-deployment-updated
        └── kustomization.yaml
```

For example, to create a Virtual Kubelet in a namespace other than the default one (in this example `sp-atl`) and update the values of `SP_CITY_CODE` and `SP_STACK_ID` environment variables, create the following `kustomization.yaml` file under the overlay directory:

```yaml
resources:
- ../../base

namespace: sp-atl

images:
- name: stackpath.com/virtual-kubelet
  newTag: 0.0.2

configMapGenerator:
- name: sp-vk-location
  behavior: replace
  literals:
    - SP_CITY_CODE=ATL

secretGenerator:
- name: sp-vk-secrets
  behavior: merge
  literals:
    - SP_STACK_ID=<another_stack_id>
```

- The resources section references the base resources that are inherited by this overlay, which includes a default Virtual Kubelet deployment configuration.
- The namespace section specifies that the Virtual Kubelet deployment will be created in the sp-atl namespace.
- The images section is used to define the version of the StackPath Virtual Kubelet image to be used.
- The configMapGenerator section replaces the existing value of SP_CITY_CODE with `ATL`, which specifies the geographic location of the edge compute infrastructure.
- The secretGenerator section merges the existing config.properties file with a new SP_STACK_ID value of <another_stack_id>. This updates the StackPath stack ID specified in `config.properties`.

To deploy overlay, run the following command:

```bash
kubectl apply -k overlays/vk-deployment-updated
```

## Configuring Pods to Use Virtual Kubelet

Now that you've created a Virtual Kubelet pod after completing the steps above, you're ready to move on to the next step. Once this pod is running, you can then create a standard pod and StackPath workload.

To use the Virtual Kubelet deployment to deploy workloads in the StackPath Edge Compute infrastructure, configure your pods to use the virtual-kubelet.io/provider toleration and type: virtual-kubelet node selector.

Here is an example configuration that will create the simplest possible container in the sp-atl namespace by providing only a name (my-pod) and image (my-image):

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
  namespace: sp-atl
spec:
  containers:
  - name: my-container
    image: my-image
  tolerations:
  - key: virtual-kubelet.io/provider
    operator: Equal
    value: stackpath
    effect: NoSchedule
  nodeSelector: 
    kubernetes.io/role: agent
    type: virtual-kubelet
```

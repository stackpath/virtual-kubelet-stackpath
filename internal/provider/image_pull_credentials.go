// Package provider implements the stackpath virtual kubelet provider
package provider

import (
	"encoding/json"
	"fmt"

	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_models"
	"github.com/virtual-kubelet/virtual-kubelet/errdefs"
	v1 "k8s.io/api/core/v1"
)

type dockerConfig map[string]dockerConfigEntry

type dockerConfigEntry struct {
	// +optional
	Username string `json:"username,omitempty"`
	// +optional
	Password string `json:"password,omitempty"`
	// +optional
	Email string `json:"email,omitempty"`
	// +optional
	Auth string `json:"auth,omitempty"`
}

type dockerConfigJSON struct {
	Auths dockerConfig `json:"auths"`
	// +optional
	HTTPHeaders map[string]string `json:"HttpHeaders,omitempty"`
}

func (p *StackpathProvider) getImagePullCredentialsFrom(namespace string, k8sImagePullSecrets []v1.LocalObjectReference) ([]*workload_models.V1ImagePullCredential, error) {
	imagePullCredentials := workload_models.V1WrappedImagePullCredentials{}
	// if there are image pull credentials, fetch the secret and fill the details.
	if len(k8sImagePullSecrets) != 0 {
		for _, k8sImagePullCredential := range k8sImagePullSecrets {
			imagePullCredential, err := p.getImagePullCredentialFrom(namespace, k8sImagePullCredential)
			if err != nil {
				if errdefs.IsNotFound(err) {
					// if the secret isn't found we don't want to fail, just skip it
					continue
				}
				return nil, err
			}
			imagePullCredentials = append(imagePullCredentials, imagePullCredential...)
		}
	}
	return imagePullCredentials, nil
}

func (p *StackpathProvider) getImagePullCredentialFrom(namespace string, k8sImagePullSecret v1.LocalObjectReference) ([]*workload_models.V1ImagePullCredential, error) {
	// get the secret from the secret lister
	imagePullCredentials := make([]*workload_models.V1ImagePullCredential, 0)
	k8sSecret, err := p.secretLister.Secrets(namespace).Get(k8sImagePullSecret.Name)
	if err != nil {
		if errdefs.IsNotFound(err) {
			p.logger.Warnf("image pull secret %s not found", k8sImagePullSecret.Name)
			return nil, err
		}
		p.logger.Errorf("error getting image pull secret %s", k8sImagePullSecret)
		return nil, err
	}

	if k8sSecret == nil {
		return nil, fmt.Errorf("error getting image pull secret %s", k8sImagePullSecret)
	}
	switch k8sSecret.Type {
	case v1.SecretTypeDockercfg:
		return nil, fmt.Errorf("legacy format kubernetes.io/dockercfg is not supported, please use kubernetes.io/dockerconfigjson")
	case v1.SecretTypeDockerConfigJson:
		ipcs, err := getImagePullCredentialsFromJSONConfig(k8sSecret)
		if err != nil {
			return nil, err
		}
		imagePullCredentials = append(imagePullCredentials, ipcs...)
	default:
		return nil, fmt.Errorf("image pull secret type is not kubernetes.io/dockerconfigjson")
	}

	return imagePullCredentials, nil
}

func getImagePullCredentialsFromJSONConfig(secret *v1.Secret) ([]*workload_models.V1ImagePullCredential, error) {
	imagePullCredentials := make([]*workload_models.V1ImagePullCredential, 0)
	repoDataJSON, ok := secret.Data[v1.DockerConfigJsonKey]
	if !ok {
		return nil, fmt.Errorf("no dockerconfigjson present in secret")
	}

	repoData := dockerConfigJSON{}
	err := json.Unmarshal(repoDataJSON, &repoData)
	if err != nil {
		return nil, err
	}

	for server, auth := range repoData.Auths {
		imagePullSecret := workload_models.V1ImagePullCredential{
			DockerRegistry: &workload_models.V1DockerRegistryCredentials{
				Server:   server,
				Username: auth.Username,
				Password: auth.Password,
				Email:    auth.Email,
			},
		}
		imagePullCredentials = append(imagePullCredentials, &imagePullSecret)
	}

	return imagePullCredentials, nil
}

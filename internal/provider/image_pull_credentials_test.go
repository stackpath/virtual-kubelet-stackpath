package provider

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_models"
	"github.com/stackpath/vk-stackpath-provider/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/virtual-kubelet/virtual-kubelet/errdefs"
	v1 "k8s.io/api/core/v1"
)

func TestImagePullCredentials(t *testing.T) {
	ctx := context.Background()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var tests = []struct {
		description         string
		err                 string
		namespace           string
		k8sImagePullSecrets []v1.LocalObjectReference
		expected            []*workload_models.V1ImagePullCredential
		secretMapListerMock func() *mocks.MockSecretLister
	}{
		{
			description:         "test empty image pull secret return empty image pull credentials",
			err:                 "",
			namespace:           "test",
			k8sImagePullSecrets: []v1.LocalObjectReference{},
			expected:            []*workload_models.V1ImagePullCredential{},
			secretMapListerMock: func() *mocks.MockSecretLister { return mocks.NewMockSecretLister(mockController) },
		},
		{
			description: "test valid credential work",
			err:         "",
			namespace:   "test",
			k8sImagePullSecrets: []v1.LocalObjectReference{
				{
					Name: "image-pull-secret",
				},
			},
			expected: []*workload_models.V1ImagePullCredential{
				{
					DockerRegistry: &workload_models.V1DockerRegistryCredentials{
						Username: "user",
						Password: "password",
						Email:    "user@gmail.com",
						Server:   "server",
					},
				},
			},
			secretMapListerMock: func() *mocks.MockSecretLister {
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
				secretListerMock := mocks.NewMockSecretLister(mockController)
				secretNamespaceListerMock := mocks.NewMockSecretNamespaceLister(gomock.NewController(t))
				secretNamespaceListerMock.EXPECT().Get("image-pull-secret").Return(
					&v1.Secret{
						Type: "kubernetes.io/dockerconfigjson",
						Data: map[string][]byte{
							".dockerconfigjson": secret,
						},
					},
					nil,
				)
				secretListerMock.EXPECT().Secrets("test").Return(secretNamespaceListerMock)
				return secretListerMock
			},
		},
		{
			description: "test invalid secret key name returns error",
			err:         "no dockerconfigjson present in secret",
			namespace:   "test",
			k8sImagePullSecrets: []v1.LocalObjectReference{
				{
					Name: "image-pull-secret",
				},
			},
			expected: []*workload_models.V1ImagePullCredential{
				{
					DockerRegistry: &workload_models.V1DockerRegistryCredentials{
						Username: "user",
						Password: "password",
						Email:    "user@gmail.com",
						Server:   "server",
					},
				},
			},
			secretMapListerMock: func() *mocks.MockSecretLister {
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
				secretListerMock := mocks.NewMockSecretLister(mockController)
				secretNamespaceListerMock := mocks.NewMockSecretNamespaceLister(gomock.NewController(t))
				secretNamespaceListerMock.EXPECT().Get("image-pull-secret").Return(
					&v1.Secret{
						Type: "kubernetes.io/dockerconfigjson",
						Data: map[string][]byte{
							".invalidkey": secret,
						},
					},
					nil,
				)
				secretListerMock.EXPECT().Secrets("test").Return(secretNamespaceListerMock)
				return secretListerMock
			},
		},
		{
			description: "test invalid docker config json",
			err:         "invalid character",
			namespace:   "test",
			k8sImagePullSecrets: []v1.LocalObjectReference{
				{
					Name: "image-pull-secret",
				},
			},
			expected: []*workload_models.V1ImagePullCredential{
				{
					DockerRegistry: &workload_models.V1DockerRegistryCredentials{
						Username: "user",
						Password: "password",
						Email:    "user@gmail.com",
						Server:   "server",
					},
				},
			},
			secretMapListerMock: func() *mocks.MockSecretLister {
				secretListerMock := mocks.NewMockSecretLister(mockController)
				secretNamespaceListerMock := mocks.NewMockSecretNamespaceLister(gomock.NewController(t))
				secretNamespaceListerMock.EXPECT().Get("image-pull-secret").Return(
					&v1.Secret{
						Type: "kubernetes.io/dockerconfigjson",
						Data: map[string][]byte{
							".dockerconfigjson": []byte("bla"),
						},
					},
					nil,
				)
				secretListerMock.EXPECT().Secrets("test").Return(secretNamespaceListerMock)
				return secretListerMock
			},
		},
		{
			description: "test secret not found skipped",
			err:         "",
			namespace:   "test",
			k8sImagePullSecrets: []v1.LocalObjectReference{
				{
					Name: "image-pull-secret-does-not-exist",
				},
			},
			expected: []*workload_models.V1ImagePullCredential{},
			secretMapListerMock: func() *mocks.MockSecretLister {
				secretListerMock := mocks.NewMockSecretLister(mockController)
				secretNamespaceListerMock := mocks.NewMockSecretNamespaceLister(gomock.NewController(t))
				secretNamespaceListerMock.EXPECT().Get("image-pull-secret-does-not-exist").Return(
					nil,
					errdefs.NotFound("not found"),
				)
				secretListerMock.EXPECT().Secrets("test").Return(secretNamespaceListerMock)
				return secretListerMock
			},
		},
		{
			description: "failed getting secret from lister returns error",
			err:         "lister failed",
			namespace:   "test",
			k8sImagePullSecrets: []v1.LocalObjectReference{
				{
					Name: "image-pull-secret-failed",
				},
			},
			expected: []*workload_models.V1ImagePullCredential{},
			secretMapListerMock: func() *mocks.MockSecretLister {
				secretListerMock := mocks.NewMockSecretLister(mockController)
				secretNamespaceListerMock := mocks.NewMockSecretNamespaceLister(gomock.NewController(t))
				secretNamespaceListerMock.EXPECT().Get("image-pull-secret-failed").Return(
					nil,
					errors.New("lister failed"),
				)
				secretListerMock.EXPECT().Secrets("test").Return(secretNamespaceListerMock)
				return secretListerMock
			},
		},
		{
			description: "lister returned invalid secret",
			err:         "error getting image pull secret",
			namespace:   "test",
			k8sImagePullSecrets: []v1.LocalObjectReference{
				{
					Name: "image-pull-secret-invalid",
				},
			},
			expected: []*workload_models.V1ImagePullCredential{},
			secretMapListerMock: func() *mocks.MockSecretLister {
				secretListerMock := mocks.NewMockSecretLister(mockController)
				secretNamespaceListerMock := mocks.NewMockSecretNamespaceLister(gomock.NewController(t))
				secretNamespaceListerMock.EXPECT().Get("image-pull-secret-invalid").Return(
					nil,
					nil,
				)
				secretListerMock.EXPECT().Secrets("test").Return(secretNamespaceListerMock)
				return secretListerMock
			},
		},
		{
			description: "legacy docker secret returns error",
			err:         "legacy format kubernetes.io/dockercfg is not supported",
			namespace:   "test",
			k8sImagePullSecrets: []v1.LocalObjectReference{
				{
					Name: "image-pull-secret",
				},
			},
			expected: []*workload_models.V1ImagePullCredential{},
			secretMapListerMock: func() *mocks.MockSecretLister {
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
				return secretListerMock
			},
		},
		{
			description: "other secret type returns error",
			err:         "image pull secret type is not kubernetes.io/dockerconfigjson",
			namespace:   "test",
			k8sImagePullSecrets: []v1.LocalObjectReference{
				{
					Name: "image-pull-secret",
				},
			},
			expected: []*workload_models.V1ImagePullCredential{},
			secretMapListerMock: func() *mocks.MockSecretLister {
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
						Type: "kubernetes.io/not-supported",
						Data: map[string][]byte{
							".dockercfg": secret,
						},
					},
					nil,
				)
				secretListerMock.EXPECT().Secrets("test").Return(secretNamespaceListerMock)
				return secretListerMock
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			provider, err := createTestProvider(ctx, nil, test.secretMapListerMock(), nil, nil)
			if err != nil {
				t.Fatal("failed to create the test provider", err)
			}
			imagePullCredentials, err := provider.getImagePullCredentialsFrom(test.namespace, test.k8sImagePullSecrets)
			if err != nil {
				assert.ErrorContains(t, err, test.err, test.description)
			} else {
				assert.Equal(t, test.expected, imagePullCredentials, test.description)
			}
		})
	}

}

package auth

import (
	"context"
	"fmt"

	httptransport "github.com/go-openapi/runtime/client"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	workloadsScope = "workloads"
	defaultPath    = "/"
	httpProtocol   = "https"
)

type authenticator struct {
	conf *clientcredentials.Config
}

// newAuthenticator returns a new authenticator for StackPath API
func newAuthenticator(ctx context.Context, clientID string, clientSecret string, apiHost string) *authenticator {

	oauthConfig := clientcredentials.Config{
		AuthStyle:    oauth2.AuthStyleInParams,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     fmt.Sprintf("%s://%s/identity/v1/oauth2/token", httpProtocol, apiHost),
		Scopes:       []string{workloadsScope},
	}

	return &authenticator{
		conf: &oauthConfig,
	}
}

// getAccessToken retrieves an access token to authenticate with the StackPath API
func (a *authenticator) GetAccessToken(ctx context.Context) (oauth2.TokenSource, error) {
	tokenSource := a.conf.TokenSource(ctx)

	return tokenSource, nil
}

// NewRuntime creates a new runtime that represents an API client that uses the transport
// to make HTTP requests based on a swagger specification.
func NewRuntime(ctx context.Context, clientID string, clientSecret string, apiHost string, version string) (*httptransport.Runtime, error) {
	auth := newAuthenticator(ctx, clientID, clientSecret, apiHost)
	tokenSource, err := auth.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	// Create a new http Client that will pull authentication tokens
	// from the configured token source
	client := oauth2.NewClient(ctx, tokenSource)

	// Create a new openAPI runtime
	runtime := httptransport.NewWithClient(apiHost, defaultPath, []string{httpProtocol}, client)
	runtime.Transport = NewUserAgentTransport(runtime.Transport, version)

	return runtime, nil
}

package auth

import (
	"fmt"
	"net/http"
)

const userAgentFormat = "Virtual Kubelet Provider/%s"

// UserAgentTransport is an http RoundTripper that sets a descriptive User-Agent
// header for all StackPath API requests.
type UserAgentTransport struct {
	providerVersion string
	http.RoundTripper
	parent http.RoundTripper
}

// NewUserAgentTransport builds a new UserAgentTransport around the underlying
// RoundTripper.
func NewUserAgentTransport(parent http.RoundTripper, providerVersion string) *UserAgentTransport {
	return &UserAgentTransport{parent: parent, providerVersion: providerVersion}
}

// RoundTrip implements the http.RoundTripper interface, setting a User-Agent
// header on the HTTP request.
func (t *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", fmt.Sprintf(userAgentFormat, t.providerVersion))
	return t.parent.RoundTrip(req)
}

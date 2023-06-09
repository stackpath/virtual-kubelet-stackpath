// Code generated by go-swagger; DO NOT EDIT.

package workload

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewWatchNetworks2Params creates a new WatchNetworks2Params object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewWatchNetworks2Params() *WatchNetworks2Params {
	return &WatchNetworks2Params{
		timeout: cr.DefaultTimeout,
	}
}

// NewWatchNetworks2ParamsWithTimeout creates a new WatchNetworks2Params object
// with the ability to set a timeout on a request.
func NewWatchNetworks2ParamsWithTimeout(timeout time.Duration) *WatchNetworks2Params {
	return &WatchNetworks2Params{
		timeout: timeout,
	}
}

// NewWatchNetworks2ParamsWithContext creates a new WatchNetworks2Params object
// with the ability to set a context for a request.
func NewWatchNetworks2ParamsWithContext(ctx context.Context) *WatchNetworks2Params {
	return &WatchNetworks2Params{
		Context: ctx,
	}
}

// NewWatchNetworks2ParamsWithHTTPClient creates a new WatchNetworks2Params object
// with the ability to set a custom HTTPClient for a request.
func NewWatchNetworks2ParamsWithHTTPClient(client *http.Client) *WatchNetworks2Params {
	return &WatchNetworks2Params{
		HTTPClient: client,
	}
}

/*
WatchNetworks2Params contains all the parameters to send to the API endpoint

	for the watch networks2 operation.

	Typically these are written to a http.Request.
*/
type WatchNetworks2Params struct {

	// StackID.
	StackID string

	// Version.
	Version *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the watch networks2 params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WatchNetworks2Params) WithDefaults() *WatchNetworks2Params {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the watch networks2 params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WatchNetworks2Params) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the watch networks2 params
func (o *WatchNetworks2Params) WithTimeout(timeout time.Duration) *WatchNetworks2Params {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the watch networks2 params
func (o *WatchNetworks2Params) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the watch networks2 params
func (o *WatchNetworks2Params) WithContext(ctx context.Context) *WatchNetworks2Params {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the watch networks2 params
func (o *WatchNetworks2Params) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the watch networks2 params
func (o *WatchNetworks2Params) WithHTTPClient(client *http.Client) *WatchNetworks2Params {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the watch networks2 params
func (o *WatchNetworks2Params) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithStackID adds the stackID to the watch networks2 params
func (o *WatchNetworks2Params) WithStackID(stackID string) *WatchNetworks2Params {
	o.SetStackID(stackID)
	return o
}

// SetStackID adds the stackId to the watch networks2 params
func (o *WatchNetworks2Params) SetStackID(stackID string) {
	o.StackID = stackID
}

// WithVersion adds the version to the watch networks2 params
func (o *WatchNetworks2Params) WithVersion(version *string) *WatchNetworks2Params {
	o.SetVersion(version)
	return o
}

// SetVersion adds the version to the watch networks2 params
func (o *WatchNetworks2Params) SetVersion(version *string) {
	o.Version = version
}

// WriteToRequest writes these params to a swagger request
func (o *WatchNetworks2Params) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param stack_id
	if err := r.SetPathParam("stack_id", o.StackID); err != nil {
		return err
	}

	if o.Version != nil {

		// query param version
		var qrVersion string

		if o.Version != nil {
			qrVersion = *o.Version
		}
		qVersion := qrVersion
		if qVersion != "" {

			if err := r.SetQueryParam("version", qVersion); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

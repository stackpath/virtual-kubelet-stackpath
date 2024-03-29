// Code generated by go-swagger; DO NOT EDIT.

package image

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

	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
)

// NewUpdateImageParams creates a new UpdateImageParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateImageParams() *UpdateImageParams {
	return &UpdateImageParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateImageParamsWithTimeout creates a new UpdateImageParams object
// with the ability to set a timeout on a request.
func NewUpdateImageParamsWithTimeout(timeout time.Duration) *UpdateImageParams {
	return &UpdateImageParams{
		timeout: timeout,
	}
}

// NewUpdateImageParamsWithContext creates a new UpdateImageParams object
// with the ability to set a context for a request.
func NewUpdateImageParamsWithContext(ctx context.Context) *UpdateImageParams {
	return &UpdateImageParams{
		Context: ctx,
	}
}

// NewUpdateImageParamsWithHTTPClient creates a new UpdateImageParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateImageParamsWithHTTPClient(client *http.Client) *UpdateImageParams {
	return &UpdateImageParams{
		HTTPClient: client,
	}
}

/*
UpdateImageParams contains all the parameters to send to the API endpoint

	for the update image operation.

	Typically these are written to a http.Request.
*/
type UpdateImageParams struct {

	// Body.
	Body *workload_models.V1UpdateImageRequest

	/* ImageFamily.

	   An image's family
	*/
	ImageFamily string

	/* ImageStackID.

	   The ID of the stack that an image belongs to
	*/
	ImageStackID string

	/* ImageTag.

	   The image's tag
	*/
	ImageTag string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update image params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateImageParams) WithDefaults() *UpdateImageParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update image params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateImageParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update image params
func (o *UpdateImageParams) WithTimeout(timeout time.Duration) *UpdateImageParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update image params
func (o *UpdateImageParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update image params
func (o *UpdateImageParams) WithContext(ctx context.Context) *UpdateImageParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update image params
func (o *UpdateImageParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update image params
func (o *UpdateImageParams) WithHTTPClient(client *http.Client) *UpdateImageParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update image params
func (o *UpdateImageParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update image params
func (o *UpdateImageParams) WithBody(body *workload_models.V1UpdateImageRequest) *UpdateImageParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update image params
func (o *UpdateImageParams) SetBody(body *workload_models.V1UpdateImageRequest) {
	o.Body = body
}

// WithImageFamily adds the imageFamily to the update image params
func (o *UpdateImageParams) WithImageFamily(imageFamily string) *UpdateImageParams {
	o.SetImageFamily(imageFamily)
	return o
}

// SetImageFamily adds the imageFamily to the update image params
func (o *UpdateImageParams) SetImageFamily(imageFamily string) {
	o.ImageFamily = imageFamily
}

// WithImageStackID adds the imageStackID to the update image params
func (o *UpdateImageParams) WithImageStackID(imageStackID string) *UpdateImageParams {
	o.SetImageStackID(imageStackID)
	return o
}

// SetImageStackID adds the imageStackId to the update image params
func (o *UpdateImageParams) SetImageStackID(imageStackID string) {
	o.ImageStackID = imageStackID
}

// WithImageTag adds the imageTag to the update image params
func (o *UpdateImageParams) WithImageTag(imageTag string) *UpdateImageParams {
	o.SetImageTag(imageTag)
	return o
}

// SetImageTag adds the imageTag to the update image params
func (o *UpdateImageParams) SetImageTag(imageTag string) {
	o.ImageTag = imageTag
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateImageParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param image.family
	if err := r.SetPathParam("image.family", o.ImageFamily); err != nil {
		return err
	}

	// path param image.stack_id
	if err := r.SetPathParam("image.stack_id", o.ImageStackID); err != nil {
		return err
	}

	// path param image.tag
	if err := r.SetPathParam("image.tag", o.ImageTag); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

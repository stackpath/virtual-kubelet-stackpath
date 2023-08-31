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
	"github.com/go-openapi/swag"
)

// NewGetImagesParams creates a new GetImagesParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetImagesParams() *GetImagesParams {
	return &GetImagesParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetImagesParamsWithTimeout creates a new GetImagesParams object
// with the ability to set a timeout on a request.
func NewGetImagesParamsWithTimeout(timeout time.Duration) *GetImagesParams {
	return &GetImagesParams{
		timeout: timeout,
	}
}

// NewGetImagesParamsWithContext creates a new GetImagesParams object
// with the ability to set a context for a request.
func NewGetImagesParamsWithContext(ctx context.Context) *GetImagesParams {
	return &GetImagesParams{
		Context: ctx,
	}
}

// NewGetImagesParamsWithHTTPClient creates a new GetImagesParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetImagesParamsWithHTTPClient(client *http.Client) *GetImagesParams {
	return &GetImagesParams{
		HTTPClient: client,
	}
}

/*
GetImagesParams contains all the parameters to send to the API endpoint

	for the get images operation.

	Typically these are written to a http.Request.
*/
type GetImagesParams struct {

	/* Deprecated.

	   If present and true, include deprecated images in the result.

	   Format: boolean
	*/
	Deprecated *bool

	/* PageRequestAfter.

	   The cursor value after which data will be returned.
	*/
	PageRequestAfter *string

	/* PageRequestFilter.

	   SQL-style constraint filters.
	*/
	PageRequestFilter *string

	/* PageRequestFirst.

	   The number of items desired.
	*/
	PageRequestFirst *string

	/* PageRequestSortBy.

	   Sort the response by the given field.
	*/
	PageRequestSortBy *string

	/* StackID.

	   The ID or slug of the stack to retrieve images for
	*/
	StackID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get images params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetImagesParams) WithDefaults() *GetImagesParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get images params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetImagesParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get images params
func (o *GetImagesParams) WithTimeout(timeout time.Duration) *GetImagesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get images params
func (o *GetImagesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get images params
func (o *GetImagesParams) WithContext(ctx context.Context) *GetImagesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get images params
func (o *GetImagesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get images params
func (o *GetImagesParams) WithHTTPClient(client *http.Client) *GetImagesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get images params
func (o *GetImagesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDeprecated adds the deprecated to the get images params
func (o *GetImagesParams) WithDeprecated(deprecated *bool) *GetImagesParams {
	o.SetDeprecated(deprecated)
	return o
}

// SetDeprecated adds the deprecated to the get images params
func (o *GetImagesParams) SetDeprecated(deprecated *bool) {
	o.Deprecated = deprecated
}

// WithPageRequestAfter adds the pageRequestAfter to the get images params
func (o *GetImagesParams) WithPageRequestAfter(pageRequestAfter *string) *GetImagesParams {
	o.SetPageRequestAfter(pageRequestAfter)
	return o
}

// SetPageRequestAfter adds the pageRequestAfter to the get images params
func (o *GetImagesParams) SetPageRequestAfter(pageRequestAfter *string) {
	o.PageRequestAfter = pageRequestAfter
}

// WithPageRequestFilter adds the pageRequestFilter to the get images params
func (o *GetImagesParams) WithPageRequestFilter(pageRequestFilter *string) *GetImagesParams {
	o.SetPageRequestFilter(pageRequestFilter)
	return o
}

// SetPageRequestFilter adds the pageRequestFilter to the get images params
func (o *GetImagesParams) SetPageRequestFilter(pageRequestFilter *string) {
	o.PageRequestFilter = pageRequestFilter
}

// WithPageRequestFirst adds the pageRequestFirst to the get images params
func (o *GetImagesParams) WithPageRequestFirst(pageRequestFirst *string) *GetImagesParams {
	o.SetPageRequestFirst(pageRequestFirst)
	return o
}

// SetPageRequestFirst adds the pageRequestFirst to the get images params
func (o *GetImagesParams) SetPageRequestFirst(pageRequestFirst *string) {
	o.PageRequestFirst = pageRequestFirst
}

// WithPageRequestSortBy adds the pageRequestSortBy to the get images params
func (o *GetImagesParams) WithPageRequestSortBy(pageRequestSortBy *string) *GetImagesParams {
	o.SetPageRequestSortBy(pageRequestSortBy)
	return o
}

// SetPageRequestSortBy adds the pageRequestSortBy to the get images params
func (o *GetImagesParams) SetPageRequestSortBy(pageRequestSortBy *string) {
	o.PageRequestSortBy = pageRequestSortBy
}

// WithStackID adds the stackID to the get images params
func (o *GetImagesParams) WithStackID(stackID string) *GetImagesParams {
	o.SetStackID(stackID)
	return o
}

// SetStackID adds the stackId to the get images params
func (o *GetImagesParams) SetStackID(stackID string) {
	o.StackID = stackID
}

// WriteToRequest writes these params to a swagger request
func (o *GetImagesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Deprecated != nil {

		// query param deprecated
		var qrDeprecated bool

		if o.Deprecated != nil {
			qrDeprecated = *o.Deprecated
		}
		qDeprecated := swag.FormatBool(qrDeprecated)
		if qDeprecated != "" {

			if err := r.SetQueryParam("deprecated", qDeprecated); err != nil {
				return err
			}
		}
	}

	if o.PageRequestAfter != nil {

		// query param page_request.after
		var qrPageRequestAfter string

		if o.PageRequestAfter != nil {
			qrPageRequestAfter = *o.PageRequestAfter
		}
		qPageRequestAfter := qrPageRequestAfter
		if qPageRequestAfter != "" {

			if err := r.SetQueryParam("page_request.after", qPageRequestAfter); err != nil {
				return err
			}
		}
	}

	if o.PageRequestFilter != nil {

		// query param page_request.filter
		var qrPageRequestFilter string

		if o.PageRequestFilter != nil {
			qrPageRequestFilter = *o.PageRequestFilter
		}
		qPageRequestFilter := qrPageRequestFilter
		if qPageRequestFilter != "" {

			if err := r.SetQueryParam("page_request.filter", qPageRequestFilter); err != nil {
				return err
			}
		}
	}

	if o.PageRequestFirst != nil {

		// query param page_request.first
		var qrPageRequestFirst string

		if o.PageRequestFirst != nil {
			qrPageRequestFirst = *o.PageRequestFirst
		}
		qPageRequestFirst := qrPageRequestFirst
		if qPageRequestFirst != "" {

			if err := r.SetQueryParam("page_request.first", qPageRequestFirst); err != nil {
				return err
			}
		}
	}

	if o.PageRequestSortBy != nil {

		// query param page_request.sort_by
		var qrPageRequestSortBy string

		if o.PageRequestSortBy != nil {
			qrPageRequestSortBy = *o.PageRequestSortBy
		}
		qPageRequestSortBy := qrPageRequestSortBy
		if qPageRequestSortBy != "" {

			if err := r.SetQueryParam("page_request.sort_by", qPageRequestSortBy); err != nil {
				return err
			}
		}
	}

	// path param stack_id
	if err := r.SetPathParam("stack_id", o.StackID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

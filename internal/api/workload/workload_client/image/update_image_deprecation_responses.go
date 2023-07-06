// Code generated by go-swagger; DO NOT EDIT.

package image

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
)

// UpdateImageDeprecationReader is a Reader for the UpdateImageDeprecation structure.
type UpdateImageDeprecationReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateImageDeprecationReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateImageDeprecationOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewUpdateImageDeprecationUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateImageDeprecationInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewUpdateImageDeprecationDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateImageDeprecationOK creates a UpdateImageDeprecationOK with default headers values
func NewUpdateImageDeprecationOK() *UpdateImageDeprecationOK {
	return &UpdateImageDeprecationOK{}
}

/*
UpdateImageDeprecationOK describes a response with status code 200, with default header values.

UpdateImageDeprecationOK update image deprecation o k
*/
type UpdateImageDeprecationOK struct {
	Payload *workload_models.V1UpdateImageDeprecationResponse
}

// IsSuccess returns true when this update image deprecation o k response has a 2xx status code
func (o *UpdateImageDeprecationOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this update image deprecation o k response has a 3xx status code
func (o *UpdateImageDeprecationOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update image deprecation o k response has a 4xx status code
func (o *UpdateImageDeprecationOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this update image deprecation o k response has a 5xx status code
func (o *UpdateImageDeprecationOK) IsServerError() bool {
	return false
}

// IsCode returns true when this update image deprecation o k response a status code equal to that given
func (o *UpdateImageDeprecationOK) IsCode(code int) bool {
	return code == 200
}

func (o *UpdateImageDeprecationOK) Error() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{stack_id}/images/{image_family}/{image_tag}/deprecation][%d] updateImageDeprecationOK  %+v", 200, o.Payload)
}

func (o *UpdateImageDeprecationOK) String() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{stack_id}/images/{image_family}/{image_tag}/deprecation][%d] updateImageDeprecationOK  %+v", 200, o.Payload)
}

func (o *UpdateImageDeprecationOK) GetPayload() *workload_models.V1UpdateImageDeprecationResponse {
	return o.Payload
}

func (o *UpdateImageDeprecationOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.V1UpdateImageDeprecationResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateImageDeprecationUnauthorized creates a UpdateImageDeprecationUnauthorized with default headers values
func NewUpdateImageDeprecationUnauthorized() *UpdateImageDeprecationUnauthorized {
	return &UpdateImageDeprecationUnauthorized{}
}

/*
UpdateImageDeprecationUnauthorized describes a response with status code 401, with default header values.

Returned when an unauthorized request is attempted.
*/
type UpdateImageDeprecationUnauthorized struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this update image deprecation unauthorized response has a 2xx status code
func (o *UpdateImageDeprecationUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update image deprecation unauthorized response has a 3xx status code
func (o *UpdateImageDeprecationUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update image deprecation unauthorized response has a 4xx status code
func (o *UpdateImageDeprecationUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this update image deprecation unauthorized response has a 5xx status code
func (o *UpdateImageDeprecationUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this update image deprecation unauthorized response a status code equal to that given
func (o *UpdateImageDeprecationUnauthorized) IsCode(code int) bool {
	return code == 401
}

func (o *UpdateImageDeprecationUnauthorized) Error() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{stack_id}/images/{image_family}/{image_tag}/deprecation][%d] updateImageDeprecationUnauthorized  %+v", 401, o.Payload)
}

func (o *UpdateImageDeprecationUnauthorized) String() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{stack_id}/images/{image_family}/{image_tag}/deprecation][%d] updateImageDeprecationUnauthorized  %+v", 401, o.Payload)
}

func (o *UpdateImageDeprecationUnauthorized) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *UpdateImageDeprecationUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateImageDeprecationInternalServerError creates a UpdateImageDeprecationInternalServerError with default headers values
func NewUpdateImageDeprecationInternalServerError() *UpdateImageDeprecationInternalServerError {
	return &UpdateImageDeprecationInternalServerError{}
}

/*
UpdateImageDeprecationInternalServerError describes a response with status code 500, with default header values.

Internal server error.
*/
type UpdateImageDeprecationInternalServerError struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this update image deprecation internal server error response has a 2xx status code
func (o *UpdateImageDeprecationInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update image deprecation internal server error response has a 3xx status code
func (o *UpdateImageDeprecationInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update image deprecation internal server error response has a 4xx status code
func (o *UpdateImageDeprecationInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this update image deprecation internal server error response has a 5xx status code
func (o *UpdateImageDeprecationInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this update image deprecation internal server error response a status code equal to that given
func (o *UpdateImageDeprecationInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *UpdateImageDeprecationInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{stack_id}/images/{image_family}/{image_tag}/deprecation][%d] updateImageDeprecationInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateImageDeprecationInternalServerError) String() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{stack_id}/images/{image_family}/{image_tag}/deprecation][%d] updateImageDeprecationInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateImageDeprecationInternalServerError) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *UpdateImageDeprecationInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateImageDeprecationDefault creates a UpdateImageDeprecationDefault with default headers values
func NewUpdateImageDeprecationDefault(code int) *UpdateImageDeprecationDefault {
	return &UpdateImageDeprecationDefault{
		_statusCode: code,
	}
}

/*
UpdateImageDeprecationDefault describes a response with status code -1, with default header values.

Default error structure.
*/
type UpdateImageDeprecationDefault struct {
	_statusCode int

	Payload *workload_models.StackpathapiStatus
}

// Code gets the status code for the update image deprecation default response
func (o *UpdateImageDeprecationDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this update image deprecation default response has a 2xx status code
func (o *UpdateImageDeprecationDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this update image deprecation default response has a 3xx status code
func (o *UpdateImageDeprecationDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this update image deprecation default response has a 4xx status code
func (o *UpdateImageDeprecationDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this update image deprecation default response has a 5xx status code
func (o *UpdateImageDeprecationDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this update image deprecation default response a status code equal to that given
func (o *UpdateImageDeprecationDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *UpdateImageDeprecationDefault) Error() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{stack_id}/images/{image_family}/{image_tag}/deprecation][%d] UpdateImageDeprecation default  %+v", o._statusCode, o.Payload)
}

func (o *UpdateImageDeprecationDefault) String() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{stack_id}/images/{image_family}/{image_tag}/deprecation][%d] UpdateImageDeprecation default  %+v", o._statusCode, o.Payload)
}

func (o *UpdateImageDeprecationDefault) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *UpdateImageDeprecationDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
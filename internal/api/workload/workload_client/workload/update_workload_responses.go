// Code generated by go-swagger; DO NOT EDIT.

package workload

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
)

// UpdateWorkloadReader is a Reader for the UpdateWorkload structure.
type UpdateWorkloadReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateWorkloadReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateWorkloadOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewUpdateWorkloadUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateWorkloadInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewUpdateWorkloadDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateWorkloadOK creates a UpdateWorkloadOK with default headers values
func NewUpdateWorkloadOK() *UpdateWorkloadOK {
	return &UpdateWorkloadOK{}
}

/*
UpdateWorkloadOK describes a response with status code 200, with default header values.

UpdateWorkloadOK update workload o k
*/
type UpdateWorkloadOK struct {
	Payload *workload_models.V1UpdateWorkloadResponse
}

// IsSuccess returns true when this update workload o k response has a 2xx status code
func (o *UpdateWorkloadOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this update workload o k response has a 3xx status code
func (o *UpdateWorkloadOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update workload o k response has a 4xx status code
func (o *UpdateWorkloadOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this update workload o k response has a 5xx status code
func (o *UpdateWorkloadOK) IsServerError() bool {
	return false
}

// IsCode returns true when this update workload o k response a status code equal to that given
func (o *UpdateWorkloadOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the update workload o k response
func (o *UpdateWorkloadOK) Code() int {
	return 200
}

func (o *UpdateWorkloadOK) Error() string {
	return fmt.Sprintf("[PATCH /workload/v1/stacks/{stack_id}/workloads/{workload_id}][%d] updateWorkloadOK  %+v", 200, o.Payload)
}

func (o *UpdateWorkloadOK) String() string {
	return fmt.Sprintf("[PATCH /workload/v1/stacks/{stack_id}/workloads/{workload_id}][%d] updateWorkloadOK  %+v", 200, o.Payload)
}

func (o *UpdateWorkloadOK) GetPayload() *workload_models.V1UpdateWorkloadResponse {
	return o.Payload
}

func (o *UpdateWorkloadOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.V1UpdateWorkloadResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateWorkloadUnauthorized creates a UpdateWorkloadUnauthorized with default headers values
func NewUpdateWorkloadUnauthorized() *UpdateWorkloadUnauthorized {
	return &UpdateWorkloadUnauthorized{}
}

/*
UpdateWorkloadUnauthorized describes a response with status code 401, with default header values.

Returned when an unauthorized request is attempted.
*/
type UpdateWorkloadUnauthorized struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this update workload unauthorized response has a 2xx status code
func (o *UpdateWorkloadUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update workload unauthorized response has a 3xx status code
func (o *UpdateWorkloadUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update workload unauthorized response has a 4xx status code
func (o *UpdateWorkloadUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this update workload unauthorized response has a 5xx status code
func (o *UpdateWorkloadUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this update workload unauthorized response a status code equal to that given
func (o *UpdateWorkloadUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the update workload unauthorized response
func (o *UpdateWorkloadUnauthorized) Code() int {
	return 401
}

func (o *UpdateWorkloadUnauthorized) Error() string {
	return fmt.Sprintf("[PATCH /workload/v1/stacks/{stack_id}/workloads/{workload_id}][%d] updateWorkloadUnauthorized  %+v", 401, o.Payload)
}

func (o *UpdateWorkloadUnauthorized) String() string {
	return fmt.Sprintf("[PATCH /workload/v1/stacks/{stack_id}/workloads/{workload_id}][%d] updateWorkloadUnauthorized  %+v", 401, o.Payload)
}

func (o *UpdateWorkloadUnauthorized) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *UpdateWorkloadUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateWorkloadInternalServerError creates a UpdateWorkloadInternalServerError with default headers values
func NewUpdateWorkloadInternalServerError() *UpdateWorkloadInternalServerError {
	return &UpdateWorkloadInternalServerError{}
}

/*
UpdateWorkloadInternalServerError describes a response with status code 500, with default header values.

Internal server error.
*/
type UpdateWorkloadInternalServerError struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this update workload internal server error response has a 2xx status code
func (o *UpdateWorkloadInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update workload internal server error response has a 3xx status code
func (o *UpdateWorkloadInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update workload internal server error response has a 4xx status code
func (o *UpdateWorkloadInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this update workload internal server error response has a 5xx status code
func (o *UpdateWorkloadInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this update workload internal server error response a status code equal to that given
func (o *UpdateWorkloadInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the update workload internal server error response
func (o *UpdateWorkloadInternalServerError) Code() int {
	return 500
}

func (o *UpdateWorkloadInternalServerError) Error() string {
	return fmt.Sprintf("[PATCH /workload/v1/stacks/{stack_id}/workloads/{workload_id}][%d] updateWorkloadInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateWorkloadInternalServerError) String() string {
	return fmt.Sprintf("[PATCH /workload/v1/stacks/{stack_id}/workloads/{workload_id}][%d] updateWorkloadInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateWorkloadInternalServerError) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *UpdateWorkloadInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateWorkloadDefault creates a UpdateWorkloadDefault with default headers values
func NewUpdateWorkloadDefault(code int) *UpdateWorkloadDefault {
	return &UpdateWorkloadDefault{
		_statusCode: code,
	}
}

/*
UpdateWorkloadDefault describes a response with status code -1, with default header values.

Default error structure.
*/
type UpdateWorkloadDefault struct {
	_statusCode int

	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this update workload default response has a 2xx status code
func (o *UpdateWorkloadDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this update workload default response has a 3xx status code
func (o *UpdateWorkloadDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this update workload default response has a 4xx status code
func (o *UpdateWorkloadDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this update workload default response has a 5xx status code
func (o *UpdateWorkloadDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this update workload default response a status code equal to that given
func (o *UpdateWorkloadDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the update workload default response
func (o *UpdateWorkloadDefault) Code() int {
	return o._statusCode
}

func (o *UpdateWorkloadDefault) Error() string {
	return fmt.Sprintf("[PATCH /workload/v1/stacks/{stack_id}/workloads/{workload_id}][%d] UpdateWorkload default  %+v", o._statusCode, o.Payload)
}

func (o *UpdateWorkloadDefault) String() string {
	return fmt.Sprintf("[PATCH /workload/v1/stacks/{stack_id}/workloads/{workload_id}][%d] UpdateWorkload default  %+v", o._statusCode, o.Payload)
}

func (o *UpdateWorkloadDefault) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *UpdateWorkloadDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

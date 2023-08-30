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

// PutWorkloadReader is a Reader for the PutWorkload structure.
type PutWorkloadReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PutWorkloadReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPutWorkloadOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewPutWorkloadUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewPutWorkloadInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewPutWorkloadDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPutWorkloadOK creates a PutWorkloadOK with default headers values
func NewPutWorkloadOK() *PutWorkloadOK {
	return &PutWorkloadOK{}
}

/*
PutWorkloadOK describes a response with status code 200, with default header values.

PutWorkloadOK put workload o k
*/
type PutWorkloadOK struct {
	Payload *workload_models.V1PutWorkloadResponse
}

// IsSuccess returns true when this put workload o k response has a 2xx status code
func (o *PutWorkloadOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this put workload o k response has a 3xx status code
func (o *PutWorkloadOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this put workload o k response has a 4xx status code
func (o *PutWorkloadOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this put workload o k response has a 5xx status code
func (o *PutWorkloadOK) IsServerError() bool {
	return false
}

// IsCode returns true when this put workload o k response a status code equal to that given
func (o *PutWorkloadOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the put workload o k response
func (o *PutWorkloadOK) Code() int {
	return 200
}

func (o *PutWorkloadOK) Error() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{workload.stack_id}/workloads/{workload.id}][%d] putWorkloadOK  %+v", 200, o.Payload)
}

func (o *PutWorkloadOK) String() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{workload.stack_id}/workloads/{workload.id}][%d] putWorkloadOK  %+v", 200, o.Payload)
}

func (o *PutWorkloadOK) GetPayload() *workload_models.V1PutWorkloadResponse {
	return o.Payload
}

func (o *PutWorkloadOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.V1PutWorkloadResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutWorkloadUnauthorized creates a PutWorkloadUnauthorized with default headers values
func NewPutWorkloadUnauthorized() *PutWorkloadUnauthorized {
	return &PutWorkloadUnauthorized{}
}

/*
PutWorkloadUnauthorized describes a response with status code 401, with default header values.

Returned when an unauthorized request is attempted.
*/
type PutWorkloadUnauthorized struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this put workload unauthorized response has a 2xx status code
func (o *PutWorkloadUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this put workload unauthorized response has a 3xx status code
func (o *PutWorkloadUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this put workload unauthorized response has a 4xx status code
func (o *PutWorkloadUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this put workload unauthorized response has a 5xx status code
func (o *PutWorkloadUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this put workload unauthorized response a status code equal to that given
func (o *PutWorkloadUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the put workload unauthorized response
func (o *PutWorkloadUnauthorized) Code() int {
	return 401
}

func (o *PutWorkloadUnauthorized) Error() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{workload.stack_id}/workloads/{workload.id}][%d] putWorkloadUnauthorized  %+v", 401, o.Payload)
}

func (o *PutWorkloadUnauthorized) String() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{workload.stack_id}/workloads/{workload.id}][%d] putWorkloadUnauthorized  %+v", 401, o.Payload)
}

func (o *PutWorkloadUnauthorized) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *PutWorkloadUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutWorkloadInternalServerError creates a PutWorkloadInternalServerError with default headers values
func NewPutWorkloadInternalServerError() *PutWorkloadInternalServerError {
	return &PutWorkloadInternalServerError{}
}

/*
PutWorkloadInternalServerError describes a response with status code 500, with default header values.

Internal server error.
*/
type PutWorkloadInternalServerError struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this put workload internal server error response has a 2xx status code
func (o *PutWorkloadInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this put workload internal server error response has a 3xx status code
func (o *PutWorkloadInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this put workload internal server error response has a 4xx status code
func (o *PutWorkloadInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this put workload internal server error response has a 5xx status code
func (o *PutWorkloadInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this put workload internal server error response a status code equal to that given
func (o *PutWorkloadInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the put workload internal server error response
func (o *PutWorkloadInternalServerError) Code() int {
	return 500
}

func (o *PutWorkloadInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{workload.stack_id}/workloads/{workload.id}][%d] putWorkloadInternalServerError  %+v", 500, o.Payload)
}

func (o *PutWorkloadInternalServerError) String() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{workload.stack_id}/workloads/{workload.id}][%d] putWorkloadInternalServerError  %+v", 500, o.Payload)
}

func (o *PutWorkloadInternalServerError) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *PutWorkloadInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutWorkloadDefault creates a PutWorkloadDefault with default headers values
func NewPutWorkloadDefault(code int) *PutWorkloadDefault {
	return &PutWorkloadDefault{
		_statusCode: code,
	}
}

/*
PutWorkloadDefault describes a response with status code -1, with default header values.

Default error structure.
*/
type PutWorkloadDefault struct {
	_statusCode int

	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this put workload default response has a 2xx status code
func (o *PutWorkloadDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this put workload default response has a 3xx status code
func (o *PutWorkloadDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this put workload default response has a 4xx status code
func (o *PutWorkloadDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this put workload default response has a 5xx status code
func (o *PutWorkloadDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this put workload default response a status code equal to that given
func (o *PutWorkloadDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the put workload default response
func (o *PutWorkloadDefault) Code() int {
	return o._statusCode
}

func (o *PutWorkloadDefault) Error() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{workload.stack_id}/workloads/{workload.id}][%d] PutWorkload default  %+v", o._statusCode, o.Payload)
}

func (o *PutWorkloadDefault) String() string {
	return fmt.Sprintf("[PUT /workload/v1/stacks/{workload.stack_id}/workloads/{workload.id}][%d] PutWorkload default  %+v", o._statusCode, o.Payload)
}

func (o *PutWorkloadDefault) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *PutWorkloadDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

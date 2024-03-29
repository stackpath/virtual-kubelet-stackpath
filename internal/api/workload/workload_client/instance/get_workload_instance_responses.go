// Code generated by go-swagger; DO NOT EDIT.

package instance

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
)

// GetWorkloadInstanceReader is a Reader for the GetWorkloadInstance structure.
type GetWorkloadInstanceReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetWorkloadInstanceReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetWorkloadInstanceOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetWorkloadInstanceUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetWorkloadInstanceInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetWorkloadInstanceDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetWorkloadInstanceOK creates a GetWorkloadInstanceOK with default headers values
func NewGetWorkloadInstanceOK() *GetWorkloadInstanceOK {
	return &GetWorkloadInstanceOK{}
}

/*
GetWorkloadInstanceOK describes a response with status code 200, with default header values.

GetWorkloadInstanceOK get workload instance o k
*/
type GetWorkloadInstanceOK struct {
	Payload *workload_models.V1GetWorkloadInstanceResponse
}

// IsSuccess returns true when this get workload instance o k response has a 2xx status code
func (o *GetWorkloadInstanceOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get workload instance o k response has a 3xx status code
func (o *GetWorkloadInstanceOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get workload instance o k response has a 4xx status code
func (o *GetWorkloadInstanceOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get workload instance o k response has a 5xx status code
func (o *GetWorkloadInstanceOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get workload instance o k response a status code equal to that given
func (o *GetWorkloadInstanceOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get workload instance o k response
func (o *GetWorkloadInstanceOK) Code() int {
	return 200
}

func (o *GetWorkloadInstanceOK) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}][%d] getWorkloadInstanceOK  %+v", 200, o.Payload)
}

func (o *GetWorkloadInstanceOK) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}][%d] getWorkloadInstanceOK  %+v", 200, o.Payload)
}

func (o *GetWorkloadInstanceOK) GetPayload() *workload_models.V1GetWorkloadInstanceResponse {
	return o.Payload
}

func (o *GetWorkloadInstanceOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.V1GetWorkloadInstanceResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkloadInstanceUnauthorized creates a GetWorkloadInstanceUnauthorized with default headers values
func NewGetWorkloadInstanceUnauthorized() *GetWorkloadInstanceUnauthorized {
	return &GetWorkloadInstanceUnauthorized{}
}

/*
GetWorkloadInstanceUnauthorized describes a response with status code 401, with default header values.

Returned when an unauthorized request is attempted.
*/
type GetWorkloadInstanceUnauthorized struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get workload instance unauthorized response has a 2xx status code
func (o *GetWorkloadInstanceUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get workload instance unauthorized response has a 3xx status code
func (o *GetWorkloadInstanceUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get workload instance unauthorized response has a 4xx status code
func (o *GetWorkloadInstanceUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this get workload instance unauthorized response has a 5xx status code
func (o *GetWorkloadInstanceUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this get workload instance unauthorized response a status code equal to that given
func (o *GetWorkloadInstanceUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the get workload instance unauthorized response
func (o *GetWorkloadInstanceUnauthorized) Code() int {
	return 401
}

func (o *GetWorkloadInstanceUnauthorized) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}][%d] getWorkloadInstanceUnauthorized  %+v", 401, o.Payload)
}

func (o *GetWorkloadInstanceUnauthorized) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}][%d] getWorkloadInstanceUnauthorized  %+v", 401, o.Payload)
}

func (o *GetWorkloadInstanceUnauthorized) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetWorkloadInstanceUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkloadInstanceInternalServerError creates a GetWorkloadInstanceInternalServerError with default headers values
func NewGetWorkloadInstanceInternalServerError() *GetWorkloadInstanceInternalServerError {
	return &GetWorkloadInstanceInternalServerError{}
}

/*
GetWorkloadInstanceInternalServerError describes a response with status code 500, with default header values.

Internal server error.
*/
type GetWorkloadInstanceInternalServerError struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get workload instance internal server error response has a 2xx status code
func (o *GetWorkloadInstanceInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get workload instance internal server error response has a 3xx status code
func (o *GetWorkloadInstanceInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get workload instance internal server error response has a 4xx status code
func (o *GetWorkloadInstanceInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get workload instance internal server error response has a 5xx status code
func (o *GetWorkloadInstanceInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get workload instance internal server error response a status code equal to that given
func (o *GetWorkloadInstanceInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get workload instance internal server error response
func (o *GetWorkloadInstanceInternalServerError) Code() int {
	return 500
}

func (o *GetWorkloadInstanceInternalServerError) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}][%d] getWorkloadInstanceInternalServerError  %+v", 500, o.Payload)
}

func (o *GetWorkloadInstanceInternalServerError) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}][%d] getWorkloadInstanceInternalServerError  %+v", 500, o.Payload)
}

func (o *GetWorkloadInstanceInternalServerError) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetWorkloadInstanceInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkloadInstanceDefault creates a GetWorkloadInstanceDefault with default headers values
func NewGetWorkloadInstanceDefault(code int) *GetWorkloadInstanceDefault {
	return &GetWorkloadInstanceDefault{
		_statusCode: code,
	}
}

/*
GetWorkloadInstanceDefault describes a response with status code -1, with default header values.

Default error structure.
*/
type GetWorkloadInstanceDefault struct {
	_statusCode int

	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get workload instance default response has a 2xx status code
func (o *GetWorkloadInstanceDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get workload instance default response has a 3xx status code
func (o *GetWorkloadInstanceDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get workload instance default response has a 4xx status code
func (o *GetWorkloadInstanceDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get workload instance default response has a 5xx status code
func (o *GetWorkloadInstanceDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get workload instance default response a status code equal to that given
func (o *GetWorkloadInstanceDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get workload instance default response
func (o *GetWorkloadInstanceDefault) Code() int {
	return o._statusCode
}

func (o *GetWorkloadInstanceDefault) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}][%d] GetWorkloadInstance default  %+v", o._statusCode, o.Payload)
}

func (o *GetWorkloadInstanceDefault) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}][%d] GetWorkloadInstance default  %+v", o._statusCode, o.Payload)
}

func (o *GetWorkloadInstanceDefault) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetWorkloadInstanceDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

package instances

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
)

// GetWorkloadInstancesReader is a Reader for the GetWorkloadInstances structure.
type GetWorkloadInstancesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetWorkloadInstancesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetWorkloadInstancesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetWorkloadInstancesUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetWorkloadInstancesInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetWorkloadInstancesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetWorkloadInstancesOK creates a GetWorkloadInstancesOK with default headers values
func NewGetWorkloadInstancesOK() *GetWorkloadInstancesOK {
	return &GetWorkloadInstancesOK{}
}

/*
GetWorkloadInstancesOK describes a response with status code 200, with default header values.

GetWorkloadInstancesOK get workload instances o k
*/
type GetWorkloadInstancesOK struct {
	Payload *workload_models.V1GetWorkloadInstancesResponse
}

// IsSuccess returns true when this get workload instances o k response has a 2xx status code
func (o *GetWorkloadInstancesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get workload instances o k response has a 3xx status code
func (o *GetWorkloadInstancesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get workload instances o k response has a 4xx status code
func (o *GetWorkloadInstancesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get workload instances o k response has a 5xx status code
func (o *GetWorkloadInstancesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get workload instances o k response a status code equal to that given
func (o *GetWorkloadInstancesOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get workload instances o k response
func (o *GetWorkloadInstancesOK) Code() int {
	return 200
}

func (o *GetWorkloadInstancesOK) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances][%d] getWorkloadInstancesOK  %+v", 200, o.Payload)
}

func (o *GetWorkloadInstancesOK) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances][%d] getWorkloadInstancesOK  %+v", 200, o.Payload)
}

func (o *GetWorkloadInstancesOK) GetPayload() *workload_models.V1GetWorkloadInstancesResponse {
	return o.Payload
}

func (o *GetWorkloadInstancesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.V1GetWorkloadInstancesResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkloadInstancesUnauthorized creates a GetWorkloadInstancesUnauthorized with default headers values
func NewGetWorkloadInstancesUnauthorized() *GetWorkloadInstancesUnauthorized {
	return &GetWorkloadInstancesUnauthorized{}
}

/*
GetWorkloadInstancesUnauthorized describes a response with status code 401, with default header values.

Returned when an unauthorized request is attempted.
*/
type GetWorkloadInstancesUnauthorized struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get workload instances unauthorized response has a 2xx status code
func (o *GetWorkloadInstancesUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get workload instances unauthorized response has a 3xx status code
func (o *GetWorkloadInstancesUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get workload instances unauthorized response has a 4xx status code
func (o *GetWorkloadInstancesUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this get workload instances unauthorized response has a 5xx status code
func (o *GetWorkloadInstancesUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this get workload instances unauthorized response a status code equal to that given
func (o *GetWorkloadInstancesUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the get workload instances unauthorized response
func (o *GetWorkloadInstancesUnauthorized) Code() int {
	return 401
}

func (o *GetWorkloadInstancesUnauthorized) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances][%d] getWorkloadInstancesUnauthorized  %+v", 401, o.Payload)
}

func (o *GetWorkloadInstancesUnauthorized) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances][%d] getWorkloadInstancesUnauthorized  %+v", 401, o.Payload)
}

func (o *GetWorkloadInstancesUnauthorized) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetWorkloadInstancesUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkloadInstancesInternalServerError creates a GetWorkloadInstancesInternalServerError with default headers values
func NewGetWorkloadInstancesInternalServerError() *GetWorkloadInstancesInternalServerError {
	return &GetWorkloadInstancesInternalServerError{}
}

/*
GetWorkloadInstancesInternalServerError describes a response with status code 500, with default header values.

Internal server error.
*/
type GetWorkloadInstancesInternalServerError struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get workload instances internal server error response has a 2xx status code
func (o *GetWorkloadInstancesInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get workload instances internal server error response has a 3xx status code
func (o *GetWorkloadInstancesInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get workload instances internal server error response has a 4xx status code
func (o *GetWorkloadInstancesInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get workload instances internal server error response has a 5xx status code
func (o *GetWorkloadInstancesInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get workload instances internal server error response a status code equal to that given
func (o *GetWorkloadInstancesInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get workload instances internal server error response
func (o *GetWorkloadInstancesInternalServerError) Code() int {
	return 500
}

func (o *GetWorkloadInstancesInternalServerError) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances][%d] getWorkloadInstancesInternalServerError  %+v", 500, o.Payload)
}

func (o *GetWorkloadInstancesInternalServerError) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances][%d] getWorkloadInstancesInternalServerError  %+v", 500, o.Payload)
}

func (o *GetWorkloadInstancesInternalServerError) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetWorkloadInstancesInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkloadInstancesDefault creates a GetWorkloadInstancesDefault with default headers values
func NewGetWorkloadInstancesDefault(code int) *GetWorkloadInstancesDefault {
	return &GetWorkloadInstancesDefault{
		_statusCode: code,
	}
}

/*
GetWorkloadInstancesDefault describes a response with status code -1, with default header values.

Default error structure.
*/
type GetWorkloadInstancesDefault struct {
	_statusCode int

	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get workload instances default response has a 2xx status code
func (o *GetWorkloadInstancesDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get workload instances default response has a 3xx status code
func (o *GetWorkloadInstancesDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get workload instances default response has a 4xx status code
func (o *GetWorkloadInstancesDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get workload instances default response has a 5xx status code
func (o *GetWorkloadInstancesDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get workload instances default response a status code equal to that given
func (o *GetWorkloadInstancesDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get workload instances default response
func (o *GetWorkloadInstancesDefault) Code() int {
	return o._statusCode
}

func (o *GetWorkloadInstancesDefault) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances][%d] GetWorkloadInstances default  %+v", o._statusCode, o.Payload)
}

func (o *GetWorkloadInstancesDefault) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances][%d] GetWorkloadInstances default  %+v", o._statusCode, o.Payload)
}

func (o *GetWorkloadInstancesDefault) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetWorkloadInstancesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

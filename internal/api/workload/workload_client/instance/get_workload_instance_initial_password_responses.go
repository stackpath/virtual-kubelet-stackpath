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

// GetWorkloadInstanceInitialPasswordReader is a Reader for the GetWorkloadInstanceInitialPassword structure.
type GetWorkloadInstanceInitialPasswordReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetWorkloadInstanceInitialPasswordReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetWorkloadInstanceInitialPasswordOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetWorkloadInstanceInitialPasswordUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetWorkloadInstanceInitialPasswordInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetWorkloadInstanceInitialPasswordDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetWorkloadInstanceInitialPasswordOK creates a GetWorkloadInstanceInitialPasswordOK with default headers values
func NewGetWorkloadInstanceInitialPasswordOK() *GetWorkloadInstanceInitialPasswordOK {
	return &GetWorkloadInstanceInitialPasswordOK{}
}

/*
GetWorkloadInstanceInitialPasswordOK describes a response with status code 200, with default header values.

GetWorkloadInstanceInitialPasswordOK get workload instance initial password o k
*/
type GetWorkloadInstanceInitialPasswordOK struct {
	Payload *workload_models.V1GetWorkloadInstanceInitialPasswordResponse
}

// IsSuccess returns true when this get workload instance initial password o k response has a 2xx status code
func (o *GetWorkloadInstanceInitialPasswordOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get workload instance initial password o k response has a 3xx status code
func (o *GetWorkloadInstanceInitialPasswordOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get workload instance initial password o k response has a 4xx status code
func (o *GetWorkloadInstanceInitialPasswordOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get workload instance initial password o k response has a 5xx status code
func (o *GetWorkloadInstanceInitialPasswordOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get workload instance initial password o k response a status code equal to that given
func (o *GetWorkloadInstanceInitialPasswordOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get workload instance initial password o k response
func (o *GetWorkloadInstanceInitialPasswordOK) Code() int {
	return 200
}

func (o *GetWorkloadInstanceInitialPasswordOK) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/passwords/initial][%d] getWorkloadInstanceInitialPasswordOK  %+v", 200, o.Payload)
}

func (o *GetWorkloadInstanceInitialPasswordOK) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/passwords/initial][%d] getWorkloadInstanceInitialPasswordOK  %+v", 200, o.Payload)
}

func (o *GetWorkloadInstanceInitialPasswordOK) GetPayload() *workload_models.V1GetWorkloadInstanceInitialPasswordResponse {
	return o.Payload
}

func (o *GetWorkloadInstanceInitialPasswordOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.V1GetWorkloadInstanceInitialPasswordResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkloadInstanceInitialPasswordUnauthorized creates a GetWorkloadInstanceInitialPasswordUnauthorized with default headers values
func NewGetWorkloadInstanceInitialPasswordUnauthorized() *GetWorkloadInstanceInitialPasswordUnauthorized {
	return &GetWorkloadInstanceInitialPasswordUnauthorized{}
}

/*
GetWorkloadInstanceInitialPasswordUnauthorized describes a response with status code 401, with default header values.

Returned when an unauthorized request is attempted.
*/
type GetWorkloadInstanceInitialPasswordUnauthorized struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get workload instance initial password unauthorized response has a 2xx status code
func (o *GetWorkloadInstanceInitialPasswordUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get workload instance initial password unauthorized response has a 3xx status code
func (o *GetWorkloadInstanceInitialPasswordUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get workload instance initial password unauthorized response has a 4xx status code
func (o *GetWorkloadInstanceInitialPasswordUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this get workload instance initial password unauthorized response has a 5xx status code
func (o *GetWorkloadInstanceInitialPasswordUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this get workload instance initial password unauthorized response a status code equal to that given
func (o *GetWorkloadInstanceInitialPasswordUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the get workload instance initial password unauthorized response
func (o *GetWorkloadInstanceInitialPasswordUnauthorized) Code() int {
	return 401
}

func (o *GetWorkloadInstanceInitialPasswordUnauthorized) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/passwords/initial][%d] getWorkloadInstanceInitialPasswordUnauthorized  %+v", 401, o.Payload)
}

func (o *GetWorkloadInstanceInitialPasswordUnauthorized) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/passwords/initial][%d] getWorkloadInstanceInitialPasswordUnauthorized  %+v", 401, o.Payload)
}

func (o *GetWorkloadInstanceInitialPasswordUnauthorized) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetWorkloadInstanceInitialPasswordUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkloadInstanceInitialPasswordInternalServerError creates a GetWorkloadInstanceInitialPasswordInternalServerError with default headers values
func NewGetWorkloadInstanceInitialPasswordInternalServerError() *GetWorkloadInstanceInitialPasswordInternalServerError {
	return &GetWorkloadInstanceInitialPasswordInternalServerError{}
}

/*
GetWorkloadInstanceInitialPasswordInternalServerError describes a response with status code 500, with default header values.

Internal server error.
*/
type GetWorkloadInstanceInitialPasswordInternalServerError struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get workload instance initial password internal server error response has a 2xx status code
func (o *GetWorkloadInstanceInitialPasswordInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get workload instance initial password internal server error response has a 3xx status code
func (o *GetWorkloadInstanceInitialPasswordInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get workload instance initial password internal server error response has a 4xx status code
func (o *GetWorkloadInstanceInitialPasswordInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get workload instance initial password internal server error response has a 5xx status code
func (o *GetWorkloadInstanceInitialPasswordInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get workload instance initial password internal server error response a status code equal to that given
func (o *GetWorkloadInstanceInitialPasswordInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get workload instance initial password internal server error response
func (o *GetWorkloadInstanceInitialPasswordInternalServerError) Code() int {
	return 500
}

func (o *GetWorkloadInstanceInitialPasswordInternalServerError) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/passwords/initial][%d] getWorkloadInstanceInitialPasswordInternalServerError  %+v", 500, o.Payload)
}

func (o *GetWorkloadInstanceInitialPasswordInternalServerError) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/passwords/initial][%d] getWorkloadInstanceInitialPasswordInternalServerError  %+v", 500, o.Payload)
}

func (o *GetWorkloadInstanceInitialPasswordInternalServerError) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetWorkloadInstanceInitialPasswordInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkloadInstanceInitialPasswordDefault creates a GetWorkloadInstanceInitialPasswordDefault with default headers values
func NewGetWorkloadInstanceInitialPasswordDefault(code int) *GetWorkloadInstanceInitialPasswordDefault {
	return &GetWorkloadInstanceInitialPasswordDefault{
		_statusCode: code,
	}
}

/*
GetWorkloadInstanceInitialPasswordDefault describes a response with status code -1, with default header values.

Default error structure.
*/
type GetWorkloadInstanceInitialPasswordDefault struct {
	_statusCode int

	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get workload instance initial password default response has a 2xx status code
func (o *GetWorkloadInstanceInitialPasswordDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get workload instance initial password default response has a 3xx status code
func (o *GetWorkloadInstanceInitialPasswordDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get workload instance initial password default response has a 4xx status code
func (o *GetWorkloadInstanceInitialPasswordDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get workload instance initial password default response has a 5xx status code
func (o *GetWorkloadInstanceInitialPasswordDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get workload instance initial password default response a status code equal to that given
func (o *GetWorkloadInstanceInitialPasswordDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get workload instance initial password default response
func (o *GetWorkloadInstanceInitialPasswordDefault) Code() int {
	return o._statusCode
}

func (o *GetWorkloadInstanceInitialPasswordDefault) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/passwords/initial][%d] GetWorkloadInstanceInitialPassword default  %+v", o._statusCode, o.Payload)
}

func (o *GetWorkloadInstanceInitialPasswordDefault) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/passwords/initial][%d] GetWorkloadInstanceInitialPassword default  %+v", o._statusCode, o.Payload)
}

func (o *GetWorkloadInstanceInitialPasswordDefault) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetWorkloadInstanceInitialPasswordDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

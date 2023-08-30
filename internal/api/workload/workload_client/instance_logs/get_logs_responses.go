// Code generated by go-swagger; DO NOT EDIT.

package instance_logs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
)

// GetLogsReader is a Reader for the GetLogs structure.
type GetLogsReader struct {
	formats strfmt.Registry
	writer  io.Writer
}

// ReadResponse reads a server response into the received o.
func (o *GetLogsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetLogsOK(o.writer)
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetLogsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetLogsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetLogsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetLogsOK creates a GetLogsOK with default headers values
func NewGetLogsOK(writer io.Writer) *GetLogsOK {
	return &GetLogsOK{

		Payload: writer,
	}
}

/*
GetLogsOK describes a response with status code 200, with default header values.

(streaming responses)
*/
type GetLogsOK struct {
	Payload io.Writer
}

// IsSuccess returns true when this get logs o k response has a 2xx status code
func (o *GetLogsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get logs o k response has a 3xx status code
func (o *GetLogsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get logs o k response has a 4xx status code
func (o *GetLogsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get logs o k response has a 5xx status code
func (o *GetLogsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get logs o k response a status code equal to that given
func (o *GetLogsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get logs o k response
func (o *GetLogsOK) Code() int {
	return 200
}

func (o *GetLogsOK) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/logs][%d] getLogsOK  %+v", 200, o.Payload)
}

func (o *GetLogsOK) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/logs][%d] getLogsOK  %+v", 200, o.Payload)
}

func (o *GetLogsOK) GetPayload() io.Writer {
	return o.Payload
}

func (o *GetLogsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetLogsUnauthorized creates a GetLogsUnauthorized with default headers values
func NewGetLogsUnauthorized() *GetLogsUnauthorized {
	return &GetLogsUnauthorized{}
}

/*
GetLogsUnauthorized describes a response with status code 401, with default header values.

Returned when an unauthorized request is attempted.
*/
type GetLogsUnauthorized struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get logs unauthorized response has a 2xx status code
func (o *GetLogsUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get logs unauthorized response has a 3xx status code
func (o *GetLogsUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get logs unauthorized response has a 4xx status code
func (o *GetLogsUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this get logs unauthorized response has a 5xx status code
func (o *GetLogsUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this get logs unauthorized response a status code equal to that given
func (o *GetLogsUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the get logs unauthorized response
func (o *GetLogsUnauthorized) Code() int {
	return 401
}

func (o *GetLogsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/logs][%d] getLogsUnauthorized  %+v", 401, o.Payload)
}

func (o *GetLogsUnauthorized) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/logs][%d] getLogsUnauthorized  %+v", 401, o.Payload)
}

func (o *GetLogsUnauthorized) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetLogsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetLogsInternalServerError creates a GetLogsInternalServerError with default headers values
func NewGetLogsInternalServerError() *GetLogsInternalServerError {
	return &GetLogsInternalServerError{}
}

/*
GetLogsInternalServerError describes a response with status code 500, with default header values.

Internal server error.
*/
type GetLogsInternalServerError struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get logs internal server error response has a 2xx status code
func (o *GetLogsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get logs internal server error response has a 3xx status code
func (o *GetLogsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get logs internal server error response has a 4xx status code
func (o *GetLogsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get logs internal server error response has a 5xx status code
func (o *GetLogsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get logs internal server error response a status code equal to that given
func (o *GetLogsInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get logs internal server error response
func (o *GetLogsInternalServerError) Code() int {
	return 500
}

func (o *GetLogsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/logs][%d] getLogsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetLogsInternalServerError) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/logs][%d] getLogsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetLogsInternalServerError) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetLogsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetLogsDefault creates a GetLogsDefault with default headers values
func NewGetLogsDefault(code int) *GetLogsDefault {
	return &GetLogsDefault{
		_statusCode: code,
	}
}

/*
GetLogsDefault describes a response with status code -1, with default header values.

Default error structure.
*/
type GetLogsDefault struct {
	_statusCode int

	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get logs default response has a 2xx status code
func (o *GetLogsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get logs default response has a 3xx status code
func (o *GetLogsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get logs default response has a 4xx status code
func (o *GetLogsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get logs default response has a 5xx status code
func (o *GetLogsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get logs default response a status code equal to that given
func (o *GetLogsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get logs default response
func (o *GetLogsDefault) Code() int {
	return o._statusCode
}

func (o *GetLogsDefault) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/logs][%d] GetLogs default  %+v", o._statusCode, o.Payload)
}

func (o *GetLogsDefault) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/workloads/{workload_id}/instances/{instance_name}/logs][%d] GetLogs default  %+v", o._statusCode, o.Payload)
}

func (o *GetLogsDefault) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetLogsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

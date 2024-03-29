// Code generated by go-swagger; DO NOT EDIT.

package metrics

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/stackpath/virtual-kubelet-stackpath/internal/api/workload/workload_models"
)

// GetMetricsReader is a Reader for the GetMetrics structure.
type GetMetricsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetMetricsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetMetricsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetMetricsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetMetricsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetMetricsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetMetricsOK creates a GetMetricsOK with default headers values
func NewGetMetricsOK() *GetMetricsOK {
	return &GetMetricsOK{}
}

/*
GetMetricsOK describes a response with status code 200, with default header values.

GetMetricsOK get metrics o k
*/
type GetMetricsOK struct {
	Payload *workload_models.PrometheusMetrics
}

// IsSuccess returns true when this get metrics o k response has a 2xx status code
func (o *GetMetricsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get metrics o k response has a 3xx status code
func (o *GetMetricsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get metrics o k response has a 4xx status code
func (o *GetMetricsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get metrics o k response has a 5xx status code
func (o *GetMetricsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get metrics o k response a status code equal to that given
func (o *GetMetricsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get metrics o k response
func (o *GetMetricsOK) Code() int {
	return 200
}

func (o *GetMetricsOK) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/metrics][%d] getMetricsOK  %+v", 200, o.Payload)
}

func (o *GetMetricsOK) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/metrics][%d] getMetricsOK  %+v", 200, o.Payload)
}

func (o *GetMetricsOK) GetPayload() *workload_models.PrometheusMetrics {
	return o.Payload
}

func (o *GetMetricsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.PrometheusMetrics)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetMetricsUnauthorized creates a GetMetricsUnauthorized with default headers values
func NewGetMetricsUnauthorized() *GetMetricsUnauthorized {
	return &GetMetricsUnauthorized{}
}

/*
GetMetricsUnauthorized describes a response with status code 401, with default header values.

Returned when an unauthorized request is attempted.
*/
type GetMetricsUnauthorized struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get metrics unauthorized response has a 2xx status code
func (o *GetMetricsUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get metrics unauthorized response has a 3xx status code
func (o *GetMetricsUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get metrics unauthorized response has a 4xx status code
func (o *GetMetricsUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this get metrics unauthorized response has a 5xx status code
func (o *GetMetricsUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this get metrics unauthorized response a status code equal to that given
func (o *GetMetricsUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the get metrics unauthorized response
func (o *GetMetricsUnauthorized) Code() int {
	return 401
}

func (o *GetMetricsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/metrics][%d] getMetricsUnauthorized  %+v", 401, o.Payload)
}

func (o *GetMetricsUnauthorized) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/metrics][%d] getMetricsUnauthorized  %+v", 401, o.Payload)
}

func (o *GetMetricsUnauthorized) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetMetricsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetMetricsInternalServerError creates a GetMetricsInternalServerError with default headers values
func NewGetMetricsInternalServerError() *GetMetricsInternalServerError {
	return &GetMetricsInternalServerError{}
}

/*
GetMetricsInternalServerError describes a response with status code 500, with default header values.

Internal server error.
*/
type GetMetricsInternalServerError struct {
	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get metrics internal server error response has a 2xx status code
func (o *GetMetricsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get metrics internal server error response has a 3xx status code
func (o *GetMetricsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get metrics internal server error response has a 4xx status code
func (o *GetMetricsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get metrics internal server error response has a 5xx status code
func (o *GetMetricsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get metrics internal server error response a status code equal to that given
func (o *GetMetricsInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get metrics internal server error response
func (o *GetMetricsInternalServerError) Code() int {
	return 500
}

func (o *GetMetricsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/metrics][%d] getMetricsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetMetricsInternalServerError) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/metrics][%d] getMetricsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetMetricsInternalServerError) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetMetricsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetMetricsDefault creates a GetMetricsDefault with default headers values
func NewGetMetricsDefault(code int) *GetMetricsDefault {
	return &GetMetricsDefault{
		_statusCode: code,
	}
}

/*
GetMetricsDefault describes a response with status code -1, with default header values.

Default error structure.
*/
type GetMetricsDefault struct {
	_statusCode int

	Payload *workload_models.StackpathapiStatus
}

// IsSuccess returns true when this get metrics default response has a 2xx status code
func (o *GetMetricsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get metrics default response has a 3xx status code
func (o *GetMetricsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get metrics default response has a 4xx status code
func (o *GetMetricsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get metrics default response has a 5xx status code
func (o *GetMetricsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get metrics default response a status code equal to that given
func (o *GetMetricsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get metrics default response
func (o *GetMetricsDefault) Code() int {
	return o._statusCode
}

func (o *GetMetricsDefault) Error() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/metrics][%d] GetMetrics default  %+v", o._statusCode, o.Payload)
}

func (o *GetMetricsDefault) String() string {
	return fmt.Sprintf("[GET /workload/v1/stacks/{stack_id}/metrics][%d] GetMetrics default  %+v", o._statusCode, o.Payload)
}

func (o *GetMetricsDefault) GetPayload() *workload_models.StackpathapiStatus {
	return o.Payload
}

func (o *GetMetricsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(workload_models.StackpathapiStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

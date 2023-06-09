// Code generated by go-swagger; DO NOT EDIT.

package workload_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// PrometheusMetricsStatus A metrics query's resulting status
//
// swagger:model prometheusMetricsStatus
type PrometheusMetricsStatus string

func NewPrometheusMetricsStatus(value PrometheusMetricsStatus) *PrometheusMetricsStatus {
	return &value
}

// Pointer returns a pointer to a freshly-allocated PrometheusMetricsStatus.
func (m PrometheusMetricsStatus) Pointer() *PrometheusMetricsStatus {
	return &m
}

const (

	// PrometheusMetricsStatusSUCCESS captures enum value "SUCCESS"
	PrometheusMetricsStatusSUCCESS PrometheusMetricsStatus = "SUCCESS"

	// PrometheusMetricsStatusERROR captures enum value "ERROR"
	PrometheusMetricsStatusERROR PrometheusMetricsStatus = "ERROR"
)

// for schema
var prometheusMetricsStatusEnum []interface{}

func init() {
	var res []PrometheusMetricsStatus
	if err := json.Unmarshal([]byte(`["SUCCESS","ERROR"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		prometheusMetricsStatusEnum = append(prometheusMetricsStatusEnum, v)
	}
}

func (m PrometheusMetricsStatus) validatePrometheusMetricsStatusEnum(path, location string, value PrometheusMetricsStatus) error {
	if err := validate.EnumCase(path, location, value, prometheusMetricsStatusEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this prometheus metrics status
func (m PrometheusMetricsStatus) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validatePrometheusMetricsStatusEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this prometheus metrics status based on context it is used
func (m PrometheusMetricsStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

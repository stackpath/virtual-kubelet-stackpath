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

// V1ConditionStatus v1 condition status
//
// swagger:model v1ConditionStatus
type V1ConditionStatus string

func NewV1ConditionStatus(value V1ConditionStatus) *V1ConditionStatus {
	return &value
}

// Pointer returns a pointer to a freshly-allocated V1ConditionStatus.
func (m V1ConditionStatus) Pointer() *V1ConditionStatus {
	return &m
}

const (

	// V1ConditionStatusCONDITIONSTATUSUNKNOWN captures enum value "CONDITION_STATUS_UNKNOWN"
	V1ConditionStatusCONDITIONSTATUSUNKNOWN V1ConditionStatus = "CONDITION_STATUS_UNKNOWN"

	// V1ConditionStatusCONDITIONSTATUSTRUE captures enum value "CONDITION_STATUS_TRUE"
	V1ConditionStatusCONDITIONSTATUSTRUE V1ConditionStatus = "CONDITION_STATUS_TRUE"

	// V1ConditionStatusCONDITIONSTATUSFALSE captures enum value "CONDITION_STATUS_FALSE"
	V1ConditionStatusCONDITIONSTATUSFALSE V1ConditionStatus = "CONDITION_STATUS_FALSE"
)

// for schema
var v1ConditionStatusEnum []interface{}

func init() {
	var res []V1ConditionStatus
	if err := json.Unmarshal([]byte(`["CONDITION_STATUS_UNKNOWN","CONDITION_STATUS_TRUE","CONDITION_STATUS_FALSE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		v1ConditionStatusEnum = append(v1ConditionStatusEnum, v)
	}
}

func (m V1ConditionStatus) validateV1ConditionStatusEnum(path, location string, value V1ConditionStatus) error {
	if err := validate.EnumCase(path, location, value, v1ConditionStatusEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this v1 condition status
func (m V1ConditionStatus) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateV1ConditionStatusEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this v1 condition status based on context it is used
func (m V1ConditionStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

package workload_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1DNSConfig Specification to control DNS settings for workload instanceIt is used to control resolver config such as DNS nameservers, search domains and resolver options for workload instance
//
// swagger:model v1DNSConfig
type V1DNSConfig struct {

	// A list of DNS name server IP addresses
	Nameservers []string `json:"nameservers"`

	// A list of DNS resolver options
	Options []*V1DNSConfigOption `json:"options"`

	// A list of DNS search domains for hostname lookup
	Searches []string `json:"searches"`
}

// Validate validates this v1 DNS config
func (m *V1DNSConfig) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateOptions(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1DNSConfig) validateOptions(formats strfmt.Registry) error {
	if swag.IsZero(m.Options) { // not required
		return nil
	}

	for i := 0; i < len(m.Options); i++ {
		if swag.IsZero(m.Options[i]) { // not required
			continue
		}

		if m.Options[i] != nil {
			if err := m.Options[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("options" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("options" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this v1 DNS config based on the context it is used
func (m *V1DNSConfig) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateOptions(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1DNSConfig) contextValidateOptions(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Options); i++ {

		if m.Options[i] != nil {
			if err := m.Options[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("options" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("options" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1DNSConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1DNSConfig) UnmarshalBinary(b []byte) error {
	var res V1DNSConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

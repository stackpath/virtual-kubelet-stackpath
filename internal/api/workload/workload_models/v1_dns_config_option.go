// Code generated by go-swagger; DO NOT EDIT.

package workload_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1DNSConfigOption Specification for DNS resolver options, More information about DNS resolveroptions are available here - https://www.man7.org/linux/man-pages/man5/resolv.conf.5.html
//
// swagger:model v1DNSConfigOption
type V1DNSConfigOption struct {

	// DNS config name
	Name string `json:"name,omitempty"`

	// DNS config value
	Value string `json:"value,omitempty"`
}

// Validate validates this v1 DNS config option
func (m *V1DNSConfigOption) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this v1 DNS config option based on context it is used
func (m *V1DNSConfigOption) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *V1DNSConfigOption) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1DNSConfigOption) UnmarshalBinary(b []byte) error {
	var res V1DNSConfigOption
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
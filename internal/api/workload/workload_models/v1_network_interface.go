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

// V1NetworkInterface Network interfaces that will be created on instances in the workload
//
// swagger:model v1NetworkInterface
type V1NetworkInterface struct {

	// Provide one to one NAT for this specific network’s interface
	//
	// This is an optional property used to enable OR disable the NAT'ing of the specifc network interface.NAT is enabled by default on the first / primary interface and disabled on secondary / multi interfacesUser can disable NAT on the first / primary interface, by marking this property - falseUser can enable NAT on secondary / multi interface, by marking this property - trueExcludeNAT Annotation supercedes this property
	EnableOneToOneNat bool `json:"enableOneToOneNat,omitempty"`

	// A list of IP families to use for interface ip assignments
	//
	// This is an optional property and supports ['IPv4'] or ['IPv4', 'IPv6'] list
	IPFamilies []*V1IPFamily `json:"ipFamilies"`

	// An IPv6 subnet interface's slug. This is an optional property used to attach a specific network interface to a ipv6 subnet.
	IPV6Subnet string `json:"ipv6Subnet,omitempty"`

	// A network interface's slug
	Network string `json:"network,omitempty"`

	// An IPv4 subnet interface's slug. This is an optional property used to attach a specific network interface to a IPv4 subnet.
	Subnet string `json:"subnet,omitempty"`
}

// Validate validates this v1 network interface
func (m *V1NetworkInterface) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateIPFamilies(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1NetworkInterface) validateIPFamilies(formats strfmt.Registry) error {
	if swag.IsZero(m.IPFamilies) { // not required
		return nil
	}

	for i := 0; i < len(m.IPFamilies); i++ {
		if swag.IsZero(m.IPFamilies[i]) { // not required
			continue
		}

		if m.IPFamilies[i] != nil {
			if err := m.IPFamilies[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("ipFamilies" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("ipFamilies" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this v1 network interface based on the context it is used
func (m *V1NetworkInterface) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateIPFamilies(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1NetworkInterface) contextValidateIPFamilies(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.IPFamilies); i++ {

		if m.IPFamilies[i] != nil {
			if err := m.IPFamilies[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("ipFamilies" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("ipFamilies" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1NetworkInterface) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1NetworkInterface) UnmarshalBinary(b []byte) error {
	var res V1NetworkInterface
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

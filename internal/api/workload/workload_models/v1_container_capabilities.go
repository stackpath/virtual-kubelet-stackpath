// Code generated by go-swagger; DO NOT EDIT.

package workload_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1ContainerCapabilities The specification for capabilities to add/drop when running containersSupported capabilities are [cap_chown, cap_dac_override, cap_dac_read_search, cap_fowner, cap_fsetid, cap_kill, cap_setgid, cap_setuid, cap_setpcap, cap_linux_immutable,cap_net_bind_service, cap_net_broadcast, cap_net_admin, cap_net_raw, cap_ipc_lock, cap_ipc_owner,cap_sys_chroot, cap_sys_ptrace, cap_sys_pacct, cap_sys_nice, cap_mknod, cap_lease, cap_setfcap]
//
// swagger:model v1ContainerCapabilities
type V1ContainerCapabilities struct {

	// List of supported capabilities to add in container
	Add []string `json:"add"`

	// List of supported capabilities to drop from container
	Drop []string `json:"drop"`
}

// Validate validates this v1 container capabilities
func (m *V1ContainerCapabilities) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this v1 container capabilities based on context it is used
func (m *V1ContainerCapabilities) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *V1ContainerCapabilities) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1ContainerCapabilities) UnmarshalBinary(b []byte) error {
	var res V1ContainerCapabilities
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

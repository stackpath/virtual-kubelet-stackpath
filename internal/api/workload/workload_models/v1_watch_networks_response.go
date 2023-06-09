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

// V1WatchNetworksResponse v1 watch networks response
//
// swagger:model v1WatchNetworksResponse
type V1WatchNetworksResponse struct {

	// cluster
	Cluster string `json:"cluster,omitempty"`

	// event type
	EventType *WatchNetworksResponseEventType `json:"eventType,omitempty"`

	// instance conditions
	InstanceConditions []*V1InstanceCondition `json:"instanceConditions"`

	// instance name
	InstanceName string `json:"instanceName,omitempty"`

	// network
	Network string `json:"network,omitempty"`

	// network interface
	NetworkInterface *V1WatchNetworksResponseNetworkInterfaceStatus `json:"networkInterface,omitempty"`

	// stack Id
	StackID string `json:"stackId,omitempty"`

	// version
	Version string `json:"version,omitempty"`

	// watcher state
	WatcherState *V1WatcherState `json:"watcherState,omitempty"`
}

// Validate validates this v1 watch networks response
func (m *V1WatchNetworksResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEventType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInstanceConditions(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworkInterface(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWatcherState(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1WatchNetworksResponse) validateEventType(formats strfmt.Registry) error {
	if swag.IsZero(m.EventType) { // not required
		return nil
	}

	if m.EventType != nil {
		if err := m.EventType.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("eventType")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("eventType")
			}
			return err
		}
	}

	return nil
}

func (m *V1WatchNetworksResponse) validateInstanceConditions(formats strfmt.Registry) error {
	if swag.IsZero(m.InstanceConditions) { // not required
		return nil
	}

	for i := 0; i < len(m.InstanceConditions); i++ {
		if swag.IsZero(m.InstanceConditions[i]) { // not required
			continue
		}

		if m.InstanceConditions[i] != nil {
			if err := m.InstanceConditions[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("instanceConditions" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("instanceConditions" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *V1WatchNetworksResponse) validateNetworkInterface(formats strfmt.Registry) error {
	if swag.IsZero(m.NetworkInterface) { // not required
		return nil
	}

	if m.NetworkInterface != nil {
		if err := m.NetworkInterface.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("networkInterface")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("networkInterface")
			}
			return err
		}
	}

	return nil
}

func (m *V1WatchNetworksResponse) validateWatcherState(formats strfmt.Registry) error {
	if swag.IsZero(m.WatcherState) { // not required
		return nil
	}

	if m.WatcherState != nil {
		if err := m.WatcherState.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("watcherState")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("watcherState")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this v1 watch networks response based on the context it is used
func (m *V1WatchNetworksResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateEventType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateInstanceConditions(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNetworkInterface(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateWatcherState(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1WatchNetworksResponse) contextValidateEventType(ctx context.Context, formats strfmt.Registry) error {

	if m.EventType != nil {
		if err := m.EventType.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("eventType")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("eventType")
			}
			return err
		}
	}

	return nil
}

func (m *V1WatchNetworksResponse) contextValidateInstanceConditions(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.InstanceConditions); i++ {

		if m.InstanceConditions[i] != nil {
			if err := m.InstanceConditions[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("instanceConditions" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("instanceConditions" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *V1WatchNetworksResponse) contextValidateNetworkInterface(ctx context.Context, formats strfmt.Registry) error {

	if m.NetworkInterface != nil {
		if err := m.NetworkInterface.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("networkInterface")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("networkInterface")
			}
			return err
		}
	}

	return nil
}

func (m *V1WatchNetworksResponse) contextValidateWatcherState(ctx context.Context, formats strfmt.Registry) error {

	if m.WatcherState != nil {
		if err := m.WatcherState.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("watcherState")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("watcherState")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1WatchNetworksResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1WatchNetworksResponse) UnmarshalBinary(b []byte) error {
	var res V1WatchNetworksResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

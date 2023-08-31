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

// V1ContainerSpec The specification for the desired state of a container in a workload
//
// swagger:model v1ContainerSpec
type V1ContainerSpec struct {

	// Arguments to the container entry point
	Args []string `json:"args"`

	// The commands that start a container
	Command []string `json:"command"`

	// env
	Env V1EnvironmentVariableMapEntry `json:"env,omitempty"`

	// The location of a Docker image to run as a container
	Image string `json:"image,omitempty"`

	// image pull policy
	ImagePullPolicy *V1ContainerImagePullPolicy `json:"imagePullPolicy,omitempty"`

	// lifecycle
	Lifecycle *V1ContainerLifecycle `json:"lifecycle,omitempty"`

	// liveness probe
	LivenessProbe *V1Probe `json:"livenessProbe,omitempty"`

	// ports
	Ports V1InstancePortMapEntry `json:"ports,omitempty"`

	// readiness probe
	ReadinessProbe *V1Probe `json:"readinessProbe,omitempty"`

	// resources
	Resources *V1ResourceRequirements `json:"resources,omitempty"`

	// security context
	SecurityContext *V1ContainerSecurityContext `json:"securityContext,omitempty"`

	// startup probe
	StartupProbe *V1Probe `json:"startupProbe,omitempty"`

	// Mounted file path at which the container's termination message will be written
	TerminationMessagePath string `json:"terminationMessagePath,omitempty"`

	// termination message policy
	TerminationMessagePolicy *V1ContainerTerminationMessagePolicy `json:"terminationMessagePolicy,omitempty"`

	// Volumes to mount in the container
	VolumeMounts []*V1InstanceVolumeMount `json:"volumeMounts"`

	// Container's working directory
	WorkingDir string `json:"workingDir,omitempty"`
}

// Validate validates this v1 container spec
func (m *V1ContainerSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEnv(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateImagePullPolicy(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLifecycle(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLivenessProbe(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePorts(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateReadinessProbe(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResources(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSecurityContext(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStartupProbe(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTerminationMessagePolicy(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVolumeMounts(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1ContainerSpec) validateEnv(formats strfmt.Registry) error {
	if swag.IsZero(m.Env) { // not required
		return nil
	}

	if m.Env != nil {
		if err := m.Env.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("env")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("env")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) validateImagePullPolicy(formats strfmt.Registry) error {
	if swag.IsZero(m.ImagePullPolicy) { // not required
		return nil
	}

	if m.ImagePullPolicy != nil {
		if err := m.ImagePullPolicy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("imagePullPolicy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("imagePullPolicy")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) validateLifecycle(formats strfmt.Registry) error {
	if swag.IsZero(m.Lifecycle) { // not required
		return nil
	}

	if m.Lifecycle != nil {
		if err := m.Lifecycle.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("lifecycle")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("lifecycle")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) validateLivenessProbe(formats strfmt.Registry) error {
	if swag.IsZero(m.LivenessProbe) { // not required
		return nil
	}

	if m.LivenessProbe != nil {
		if err := m.LivenessProbe.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("livenessProbe")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("livenessProbe")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) validatePorts(formats strfmt.Registry) error {
	if swag.IsZero(m.Ports) { // not required
		return nil
	}

	if m.Ports != nil {
		if err := m.Ports.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("ports")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("ports")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) validateReadinessProbe(formats strfmt.Registry) error {
	if swag.IsZero(m.ReadinessProbe) { // not required
		return nil
	}

	if m.ReadinessProbe != nil {
		if err := m.ReadinessProbe.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("readinessProbe")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("readinessProbe")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) validateResources(formats strfmt.Registry) error {
	if swag.IsZero(m.Resources) { // not required
		return nil
	}

	if m.Resources != nil {
		if err := m.Resources.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("resources")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("resources")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) validateSecurityContext(formats strfmt.Registry) error {
	if swag.IsZero(m.SecurityContext) { // not required
		return nil
	}

	if m.SecurityContext != nil {
		if err := m.SecurityContext.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("securityContext")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("securityContext")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) validateStartupProbe(formats strfmt.Registry) error {
	if swag.IsZero(m.StartupProbe) { // not required
		return nil
	}

	if m.StartupProbe != nil {
		if err := m.StartupProbe.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("startupProbe")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("startupProbe")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) validateTerminationMessagePolicy(formats strfmt.Registry) error {
	if swag.IsZero(m.TerminationMessagePolicy) { // not required
		return nil
	}

	if m.TerminationMessagePolicy != nil {
		if err := m.TerminationMessagePolicy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("terminationMessagePolicy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("terminationMessagePolicy")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) validateVolumeMounts(formats strfmt.Registry) error {
	if swag.IsZero(m.VolumeMounts) { // not required
		return nil
	}

	for i := 0; i < len(m.VolumeMounts); i++ {
		if swag.IsZero(m.VolumeMounts[i]) { // not required
			continue
		}

		if m.VolumeMounts[i] != nil {
			if err := m.VolumeMounts[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("volumeMounts" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("volumeMounts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this v1 container spec based on the context it is used
func (m *V1ContainerSpec) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateEnv(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateImagePullPolicy(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateLifecycle(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateLivenessProbe(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePorts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateReadinessProbe(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateResources(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSecurityContext(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStartupProbe(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTerminationMessagePolicy(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateVolumeMounts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1ContainerSpec) contextValidateEnv(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Env.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("env")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("env")
		}
		return err
	}

	return nil
}

func (m *V1ContainerSpec) contextValidateImagePullPolicy(ctx context.Context, formats strfmt.Registry) error {

	if m.ImagePullPolicy != nil {
		if err := m.ImagePullPolicy.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("imagePullPolicy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("imagePullPolicy")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) contextValidateLifecycle(ctx context.Context, formats strfmt.Registry) error {

	if m.Lifecycle != nil {
		if err := m.Lifecycle.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("lifecycle")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("lifecycle")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) contextValidateLivenessProbe(ctx context.Context, formats strfmt.Registry) error {

	if m.LivenessProbe != nil {
		if err := m.LivenessProbe.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("livenessProbe")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("livenessProbe")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) contextValidatePorts(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Ports.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("ports")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("ports")
		}
		return err
	}

	return nil
}

func (m *V1ContainerSpec) contextValidateReadinessProbe(ctx context.Context, formats strfmt.Registry) error {

	if m.ReadinessProbe != nil {
		if err := m.ReadinessProbe.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("readinessProbe")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("readinessProbe")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) contextValidateResources(ctx context.Context, formats strfmt.Registry) error {

	if m.Resources != nil {
		if err := m.Resources.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("resources")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("resources")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) contextValidateSecurityContext(ctx context.Context, formats strfmt.Registry) error {

	if m.SecurityContext != nil {
		if err := m.SecurityContext.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("securityContext")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("securityContext")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) contextValidateStartupProbe(ctx context.Context, formats strfmt.Registry) error {

	if m.StartupProbe != nil {
		if err := m.StartupProbe.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("startupProbe")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("startupProbe")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) contextValidateTerminationMessagePolicy(ctx context.Context, formats strfmt.Registry) error {

	if m.TerminationMessagePolicy != nil {
		if err := m.TerminationMessagePolicy.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("terminationMessagePolicy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("terminationMessagePolicy")
			}
			return err
		}
	}

	return nil
}

func (m *V1ContainerSpec) contextValidateVolumeMounts(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.VolumeMounts); i++ {

		if m.VolumeMounts[i] != nil {
			if err := m.VolumeMounts[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("volumeMounts" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("volumeMounts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1ContainerSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1ContainerSpec) UnmarshalBinary(b []byte) error {
	var res V1ContainerSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// UpdateNoteReq update note req
// swagger:model UpdateNoteReq
type UpdateNoteReq struct {

	// importance
	// Enum: [LOW MEDIUM HIGH]
	Importance string `json:"importance,omitempty"`

	// message
	Message *string `json:"message,omitempty"`
}

// Validate validates this update note req
func (m *UpdateNoteReq) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateImportance(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var updateNoteReqTypeImportancePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["LOW","MEDIUM","HIGH"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		updateNoteReqTypeImportancePropEnum = append(updateNoteReqTypeImportancePropEnum, v)
	}
}

const (

	// UpdateNoteReqImportanceLOW captures enum value "LOW"
	UpdateNoteReqImportanceLOW string = "LOW"

	// UpdateNoteReqImportanceMEDIUM captures enum value "MEDIUM"
	UpdateNoteReqImportanceMEDIUM string = "MEDIUM"

	// UpdateNoteReqImportanceHIGH captures enum value "HIGH"
	UpdateNoteReqImportanceHIGH string = "HIGH"
)

// prop value enum
func (m *UpdateNoteReq) validateImportanceEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, updateNoteReqTypeImportancePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *UpdateNoteReq) validateImportance(formats strfmt.Registry) error {

	if swag.IsZero(m.Importance) { // not required
		return nil
	}

	// value enum
	if err := m.validateImportanceEnum("importance", "body", m.Importance); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *UpdateNoteReq) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UpdateNoteReq) UnmarshalBinary(b []byte) error {
	var res UpdateNoteReq
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NetworkCreate network create
// swagger:model NetworkCreate
type NetworkCreate struct {

	// Network in CIDR notation (192.168.0.0/24)
	Cidr string `json:"cidr,omitempty"`

	// DNS Servers
	DNSServers []string `json:"dnsServers"`

	// Gateway IP Address
	Gateway string `json:"gateway,omitempty"`

	// IP Address Ranges
	IPAddressRanges []*IPAddressRange `json:"ipAddressRanges"`

	// Enable MTU Jumbo Network
	Jumbo bool `json:"jumbo,omitempty"`

	// Network Name
	Name string `json:"name,omitempty"`

	// Type of Network - 'vlan' (private network) 'pub-vlan' (public network)
	// Required: true
	// Enum: [vlan pub-vlan]
	Type *string `json:"type"`
}

// Validate validates this network create
func (m *NetworkCreate) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateIPAddressRanges(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *NetworkCreate) validateIPAddressRanges(formats strfmt.Registry) error {

	if swag.IsZero(m.IPAddressRanges) { // not required
		return nil
	}

	for i := 0; i < len(m.IPAddressRanges); i++ {
		if swag.IsZero(m.IPAddressRanges[i]) { // not required
			continue
		}

		if m.IPAddressRanges[i] != nil {
			if err := m.IPAddressRanges[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("ipAddressRanges" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

var networkCreateTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["vlan","pub-vlan"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		networkCreateTypeTypePropEnum = append(networkCreateTypeTypePropEnum, v)
	}
}

const (

	// NetworkCreateTypeVlan captures enum value "vlan"
	NetworkCreateTypeVlan string = "vlan"

	// NetworkCreateTypePubVlan captures enum value "pub-vlan"
	NetworkCreateTypePubVlan string = "pub-vlan"
)

// prop value enum
func (m *NetworkCreate) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, networkCreateTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *NetworkCreate) validateType(formats strfmt.Registry) error {

	if err := validate.Required("type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", *m.Type); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *NetworkCreate) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NetworkCreate) UnmarshalBinary(b []byte) error {
	var res NetworkCreate
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
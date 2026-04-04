package models

import (
	"errors"
	"fmt"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

var errInvalidVenueDescriptionItemType = errors.New("invalid apiv1.VenueDescriptionItemType_name")

// VenueDescriptionItemType describes the visual treatment of a venue description item.
//
//nolint:recvcheck
type VenueDescriptionItemType int

// VenueDescriptionItemType values.
const (
	VenueDescriptionItemTypeUnspecified VenueDescriptionItemType = iota
	VenueDescriptionItemTypeDefault
	VenueDescriptionItemTypeWarning
	VenueDescriptionItemTypeRule
)

// ProtoEnum returns an enum compatible with the generated proto code.
func (t *VenueDescriptionItemType) ProtoEnum() (apiv1.VenueDescriptionItemType, error) {
	if t == nil {
		return apiv1.VenueDescriptionItemType(0), fmt.Errorf("convert nil VenueDescriptionItemType to protobuf enum: %w", errInvalidVenueDescriptionItemType)
	}

	if val, ok := apiv1.VenueDescriptionItemType_value[t.String()]; ok {
		return apiv1.VenueDescriptionItemType(val), nil
	}

	return apiv1.VenueDescriptionItemType(0), fmt.Errorf("convert VenueDescriptionItemType(%q) to protobuf enum: %w", t.String(), errInvalidVenueDescriptionItemType)
}

// FromProtoEnum parses an enum from a proto message into an internal representation.
func (t *VenueDescriptionItemType) FromProtoEnum(pe apiv1.VenueDescriptionItemType) error {
	var err error

	*t, err = VenueDescriptionItemTypeString(pe.String())

	return err
}

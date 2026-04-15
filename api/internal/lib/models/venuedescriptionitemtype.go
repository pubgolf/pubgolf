package models

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

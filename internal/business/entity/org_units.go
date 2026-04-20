package entity

import "github.com/google/uuid"

const (
	ProfileTypeBrand = "BRAND"
	ProfileTypeVenue = "VENUE"
)

type OrgUnit struct {
	Id           int
	OrgAccountId uuid.UUID
	ProfileType  string
	Name         string
	Avatar       *string
	Banner       *string
	Description  *string
	ParentId     *int
	Latitude     *float32
	Longitude    *float32
}

type UpdateOrgUnitInput struct {
	Name        *string
	Avatar      *string
	Banner      *string
	Description *string
	Latitude    *float32
	Longitude   *float32
}

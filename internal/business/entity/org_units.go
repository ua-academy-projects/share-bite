package entity

import "github.com/google/uuid"

const (
	ProfileTypeBrand = "BRAND"
	ProfileTypeVenue = "VENUE"
)

type OrgStatus string

const (
	OrgStatusPending  OrgStatus = "pending"
	OrgStatusVerified OrgStatus = "verified"
	OrgStatusRejected OrgStatus = "rejected"
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
	Tags         []string
	Status       OrgStatus
}

type UpdateOrgUnitInput struct {
	Name        *string
	Avatar      *string
	Banner      *string
	Description *string
	Latitude    *float32
	Longitude   *float32
}
type OrgUnitWithDistance struct {
	OrgUnit  OrgUnit
	Distance float64
}

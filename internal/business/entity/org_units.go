package entity

import "github.com/google/uuid"

type OrgUnit struct {
	Id           int
	OrgAccountId uuid.UUID

	ProfileType string
	Name        string
	Avatar      string
	Banner      string
	Description string

	ParentId int

	Latitude  float32
	Longitude float32
}

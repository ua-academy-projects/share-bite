package business

import (
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

type OrgUnit struct {
	Id           int       `db:"id"`
	OrgAccountId uuid.UUID `db:"org_account_id"`
	ProfileType  string    `db:"profile_type"`
	Name         string    `db:"name"`
	Avatar       *string   `db:"avatar"`
	Banner       *string   `db:"banner"`
	Description  *string   `db:"description"`
	ParentId     *int      `db:"parent_id"`
	Latitude     *float32  `db:"latitude"`
	Longitude    *float32  `db:"longitude"`
}

func (e OrgUnit) ToEntity() entity.OrgUnit {
	ou := entity.OrgUnit{
		Id:           e.Id,
		OrgAccountId: e.OrgAccountId,
		ProfileType:  e.ProfileType,
		Name:         e.Name,
	}

	if e.Avatar != nil {
		ou.Avatar = *e.Avatar
	}
	if e.Banner != nil {
		ou.Banner = *e.Banner
	}
	if e.Description != nil {
		ou.Description = *e.Description
	}
	if e.ParentId != nil {
		ou.ParentId = *e.ParentId
	}
	if e.Latitude != nil {
		ou.Latitude = *e.Latitude
	}
	if e.Longitude != nil {
		ou.Longitude = *e.Longitude
	}

	return ou
}

func executeSQLError(err error) error {
	return errwrap.Wrap("execute sql", err)
}

func scanRowError(err error) error {
	return errwrap.Wrap("scan row", err)
}

func scanRowsError(err error) error {
	return errwrap.Wrap("scan rows", err)
}

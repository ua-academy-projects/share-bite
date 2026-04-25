package business

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
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
	return entity.OrgUnit{
		Id:           e.Id,
		OrgAccountId: e.OrgAccountId,
		ProfileType:  e.ProfileType,
		Name:         e.Name,
		Avatar:       e.Avatar,
		Banner:       e.Banner,
		Description:  e.Description,
		ParentId:     e.ParentId,
		Latitude:     e.Latitude,
		Longitude:    e.Longitude,
	}
}

func executeSQLError(err error) error {
	return fmt.Errorf("execute sql: %w", err)
}

func scanRowError(err error) error {
	return fmt.Errorf("scan row: %w", err)
}

func scanRowsError(err error) error {
	return fmt.Errorf("scan rows: %w", err)
}

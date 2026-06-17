package entity

import (
	"time"

	"github.com/google/uuid"
)

type VerificationLog struct {
	ID        int64
	OrgUnitID int
	AdminID   uuid.UUID
	OldStatus *OrgStatus
	NewStatus OrgStatus
	Comment   *string
	CreatedAt time.Time
}

package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	repository "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

const (
	minBusinessNameLength = 3
	maxBusinessNameLength = 40
)

type businessRepository interface {
	Create(ctx context.Context, in entity.OrgUnit) (int, error)
	GetById(ctx context.Context, id int) (*entity.OrgUnit, error)
	UpdateOrg(ctx context.Context, id int, orgAccountID uuid.UUID, in entity.UpdateOrgUnitInput) (*entity.OrgUnit, error)
	DeleteOrg(ctx context.Context, id int, orgAccountID uuid.UUID) error
	UpdatePost(ctx context.Context, postID int64, orgID int, content string) (*entity.Post, error)
	DeletePost(ctx context.Context, id int64, orgID int) error
	GetOrgIDByUserID(ctx context.Context, userID int64) (int, error)
}

type service struct {
	businessRepo businessRepository
}

func New(businessRepo businessRepository) *service {
	return &service{
		businessRepo: businessRepo,
	}
}

func (s *service) Create(ctx context.Context, in entity.OrgUnit) (int, error) {

	nameLen := len([]rune(in.Name))
	if nameLen < minBusinessNameLength {
		return 0, apperror.BadRequest("business name cannot be less than 3 characters long")
	}
	if nameLen > maxBusinessNameLength {
		return 0, apperror.BadRequest("business name cannot be more than 40 characters long")
	}


	if in.ProfileType == "" {
		return 0, apperror.BadRequest("business type is required")
	}
	if in.ProfileType != entity.ProfileTypeBrand && in.ProfileType != entity.ProfileTypeVenue {
		return 0, apperror.BadRequest("invalid business type")
	}


	if in.ProfileType == entity.ProfileTypeBrand && in.ParentId != nil {
		return 0, apperror.BadRequest("BRAND cannot have a parent")
	}
	if in.ProfileType == entity.ProfileTypeVenue && in.ParentId == nil {
		return 0, apperror.BadRequest("VENUE must have a parent_id")
	}


	if in.ProfileType == entity.ProfileTypeVenue {
		parent, err := s.businessRepo.GetById(ctx, *in.ParentId)
		if err != nil {
			return 0, apperror.BadRequest("parent_id does not exist")
		}
		if parent.ProfileType != entity.ProfileTypeBrand {
			return 0, apperror.BadRequest("parent must be a BRAND, not a VENUE")
		}
	}

	id, err := s.businessRepo.Create(ctx, in)
	if err != nil {
		return 0, fmt.Errorf("failed to create business profile: %w", err)
	}
	return id, nil
}

func (s *service) Get(ctx context.Context, id int) (*entity.OrgUnit, error) {
	orgUnit, err := s.businessRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.OrgUnitNotFoundID(id)
		}
		return nil, errwrap.Wrap("get org unit from business repository", err)
	}

	return orgUnit, nil
}


func (s *service) UpdateOrg(ctx context.Context, id int, orgAccountID uuid.UUID, in entity.UpdateOrgUnitInput) (*entity.OrgUnit, error) {
	if in.Name == nil &&
		in.Avatar == nil &&
		in.Banner == nil &&
		in.Description == nil &&
		in.Latitude == nil &&
		in.Longitude == nil {
		return nil, apperror.BadRequest("at least one updatable field is required")
	}

	if in.Name != nil {
		nameLen := len([]rune(*in.Name))
		if nameLen < minBusinessNameLength {
			return nil, apperror.BadRequest("business name cannot be less than 3 characters long")
		}
		if nameLen > maxBusinessNameLength {
			return nil, apperror.BadRequest("business name cannot be more than 40 characters long")
		}
	}

	updated, err := s.businessRepo.UpdateOrg(ctx, id, orgAccountID, in)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.OrgUnitNotFoundID(id)
		}
		return nil, errwrap.Wrap("update org unit in business repository", err)
	}

	return updated, nil
}

func (s *service) DeleteOrg(ctx context.Context, id int, orgAccountID uuid.UUID) error {
	err := s.businessRepo.DeleteOrg(ctx, id, orgAccountID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.OrgUnitNotFoundID(id)
		}
		return errwrap.Wrap("delete org unit in business repository", err)
	}

	return nil
}


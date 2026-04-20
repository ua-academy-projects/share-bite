package business

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"

	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	repository "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

const (
	minBusinessNameLength = 3
	maxBusinessNameLength = 40
)

type businessRepository interface {
	Create(ctx context.Context, in entity.OrgUnit) (int, error)
	UpdateOrg(ctx context.Context, id int, orgAccountID uuid.UUID, in entity.UpdateOrgUnitInput) (*entity.OrgUnit, error)
	DeleteOrg(ctx context.Context, id int, orgAccountID uuid.UUID) error

	UpdatePost(ctx context.Context, postID int64, orgID int, content string) (*entity.Post, error)
	DeletePost(ctx context.Context, id int64, orgID int) error
	GetOrgIDByUserID(ctx context.Context, userID string) (int, error)
	GetPostPhotos(ctx context.Context, postID int64) ([]string, error)
	CheckOwnership(ctx context.Context, userID string, unitID int) error
	CreatePost(ctx context.Context, userID string, unitID int, description string) (*entity.Post, error)
	InsertPostImages(ctx context.Context, postID int64, URLs []string) error
	GetPosts(ctx context.Context, skip, limit int) (pagination.Result[entity.Post], error)
	GetPostByID(ctx context.Context, postID int64) (*entity.Post, error)

	CreateBox(ctx context.Context, box *entity.Box) (int64, time.Time, error)
	CreateBoxItem(ctx context.Context, boxID int64, code string) error
	GetBrandIDByOwnerUserID(ctx context.Context, userID string) (int, error)
	CreateLocation(ctx context.Context, brandID int, ownerUserID string, in dto.CreateLocationInput) (*entity.OrgUnit, error)
	UpdateLocation(ctx context.Context, locationID int, brandID int, in dto.UpdateLocationInput) (*entity.OrgUnit, error)
	DeleteLocation(ctx context.Context, locationID int, brandID int) error
	ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int) (pagination.Result[entity.BoxWithDistance], error)

	GetById(ctx context.Context, id int) (*entity.OrgUnit, error)
	ListByParentID(ctx context.Context, parentID, offset, limit int) (pagination.Result[entity.OrgUnit], error)
	GetVenuesByIDs(ctx context.Context, ids []int) ([]entity.OrgUnit, error)
	GetVenueRating(ctx context.Context, venueID int) (float32, error)
}

type service struct {
	businessRepo businessRepository
	txManager    database.TxManager
	storage      storage.ObjectStorage
}

func New(businessRepo businessRepository, txManager database.TxManager, st storage.ObjectStorage) *service {
	return &service{
		businessRepo: businessRepo,
		txManager:    txManager,
		storage:      st,
	}
}

func generateCode() string {
	return uuid.New().String()[:12]
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
		return nil, fmt.Errorf("update org unit in business repository: %w", err)
	}

	return updated, nil
}

func (s *service) DeleteOrg(ctx context.Context, id int, orgAccountID uuid.UUID) error {
	err := s.businessRepo.DeleteOrg(ctx, id, orgAccountID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.OrgUnitNotFoundID(id)
		}
		return fmt.Errorf("delete org unit in business repository: %w", err)
	}

	return nil
}

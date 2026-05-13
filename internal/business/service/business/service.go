package business

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/storage/key"
	"github.com/ua-academy-projects/share-bite/internal/storage/mediatype"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/logger"

	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	repository "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

const (
	minBusinessNameLength = 3
	maxBusinessNameLength = 40
	fileSniffSizeBytes    = 512
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

	CreateLike(ctx context.Context, postID int64, authorID string) (*entity.Like, error)
	DeleteLike(ctx context.Context, postID int64, authorID string) error
	CheckUserLiked(ctx context.Context, postID int64, authorID string) (bool, error)
	CountLikesByPost(ctx context.Context, postID int64) (int, error)
	GetLikesByPost(ctx context.Context, postID int64, limit, offset int) ([]entity.LikeWithAuthor, error)
	CreateComment(ctx context.Context, postID int64, authorID, content string) (*entity.Comment, error)
	GetCommentByID(ctx context.Context, commentID int64) (*entity.Comment, error)
	UpdateComment(ctx context.Context, commentID int64, content string) (*entity.Comment, error)
	DeleteComment(ctx context.Context, commentID int64) error
	ListCommentsWithAuthorsByPost(ctx context.Context, postID int64, limit, offset int) ([]entity.CommentWithAuthor, error)
	CountCommentsByPost(ctx context.Context, postID int64) (int, error)
	CreateBox(ctx context.Context, box *entity.Box) (int64, time.Time, error)
	CreateBoxItem(ctx context.Context, boxID int64, code string) error
	GetBrandIDByOwnerUserID(ctx context.Context, userID string) (int, error)
	CreateLocation(ctx context.Context, brandID int, ownerUserID string, in dto.CreateLocationInput) (*entity.OrgUnit, error)
	UpdateLocation(ctx context.Context, locationID int, brandID int, in dto.UpdateLocationInput) (*entity.OrgUnit, error)
	DeleteLocation(ctx context.Context, locationID int, brandID int) error
	ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int, orgID *int) (pagination.Result[entity.BoxWithDistance], error)
	GetOrgUnitTagSlugs(ctx context.Context, orgUnitID int) ([]string, error)
	GetOrgUnitTagsByOrgUnitID(ctx context.Context, ids []int) (map[int][]string, error)
	SetOrgUnitTagsByIDs(ctx context.Context, orgUnitID int, tagIDs []int) error
	ListLocationTags(ctx context.Context) ([]entity.LocationTag, error)

	GetBox(ctx context.Context, boxID int64) (*entity.Box, error)
	ReserveBoxItem(ctx context.Context, boxID int64, userID string) (string, error)

	GetById(ctx context.Context, id int) (*entity.OrgUnit, error)
	ListByParentID(ctx context.Context, parentID, offset, limit int, tags []string) (pagination.Result[entity.OrgUnit], error)
	GetVenuesByIDs(ctx context.Context, ids []int) ([]entity.OrgUnit, error)
	GetVenueRating(ctx context.Context, venueID int) (float32, error)
	ListNearbyVenues(ctx context.Context, lat, lon float64, offset, limit int) (pagination.Result[entity.OrgUnitWithDistance], error)
	SearchVenues(ctx context.Context, query string, offset, limit int, tags []string) (pagination.Result[entity.OrgUnit], error)
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

	if in.ProfileType != entity.ProfileTypeBrand {
		return 0, apperror.BadRequest("only BRAND creation is allowed via this service")
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
		in.Description == nil {
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

func (s *service) UploadAvatar(ctx context.Context, id int, orgAccountID uuid.UUID, fileHeader *multipart.FileHeader) (*entity.OrgUnit, error) {
	return s.uploadOrgAsset(ctx, id, orgAccountID, fileHeader, true)
}

func (s *service) UploadBanner(ctx context.Context, id int, orgAccountID uuid.UUID, fileHeader *multipart.FileHeader) (*entity.OrgUnit, error) {
	return s.uploadOrgAsset(ctx, id, orgAccountID, fileHeader, false)
}

func (s *service) uploadOrgAsset(ctx context.Context, id int, orgAccountID uuid.UUID, fileHeader *multipart.FileHeader, isAvatar bool) (*entity.OrgUnit, error) {
	if s.storage == nil {
		return nil, apperror.Internal("storage is not configured")
	}

	current, err := s.businessRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.OrgUnitNotFoundID(id)
		}
		return nil, fmt.Errorf("get current org unit: %w", err)
	}

	if current.ProfileType != entity.ProfileTypeBrand {
		return nil, apperror.BadRequest("only brand profiles can have avatars/banners uploaded through this endpoint")
	}

	if current.OrgAccountId != orgAccountID {
		return nil, apperror.Forbidden("you do not have permission to update this brand")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, apperror.Internal(fmt.Sprintf("open uploaded file: %v", err))
	}
	defer file.Close()

	buffer := make([]byte, fileSniffSizeBytes)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, apperror.Internal(fmt.Sprintf("read file header: %v", err))
	}

	contentType := http.DetectContentType(buffer[:n])
	if err := mediatype.DefaultImageValidator.Validate(contentType, fileHeader.Size); err != nil {
		return nil, apperror.BadRequest(err.Error())
	}

	seeker, ok := file.(io.Seeker)
	if !ok {
		return nil, apperror.Internal("uploaded file is not seekable")
	}
	if _, err := seeker.Seek(0, io.SeekStart); err != nil {
		return nil, apperror.Internal(fmt.Sprintf("seek to start of file: %v", err))
	}

	ext, ok := mediatype.ExtFromContentType(contentType)
	if !ok {
		return nil, apperror.BadRequest("unsupported image type")
	}

	var objectKey string
	if isAvatar {
		objectKey = key.BusinessAvatarKey(id, uuid.NewString(), ext)
	} else {
		objectKey = key.BusinessBannerKey(id, uuid.NewString(), ext)
	}

	if err := s.storage.Upload(ctx, objectKey, contentType, file); err != nil {
		return nil, apperror.Internal(fmt.Sprintf("upload to storage: %v", err))
	}

	objectURL := s.storage.BuildURL(objectKey)
	input := entity.UpdateOrgUnitInput{}
	if isAvatar {
		input.Avatar = &objectURL
	} else {
		input.Banner = &objectURL
	}

	updated, err := s.businessRepo.UpdateOrg(ctx, id, orgAccountID, input)
	if err != nil {
		go s.cleanupDelete(objectKey)
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.OrgUnitNotFoundID(id)
		}
		return nil, apperror.Internal(fmt.Sprintf("update org in repository: %v", err))
	}

	var oldURL *string
	if isAvatar {
		oldURL = current.Avatar
	} else {
		oldURL = current.Banner
	}

	if oldURL != nil {
		if oldKey := s.keyFromURL(*oldURL); oldKey != "" && oldKey != objectKey {
			go s.cleanupDelete(oldKey)
		}
	}

	return updated, nil
}

func (s *service) cleanupDelete(key string) {
	if s.storage == nil || key == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.storage.Delete(ctx, key); err != nil {
		logger.ErrorKV(ctx, "failed to cleanup business asset",
			"key", key,
			"error", err,
		)
	}
}

func (s *service) keyFromURL(url string) string {
	const prefix = "businesses/"
	idx := strings.Index(url, prefix)
	if idx == -1 {
		return ""
	}
	return url[idx:]
}

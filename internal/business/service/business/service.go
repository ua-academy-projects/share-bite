package business

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"

	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	repository "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/aws"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

const (
	minBusinessNameLength = 3
	maxBusinessNameLength = 40
	orgUnitEntityType     = "org_unit"
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
	GetFirstVenueIDByOwnerUserID(ctx context.Context, userID string) (int, error)
	CreateLocation(ctx context.Context, brandID int, ownerUserID string, in dto.CreateLocationInput) (*entity.OrgUnit, error)
	UpdateLocation(ctx context.Context, locationID int, brandID int, in dto.UpdateLocationInput, h3Hash *string) (*entity.OrgUnit, error)
	DeleteLocation(ctx context.Context, locationID int, brandID int) error
	ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int, orgID *int) (pagination.Result[entity.BoxWithDistance], error)
	GetOrgUnitTagSlugs(ctx context.Context, orgUnitID int) ([]string, error)
	GetOrgUnitTagsByOrgUnitID(ctx context.Context, ids []int) (map[int][]string, error)
	SetOrgUnitTagsByIDs(ctx context.Context, orgUnitID int, tagIDs []int) error
	ListLocationTags(ctx context.Context) ([]entity.LocationTag, error)
	ReplaceLocationHours(ctx context.Context, venueID int, days []dto.VenueHoursDayInput) error
	GetLocationHours(ctx context.Context, venueID int) ([]dto.VenueHoursDayInput, error)

	GetBox(ctx context.Context, boxID int64) (*entity.Box, error)
	ReserveBoxItem(ctx context.Context, boxID int64, userID string) (string, error)
	ListBoxesByVenueID(ctx context.Context, venueID int, offset, limit int) (pagination.Result[entity.Box], error)
	UpdateBox(ctx context.Context, boxID int64, input entity.BoxUpdateInput) (*entity.Box, error)
	GetBoxItems(ctx context.Context, boxID int64, offset, limit int) (pagination.Result[entity.BoxItem], error)

	GetById(ctx context.Context, id int) (*entity.OrgUnit, error)
	ListByParentID(ctx context.Context, parentID, offset, limit int, tags []string) (pagination.Result[entity.OrgUnit], error)
	GetVenuesByIDs(ctx context.Context, ids []int) ([]entity.OrgUnit, error)
	GetVenueRating(ctx context.Context, venueID int) (float32, error)
	ListNearbyVenues(ctx context.Context, lat, lon float64, offset, limit int) (pagination.Result[entity.OrgUnitWithDistance], error)
	SearchVenues(ctx context.Context, query string, offset, limit int, tags []string) (pagination.Result[entity.OrgUnit], error)
	ResubmitVerification(ctx context.Context, id int, userID string) error

	GetTopTagsByUserLikes(ctx context.Context, userID string, tagsToFetch int) ([]string, error)
	GetPostsByTag(ctx context.Context, tag string, quota int, seenCompositeIDs []string, h3Hashes []string) ([]entity.RecommendedPost, error)
	GetRandomPosts(ctx context.Context, deficit int, seenCompositeIDs []string, h3Hashes []string) ([]entity.RecommendedPost, error)
	CountPostsByTag(ctx context.Context, tag string, h3Hashes []string) (int, error)
	CountRandomPosts(ctx context.Context, h3Hashes []string) (int, error)

	GetDailySummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID) (entity.DailySummary, error)
	GetReservationSummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID *int) (entity.ReservationSummary, error)
	GetVenueActivitySummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID int) (entity.VenueActivitySummary, error)
	GetFoodBoxPerformance(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID *int) (entity.BoxPerformanceRaw, error)
	GetEngagementSummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID *int) (entity.EngagementSummaryRaw, error)
}

type H3Settings struct {
	Resolution      int
	RecommendRadius int
}

type service struct {
	businessRepo businessRepository
	txManager    database.TxManager
	storage      storage.ObjectStorage
	h3Service    aws.H3Service
	h3Config     H3Settings
	outboxWriter outbox.Writer
	emailClient  interface {
		GetUserEmail(ctx context.Context, userID, authToken string) (string, error)
	}
}

func New(businessRepo businessRepository, txManager database.TxManager, st storage.ObjectStorage, h3Service aws.H3Service, h3Config H3Settings, outboxWriter outbox.Writer, emailClient interface {
	GetUserEmail(ctx context.Context, userID, authToken string) (string, error)
}) *service {
	return &service{
		businessRepo: businessRepo,
		txManager:    txManager,
		storage:      st,
		h3Service:    h3Service,
		h3Config:     h3Config,
		outboxWriter: outboxWriter,
		emailClient:  emailClient,
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

	authToken, ok := ctx.Value(middleware.CtxAccessToken).(string)
	if !ok || authToken == "" {
		return 0, fmt.Errorf("missing access token to resolve business email")
	}

	email, err := s.emailClient.GetUserEmail(ctx, in.OrgAccountId.String(), authToken)
	if err != nil {
		return 0, fmt.Errorf("get business email: %w", err)
	}

	var id int
	err = s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		var err error
		id, err = s.businessRepo.Create(txCtx, in)
		if err != nil {
			return fmt.Errorf("failed to create business profile: %w", err)
		}

		metadata := map[string]any{
			"username": in.Name,
			"email":    email,
		}

		event := outbox.Event{
			EventType: outbox.EventTypeRegistrationConfirmed,
			Payload: outbox.Message{
				EventID:     outbox.NewEventID(strconv.Itoa(id), email),
				EventType:   outbox.EventTypeRegistrationConfirmed,
				RecipientID: in.OrgAccountId.String(),
				ActorID:     in.OrgAccountId.String(),
				EntityType:  orgUnitEntityType,
				EntityID:    strconv.Itoa(id),
				Metadata:    metadata,
				CreatedAt:   time.Now().UTC(),
			},
			SourceService: outbox.DefaultSourceService,
		}

		if err := s.outboxWriter.Enqueue(txCtx, event); err != nil {
			return fmt.Errorf("failed to enqueue registration_confirmed outbox event: %w", err)
		}

		return nil
	})

	if err != nil {
		return 0, err
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

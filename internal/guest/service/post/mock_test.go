package post

import (
	"context"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
	"io"
	"strings"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type postRepositoryMock struct {
	createFn                            func(ctx context.Context, in dto.CreatePostInput) (entity.Post, error)
	listFn                              func(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error)
	getFn                               func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error)
	getByIDFn                           func(ctx context.Context, postID string) (entity.Post, error)
	getAuthorUserIDFn                   func(ctx context.Context, postID string) (string, error)
	getAuthorCustomerIDFn               func(ctx context.Context, postID string) (string, error)
	deleteImagesByPostIDReturningKeysFn func(ctx context.Context, postID string) ([]string, error)
	updateFn                            func(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error)
	likeFn                              func(ctx context.Context, postID string, customerID string) (bool, error)
	unlikeFn                            func(ctx context.Context, postID string, customerID string) error
	updateStatusFn                      func(ctx context.Context, postID, customerID string, status entity.PostStatus) error
	getPostsByVenueIDsFn                func(ctx context.Context, venueIDs []int64, limit int) ([]entity.Post, error)
	createMentionsFn                    func(ctx context.Context, mentions []entity.PostMention) error
	listMentionsFn                      func(ctx context.Context, postIDs []string) (map[string][]entity.PostMention, error)

	createPostCollaboratorsFn      func(ctx context.Context, postID string, invitedBy string, customerIDs []string, expiresAt time.Time) error
	getPendingPostInvitationsFn    func(ctx context.Context, customerID string) ([]entity.PostCollaborator, error)
	acceptPostInvitationFn         func(ctx context.Context, collaboratorID string, customerID string) (string, error)
	declinePostInvitationFn        func(ctx context.Context, collaboratorID string, customerID string) error
	getAcceptedPostCollaboratorsFn func(ctx context.Context, postID string) ([]string, error)
	tryPublishPostIfAllAcceptedFn  func(ctx context.Context, postID string) (bool, error)
	isAcceptedCollaboratorFn       func(ctx context.Context, postID string, customerID string) (bool, error)

	deleteExpiredDraftPostsFn func(ctx context.Context) error

	lastCreateInput        dto.CreatePostInput
	lastListInput          dto.ListPostsInput
	lastGetID              string
	lastGetViewerID        string
	lastUpdateInput        entity.UpdatePostInput
	lastDeleteImagesPostID string
}

func (m *postRepositoryMock) CreateImages(ctx context.Context, images []entity.PostImage) ([]entity.PostImage, error) {
	return images, nil
}

func (m *postRepositoryMock) UpdateProcessedMetadata(ctx context.Context, imageID string, thumbnailKey string, width int, height int) error {
	return nil
}

func (m *postRepositoryMock) MarkProcessingFailed(ctx context.Context, imageID string, reason string) error {
	return nil
}

func (m *postRepositoryMock) ClaimForProcessing(ctx context.Context, imageID string) (bool, error) {
	return false, nil
}

func (m *postRepositoryMock) DeleteImagesByPostIDReturningKeys(ctx context.Context, postID string) ([]string, error) {
	m.lastDeleteImagesPostID = postID

	if m.deleteImagesByPostIDReturningKeysFn != nil {
		return m.deleteImagesByPostIDReturningKeysFn(
			ctx,
			postID,
		)
	}

	return nil, nil
}

func (m *postRepositoryMock) Create(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
	m.lastCreateInput = in
	if m.createFn != nil {
		return m.createFn(ctx, in)
	}
	return entity.Post{}, nil
}

func (m *postRepositoryMock) Update(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error) {
	m.lastUpdateInput = in
	if m.updateFn != nil {
		return m.updateFn(ctx, in)
	}
	return entity.Post{ID: in.ID}, nil
}

func (m *postRepositoryMock) List(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error) {
	m.lastListInput = in
	if m.listFn != nil {
		return m.listFn(ctx, in)
	}
	return dto.ListPostsOutput{}, nil
}

func (m *postRepositoryMock) Get(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
	m.lastGetID = postID
	m.lastGetViewerID = reqCustomerID
	if m.getFn != nil {
		return m.getFn(ctx, postID, reqCustomerID)
	}
	return entity.Post{ID: postID}, nil
}

func (m *postRepositoryMock) GetByID(ctx context.Context, postID string) (entity.Post, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, postID)
	}
	return entity.Post{ID: postID}, nil
}

func (m *postRepositoryMock) GetAuthorUserID(ctx context.Context, postID string) (string, error) {
	if m.getAuthorUserIDFn != nil {
		return m.getAuthorUserIDFn(ctx, postID)
	}
	return "", nil
}

func (m *postRepositoryMock) DeleteImagesByPostID(ctx context.Context, postID string) error {
	return nil
}

func (m *postRepositoryMock) Like(ctx context.Context, postID string, customerID string) (bool, error) {
	if m.likeFn != nil {
		return m.likeFn(ctx, postID, customerID)
	}
	return true, nil
}

func (m *postRepositoryMock) Unlike(ctx context.Context, postID string, customerID string) error {
	if m.unlikeFn != nil {
		return m.unlikeFn(ctx, postID, customerID)
	}
	return nil
}

func (m *postRepositoryMock) UpdateStatus(ctx context.Context, postID, customerID string, status entity.PostStatus) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, postID, customerID, status)
	}
	return nil
}

func (m *postRepositoryMock) GetPostsByVenueIDs(ctx context.Context, venueIDs []int64, limit int) ([]entity.Post, error) {
	if m.getPostsByVenueIDsFn != nil {
		return m.getPostsByVenueIDsFn(ctx, venueIDs, limit)
	}
	return []entity.Post{}, nil
}

type venueProviderMock struct {
	checkExistsFn     func(ctx context.Context, venueID int64) (bool, error)
	getNearbyVenuesFn func(ctx context.Context, lat, lon float64, limit int) ([]int64, error)
}

func (m *venueProviderMock) CheckExists(ctx context.Context, venueID int64) (bool, error) {
	if m.checkExistsFn != nil {
		return m.checkExistsFn(ctx, venueID)
	}
	return true, nil
}

func (m *venueProviderMock) GetNearbyVenues(ctx context.Context, lat, lon float64, limit int) ([]int64, error) {
	if m.getNearbyVenuesFn != nil {
		return m.getNearbyVenuesFn(ctx, lat, lon, limit)
	}
	return []int64{}, nil
}

type followRepoMock struct {
	getAllowedMentionsFn func(ctx context.Context, customerID string, ids []string) ([]string, error)
	getFollowersFn       func(ctx context.Context, customerID string) ([]entity.Customer, error)
}

func (m *followRepoMock) GetAllowedMentions(ctx context.Context, customerID string, ids []string) ([]string, error) {
	if m.getAllowedMentionsFn != nil {
		return m.getAllowedMentionsFn(ctx, customerID, ids)
	}
	return ids, nil
}

func (m *followRepoMock) GetFollowers(ctx context.Context, customerID string) ([]entity.Customer, error) {
	if m.getFollowersFn != nil {
		return m.getFollowersFn(ctx, customerID)
	}
	return nil, nil
}

type customerRepoMock struct {
	getByIDsFn func(ctx context.Context, customerIDs []string) ([]entity.Customer, error)
	getByIDFn  func(ctx context.Context, customerID string) (entity.Customer, error)
}

func (m *customerRepoMock) GetByIDs(ctx context.Context, customerIDs []string) ([]entity.Customer, error) {
	if m.getByIDsFn != nil {
		return m.getByIDsFn(ctx, customerIDs)
	}
	customers := make([]entity.Customer, len(customerIDs))
	for i, id := range customerIDs {
		customers[i] = entity.Customer{ID: id}
	}
	return customers, nil
}

func (m *customerRepoMock) GetByID(ctx context.Context, customerID string) (entity.Customer, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, customerID)
	}

	return entity.Customer{ID: customerID}, nil
}

type txManagerMock struct {
	readCommittedFn func(ctx context.Context, fn database.Handler) error
}

func (m *txManagerMock) ReadCommitted(ctx context.Context, fn database.Handler) error {
	if m.readCommittedFn != nil {
		return m.readCommittedFn(ctx, fn)
	}
	return fn(ctx)
}

type publisherMock struct {
	publishFn func(ctx context.Context, target string, msg notification.Message) error
}

func (m *publisherMock) Publish(ctx context.Context, target string, msg notification.Message) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, target, msg)
	}
	return nil
}

type storageMock struct {
	uploadFn          func(ctx context.Context, key string, contentType string, file io.Reader) error
	deleteFn          func(ctx context.Context, key string) error
	buildURLFn        func(key string) string
	getPresignedURLFn func(ctx context.Context, key string) (string, error)
	getFn             func(ctx context.Context, key string) (io.ReadCloser, error)
}

func (m *storageMock) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	if m.getFn != nil {
		return m.getFn(ctx, key)
	}

	return io.NopCloser(strings.NewReader("")), nil
}

func (m *storageMock) Upload(ctx context.Context, key string, contentType string, file io.Reader) error {
	if m.uploadFn != nil {
		return m.uploadFn(ctx, key, contentType, file)
	}
	return nil
}

func (m *storageMock) Delete(ctx context.Context, key string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, key)
	}
	return nil
}

func (m *storageMock) BuildURL(key string) string {
	if m.buildURLFn != nil {
		return m.buildURLFn(key)
	}
	return key
}

func (m *storageMock) GetPresignedURL(ctx context.Context, key string) (string, error) {
	if m.getPresignedURLFn != nil {
		return m.getPresignedURLFn(ctx, key)
	}
	return "http://localhost:3900/app-dev-bucket/" + key + "?signed=true", nil
}

func (m *postRepositoryMock) CreateMentions(ctx context.Context, mentions []entity.PostMention) error {
	if m.createMentionsFn != nil {
		return m.createMentionsFn(ctx, mentions)
	}
	return nil
}

func (m *postRepositoryMock) ListMentionsByPostIDs(ctx context.Context, postIDs []string) (map[string][]entity.PostMention, error) {
	if m.listMentionsFn != nil {
		return m.listMentionsFn(ctx, postIDs)
	}
	return map[string][]entity.PostMention{}, nil
}

func (m *postRepositoryMock) CreatePostCollaborators(
	ctx context.Context,
	postID string,
	invitedBy string,
	customerIDs []string,
	expiresAt time.Time,
) error {
	if m.createPostCollaboratorsFn != nil {
		return m.createPostCollaboratorsFn(ctx, postID, invitedBy, customerIDs, expiresAt)
	}

	return nil
}

func (m *postRepositoryMock) GetPendingPostInvitations(
	ctx context.Context,
	customerID string,
) ([]entity.PostCollaborator, error) {
	if m.getPendingPostInvitationsFn != nil {
		return m.getPendingPostInvitationsFn(ctx, customerID)
	}

	return []entity.PostCollaborator{}, nil
}

func (m *postRepositoryMock) AcceptPostInvitation(ctx context.Context, collaboratorID string, customerID string) (string, error) {
	if m.acceptPostInvitationFn != nil {
		return m.acceptPostInvitationFn(
			ctx,
			collaboratorID,
			customerID,
		)
	}

	return "", nil
}

func (m *postRepositoryMock) DeclinePostInvitation(ctx context.Context, collaboratorID string, customerID string) error {
	if m.declinePostInvitationFn != nil {
		return m.declinePostInvitationFn(ctx, collaboratorID, customerID)
	}

	return nil
}

func (m *postRepositoryMock) GetAcceptedPostCollaborators(
	ctx context.Context,
	postID string,
) ([]string, error) {
	if m.getAcceptedPostCollaboratorsFn != nil {
		return m.getAcceptedPostCollaboratorsFn(ctx, postID)
	}

	return []string{}, nil
}

func (m *postRepositoryMock) DeleteExpiredDraftPosts(ctx context.Context) error {
	if m.deleteExpiredDraftPostsFn != nil {
		return m.deleteExpiredDraftPostsFn(ctx)
	}

	return nil
}

func (m *postRepositoryMock) TryPublishPostIfAllAccepted(ctx context.Context, postID string) (bool, error) {
	if m.tryPublishPostIfAllAcceptedFn != nil {
		return m.tryPublishPostIfAllAcceptedFn(ctx, postID)
	}

	return false, nil
}

func (m *postRepositoryMock) GetAuthorCustomerID(ctx context.Context, postID string) (string, error) {
	if m.getAuthorCustomerIDFn != nil {
		return m.getAuthorCustomerIDFn(ctx, postID)
	}

	return "", nil
}

func (m *postRepositoryMock) IsAcceptedCollaborator(ctx context.Context, postID string, customerID string) (bool, error) {
	if m.isAcceptedCollaboratorFn != nil {
		return m.isAcceptedCollaboratorFn(
			ctx,
			postID,
			customerID,
		)
	}
	return false, nil
}

type outboxWriterMock struct {
	enqueueFn func(ctx context.Context, event outbox.Event) error
}

func (m *outboxWriterMock) Enqueue(ctx context.Context, event outbox.Event) error {
	if m.enqueueFn != nil {
		return m.enqueueFn(ctx, event)
	}
	return nil
}

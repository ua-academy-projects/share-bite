package post

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"testing"
)

type MockPostRepo struct {
	mock.Mock
}

func (m *MockPostRepo) Create(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Post), args.Error(1)
}

func (m *MockPostRepo) CreateImages(ctx context.Context, images []entity.PostImage) error {
	args := m.Called(ctx, images)
	return args.Error(0)
}

func (m *MockPostRepo) GetByID(ctx context.Context, id string) (entity.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Post), args.Error(1)
}

func (m *MockPostRepo) Update(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Post), args.Error(1)
}

type MockVenueProvider struct {
	mock.Mock
}

func (m *MockVenueProvider) CheckExists(ctx context.Context, venueID int64) (bool, error) {
	args := m.Called(ctx, venueID)
	return args.Bool(0), args.Error(1)
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Upload(ctx context.Context, key, contentType string, file any) (string, error) {
	args := m.Called(ctx, key, contentType, file)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) ReadCommitted(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx)

	if args.Error(0) != nil {
		return args.Error(0)
	}

	return fn(ctx)
}

func newTestService() (*service, *MockPostRepo, *MockVenueProvider, *MockStorage, *MockTxManager) {
	repo := new(MockPostRepo)
	venue := new(MockVenueProvider)
	storage := new(MockStorage)
	tx := new(MockTxManager)

	svc := New(repo, venue, storage, tx)

	return svc, repo, venue, storage, tx
}

func TestService_Create_SuccessWithoutImages(t *testing.T) {
	svc, repo, venue, _, tx := newTestService()

	ctx := context.Background()

	input := dto.CreatePostInput{
		CustomerID: "user-1",
		VenueID:    1,
	}

	expectedPost := entity.Post{
		ID:         "post-1",
		CustomerID: "user-1",
	}

	venue.On("CheckExists", ctx, int64(1)).Return(true, nil)
	tx.On("ReadCommitted", ctx).Return(nil)

	repo.On("Create", ctx, input).Return(expectedPost, nil)

	post, err := svc.Create(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost.ID, post.ID)

	repo.AssertExpectations(t)
}

func TestService_Create_VenueNotFound(t *testing.T) {
	svc, _, venue, _, _ := newTestService()

	ctx := context.Background()

	input := dto.CreatePostInput{
		CustomerID: "user-1",
		VenueID:    1,
	}

	venue.On("CheckExists", ctx, int64(1)).Return(false, nil)

	_, err := svc.Create(ctx, input)

	assert.Error(t, err)
}

func TestService_Create_CreateImagesFailed_RollbackUploads(t *testing.T) {
	svc, repo, venue, storage, tx := newTestService()

	ctx := context.Background()

	input := dto.CreatePostInput{
		CustomerID: "user-1",
		VenueID:    1,
		Images: []dto.UploadImageInput{
			{
				ContentType: "image/jpeg",
				File:        nil,
				FileSize:    100,
			},
		},
	}

	createdPost := entity.Post{
		ID:         "post-1",
		CustomerID: "user-1",
	}

	venue.On("CheckExists", ctx, int64(1)).Return(true, nil)
	tx.On("ReadCommitted", ctx).Return(nil)

	repo.On("Create", ctx, input).Return(createdPost, nil)

	storage.On("Upload", ctx, mock.Anything, "image/jpeg", nil).
		Return("key-1", nil)

	repo.On("CreateImages", ctx, mock.Anything).
		Return(assert.AnError)

	// ожидаем rollback
	storage.On("Delete", mock.Anything, "key-1").Return(nil)

	_, err := svc.Create(ctx, input)

	assert.Error(t, err)

	storage.AssertExpectations(t)
}

func TestService_Update_NotOwner(t *testing.T) {
	svc, repo, _, _, _ := newTestService()

	ctx := context.Background()

	input := entity.UpdatePostInput{
		ID:         "post-1",
		CustomerID: "user-2", // не владелец
	}

	currentPost := entity.Post{
		ID:         "post-1",
		CustomerID: "user-1",
	}

	repo.On("GetByID", ctx, "post-1").Return(currentPost, nil)

	_, err := svc.Update(ctx, input)

	assert.Error(t, err)
}

func TestService_Update_InvalidStatusTransition(t *testing.T) {
	svc, repo, _, _, _ := newTestService()

	ctx := context.Background()

	newStatus := entity.PostStatusDraft

	input := entity.UpdatePostInput{
		ID:         "post-1",
		CustomerID: "user-1",
		Status:     &newStatus,
	}

	currentPost := entity.Post{
		ID:         "post-1",
		CustomerID: "user-1",
		Status:     entity.PostStatusPublished,
	}

	repo.On("GetByID", ctx, "post-1").Return(currentPost, nil)

	_, err := svc.Update(ctx, input)

	assert.Error(t, err)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestService_Update_RewriteImagesSuccess(t *testing.T) {
	svc, repo, _, storage, tx := newTestService()

	ctx := context.Background()

	input := entity.UpdatePostInput{
		ID:            "post-1",
		CustomerID:    "user-1",
		RewriteImages: true,
		Images: []dto.UploadImageInput{
			{
				ContentType: "image/jpeg",
				File:        nil,
				FileSize:    100,
			},
			{
				ContentType: "image/png",
				File:        nil,
				FileSize:    200,
			},
		},
	}

	currentPost := entity.Post{
		ID:         "post-1",
		CustomerID: "user-1",
		Status:     entity.PostStatusDraft,
		Images: []entity.PostImage{
			{ObjectKey: "old-key-1"},
			{ObjectKey: "old-key-2"},
		},
	}

	updatedPost := entity.Post{
		ID:         "post-1",
		CustomerID: "user-1",
		Status:     entity.PostStatusDraft,
	}

	repo.On("GetByID", ctx, "post-1").Return(currentPost, nil)
	tx.On("ReadCommitted", ctx).Return(nil)

	storage.On("Upload", ctx, mock.Anything, "image/jpeg", nil).Return("new-key-1", nil).Once()
	storage.On("Upload", ctx, mock.Anything, "image/png", nil).Return("new-key-2", nil).Once()

	repo.On("Update", ctx, input).Return(updatedPost, nil).Once()
	repo.On("DeleteImagesByPostID", ctx, "post-1").Return(nil).Once()

	repo.On("CreateImages", ctx, mock.MatchedBy(func(images []entity.PostImage) bool {
		if len(images) != 2 {
			return false
		}

		return images[0].PostID == "post-1" &&
			images[0].ObjectKey == "new-key-1" &&
			images[0].ContentType == "image/jpeg" &&
			images[0].FileSize == 100 &&
			images[0].SortOrder == 0 &&
			images[1].PostID == "post-1" &&
			images[1].ObjectKey == "new-key-2" &&
			images[1].ContentType == "image/png" &&
			images[1].FileSize == 200 &&
			images[1].SortOrder == 1
	})).Return(nil).Once()

	storage.On("Delete", mock.Anything, "old-key-1").Return(nil).Once()
	storage.On("Delete", mock.Anything, "old-key-2").Return(nil).Once()

	post, err := svc.Update(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, "post-1", post.ID)
	assert.Len(t, post.Images, 2)
	assert.Equal(t, "new-key-1", post.Images[0].ObjectKey)
	assert.Equal(t, "new-key-2", post.Images[1].ObjectKey)

	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
	tx.AssertExpectations(t)
}

func TestService_Update_RewriteImagesCreateImagesFailed_Rollback(t *testing.T) {
	svc, repo, _, storage, tx := newTestService()

	ctx := context.Background()

	input := entity.UpdatePostInput{
		ID:            "post-1",
		CustomerID:    "user-1",
		RewriteImages: true,
		Images: []dto.UploadImageInput{
			{
				ContentType: "image/jpeg",
				File:        nil,
				FileSize:    100,
			},
		},
	}

	currentPost := entity.Post{
		ID:         "post-1",
		CustomerID: "user-1",
		Status:     entity.PostStatusDraft,
		Images: []entity.PostImage{
			{ObjectKey: "old-key-1"},
		},
	}

	updatedPost := entity.Post{
		ID:         "post-1",
		CustomerID: "user-1",
		Status:     entity.PostStatusDraft,
	}

	repo.On("GetByID", ctx, "post-1").Return(currentPost, nil)
	tx.On("ReadCommitted", ctx).Return(nil)

	storage.On("Upload", ctx, mock.Anything, "image/jpeg", nil).Return("new-key-1", nil).Once()

	repo.On("Update", ctx, input).Return(updatedPost, nil).Once()
	repo.On("DeleteImagesByPostID", ctx, "post-1").Return(nil).Once()
	repo.On("CreateImages", ctx, mock.Anything).Return(assert.AnError).Once()

	storage.On("Delete", mock.Anything, "new-key-1").Return(nil).Once()

	_, err := svc.Update(ctx, input)

	assert.Error(t, err)

	storage.AssertCalled(t, "Delete", mock.Anything, "new-key-1")
	storage.AssertNotCalled(t, "Delete", mock.Anything, "old-key-1")
	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
	tx.AssertExpectations(t)
}

func TestService_Delete_AlreadyDeleted(t *testing.T) {
	svc, repo, _, _, _ := newTestService()

	ctx := context.Background()

	currentPost := entity.Post{
		ID:         "post-1",
		CustomerID: "user-1",
		Status:     entity.PostStatusDeleted,
	}

	repo.On("GetByID", ctx, "post-1").Return(currentPost, nil)

	err := svc.Delete(ctx, "post-1", "user-1")

	assert.NoError(t, err)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
}

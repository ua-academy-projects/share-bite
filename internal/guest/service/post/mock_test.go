package post

import (
	"context"
	"io"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type postRepositoryMock struct {
	createFn           func(ctx context.Context, in dto.CreatePostInput) (entity.Post, error)
	listFn             func(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error)
	getFn              func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error)
	getByIDFn          func(ctx context.Context, postID string) (entity.Post, error)
	getAuthorUserIDFn  func(ctx context.Context, postID string) (string, error)
	updateFn           func(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error)
	likeFn             func(ctx context.Context, postID string, customerID string) error
	unlikeFn           func(ctx context.Context, postID string, customerID string) error
	updateStatusFn     func(ctx context.Context, postID, customerID string, status entity.PostStatus) error

	lastCreateInput dto.CreatePostInput
	lastListInput   dto.ListPostsInput
	lastGetID       string
	lastGetViewerID string
	lastUpdateInput entity.UpdatePostInput
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

func (m *postRepositoryMock) CreateImages(ctx context.Context, images []entity.PostImage) error {
	return nil
}

func (m *postRepositoryMock) DeleteImagesByPostID(ctx context.Context, postID string) error {
	return nil
}

func (m *postRepositoryMock) Like(ctx context.Context, postID string, customerID string) error {
	if m.likeFn != nil {
		return m.likeFn(ctx, postID, customerID)
	}
	return nil
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

type venueProviderMock struct {
	checkExistsFn func(ctx context.Context, venueID int64) (bool, error)
}

func (m *venueProviderMock) CheckExists(ctx context.Context, venueID int64) (bool, error) {
	if m.checkExistsFn != nil {
		return m.checkExistsFn(ctx, venueID)
	}
	return true, nil
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
	uploadFn func(ctx context.Context, key string, contentType string, file io.Reader) (string, error)
	deleteFn   func(ctx context.Context, key string) error
	buildURLFn func(key string) string
}

func (m *storageMock) Upload(ctx context.Context, key string, contentType string, file io.Reader) (string, error) {
	if m.uploadFn != nil {
		return m.uploadFn(ctx, key, contentType, file)
	}
	return key, nil
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

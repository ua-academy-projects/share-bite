package customer

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

type mockCustomerRepository struct {
	mock.Mock
}

type mockOutboxWriter struct {
	mock.Mock
}

func (m *mockOutboxWriter) Enqueue(ctx context.Context, event outbox.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

type mockTxManager struct {
	mock.Mock
}

type mockEmailClient struct {
	mock.Mock
}

func (m *mockTxManager) ReadCommitted(ctx context.Context, fn database.Handler) error {
	_ = m.Called(ctx, fn)
	return fn(ctx)
}

func (m *mockEmailClient) GetUserEmail(ctx context.Context, userID, authToken string) (string, error) {
	args := m.Called(ctx, userID, authToken)
	return args.String(0), args.Error(1)
}

func (m *mockCustomerRepository) GetByID(ctx context.Context, customerID string) (entity.Customer, error) {
	args := m.Called(ctx, customerID)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (m *mockCustomerRepository) Create(ctx context.Context, in entity.CreateCustomer) (string, error) {
	args := m.Called(ctx, in)
	return args.String(0), args.Error(1)
}

func (m *mockCustomerRepository) Update(ctx context.Context, in entity.UpdateCustomer) (entity.Customer, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (m *mockCustomerRepository) GetByUserID(ctx context.Context, userID string) (entity.Customer, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (m *mockCustomerRepository) GetByUserName(ctx context.Context, userName string) (entity.Customer, error) {
	args := m.Called(ctx, userName)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (m *mockCustomerRepository) GetByIDs(ctx context.Context, ids []string) ([]entity.Customer, error) {
	args := m.Called(ctx, ids)

	var res []entity.Customer
	if args.Get(0) != nil {
		res = args.Get(0).([]entity.Customer)
	}

	return res, args.Error(1)
}

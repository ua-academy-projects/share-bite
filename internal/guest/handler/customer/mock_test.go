package customer

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type mockCustomerService struct {
	mock.Mock
}

func (m *mockCustomerService) Create(ctx context.Context, in entity.CreateCustomer) (string, error) {
	args := m.Called(ctx, in)
	return args.String(0), args.Error(1)
}

func (m *mockCustomerService) Update(ctx context.Context, in entity.UpdateCustomer) (entity.Customer, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (m *mockCustomerService) GetByUserName(ctx context.Context, userName string) (entity.Customer, error) {
	args := m.Called(ctx, userName)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (m *mockCustomerService) GetByUserID(ctx context.Context, userID string) (entity.Customer, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(entity.Customer), args.Error(1)
}

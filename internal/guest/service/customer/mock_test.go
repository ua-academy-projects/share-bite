package customer

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type mockCustomerRepository struct {
	mock.Mock
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

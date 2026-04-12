package customer

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type CustomerRepository interface {
	Create(ctx context.Context, in entity.CreateCustomer) (string, error)

	Update(ctx context.Context, in entity.UpdateCustomer) (entity.Customer, error)

	GetByUserID(ctx context.Context, userID string) (entity.Customer, error)
	GetByUserName(ctx context.Context, userName string) (entity.Customer, error)
	//GetByID(ctx context.Context, customerID string) (entity.Customer, error)
}

type service struct {
	customerRepo CustomerRepository
}

func New(
	customerRepo CustomerRepository,
) *service {
	return &service{
		customerRepo: customerRepo,
	}
}

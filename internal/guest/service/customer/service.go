package customer

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

type CustomerRepository interface {
	Create(ctx context.Context, in entity.CreateCustomer) (string, error)

	Update(ctx context.Context, in entity.UpdateCustomer) (entity.Customer, error)

	GetByUserID(ctx context.Context, userID string) (entity.Customer, error)
	GetByUserName(ctx context.Context, userName string) (entity.Customer, error)
	GetByID(ctx context.Context, customerID string) (entity.Customer, error)

	GetByIDs(ctx context.Context, ids []string) ([]entity.Customer, error)
}

type emailClient interface {
	GetUserEmail(ctx context.Context, userID, authToken string) (string, error)
}

type service struct {
	customerRepo CustomerRepository
	txManager    database.TxManager
	outboxWriter outbox.Writer
	adminClient  emailClient
}

func New(
	customerRepo CustomerRepository,
	txManager database.TxManager,
	outboxWriter outbox.Writer,
	adminClient emailClient,
) *service {
	return &service{
		customerRepo: customerRepo,
		txManager:    txManager,
		outboxWriter: outboxWriter,
		adminClient:  adminClient,
	}
}

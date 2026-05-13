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

type OutboxWriter interface {
	Enqueue(ctx context.Context, event outbox.Event) error
}

type service struct {
	customerRepo CustomerRepository
	outboxWriter OutboxWriter
	txManager    database.TxManager
}

func New(
	customerRepo CustomerRepository,
	outboxWriter OutboxWriter,
	txManager database.TxManager,
) *service {
	return &service{
		customerRepo: customerRepo,
		outboxWriter: outboxWriter,
		txManager:    txManager,
	}
}

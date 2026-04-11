package customer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

const (
	constraintUserIDKey   = "customers_user_id_key"
	constraintUserNameKey = "customers_username_key"
)

type repository struct {
	db database.Client
}

func New(db database.Client) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, in entity.CreateCustomer) (string, error) {
	sql := `
        INSERT INTO guest.customers(
            user_id,
            username,
            first_name,
            last_name,
            bio
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	q := database.Query{
		Name: "customer_repository.Create",
		Sql:  sql,
	}

	args := []any{
		in.UserID,
		in.UserName,
		in.FirstName,
		in.LastName,
		in.Bio,
	}

	var customerID string
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&customerID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// unique violation
			if pgErr.Code == pgerrcode.UniqueViolation {
				switch pgErr.ConstraintName {
				case constraintUserIDKey:
					return "", apperror.ErrCustomerAlreadyExists
				case constraintUserNameKey:
					return "", apperror.CustomerUserNameTaken(in.UserName)
				}
			}
		}

		return "", scanRowError(err)
	}

	return customerID, nil
}

func (r *repository) GetByUserID(ctx context.Context, userID string) (entity.Customer, error) {
	sql := `
        SELECT 
            id,
            user_id,
            username,
            first_name,
            last_name,
            avatar_object_key,
            bio,
            created_at
        FROM guest.customers
        WHERE user_id = $1
    `

	q := database.Query{
		Name: "customer_repository.GetByUserID",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, userID)
	if err != nil {
		return entity.Customer{}, executeSQLError(err)
	}
	defer row.Close()

	var c Customer
	if err := pgxscan.ScanOne(&c, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Customer{}, apperror.CustomerNotFoundUserID(userID)
		}

		return entity.Customer{}, scanRowError(err)
	}

	return c.ToEntity(), nil
}

func (r *repository) GetByUserName(ctx context.Context, userName string) (entity.Customer, error) {
	sql := `
        SELECT 
            id,
            user_id,
            username,
            first_name,
            last_name,
            avatar_object_key,
            bio,
            created_at
        FROM guest.customers
        WHERE username = $1
    `

	q := database.Query{
		Name: "customer_repository.GetByUserName",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, userName)
	if err != nil {
		return entity.Customer{}, executeSQLError(err)
	}
	defer row.Close()

	var c Customer
	if err := pgxscan.ScanOne(&c, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Customer{}, apperror.CustomerNotFoundUserName(userName)
		}

		return entity.Customer{}, scanRowError(err)
	}

	return c.ToEntity(), nil
}

func (r *repository) Update(ctx context.Context, in entity.UpdateCustomer) (entity.Customer, error) {
	sql := `
        UPDATE guest.customers
        SET
    `
	args := pgx.NamedArgs{}
	updates := make([]string, 0, 5)

	if in.UserName != nil {
		args["username"] = *in.UserName
		updates = append(updates, "username=@username")
	}
	if in.FirstName != nil {
		args["first_name"] = *in.FirstName
		updates = append(updates, "first_name=@first_name")
	}
	if in.LastName != nil {
		args["last_name"] = *in.LastName
		updates = append(updates, "last_name=@last_name")
	}
	if in.Bio != nil {
		args["bio"] = *in.Bio
		updates = append(updates, "bio=@bio")
	}
	if in.AvatarObjectKey != nil {
		args["avatar_object_key"] = *in.AvatarObjectKey
		updates = append(updates, "avatar_object_key=@avatar_object_key")
	}

	if len(updates) == 0 {
		return entity.Customer{}, apperror.ErrEmptyUpdate
	}

	sfx := "RETURNING id,user_id,username,first_name,last_name,avatar_object_key,bio,created_at"
	sql += fmt.Sprintf("%s WHERE user_id=@user_id %s", strings.Join(updates, ", "), sfx)
	args["user_id"] = in.UserID

	q := database.Query{
		Name: "customer_repository.Update",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == constraintUserNameKey {
			return entity.Customer{}, apperror.CustomerUserNameTaken(*in.UserName)
		}
		return entity.Customer{}, executeSQLError(err)
	}
	defer row.Close()

	var c Customer
	if err := pgxscan.ScanOne(&c, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Customer{}, apperror.CustomerNotFoundUserID(in.UserID)
		}

		return entity.Customer{}, scanRowError(err)
	}

	return c.ToEntity(), nil
}
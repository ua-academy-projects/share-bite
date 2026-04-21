package collection

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

const (
	constraintCollectionVenuesPkey        = "collection_venues_pkey"
	constraintCollectionCollaboratorsPkey = "collection_collaborators_pkey"
)

type repository struct {
	db database.Client
}

func New(db database.Client) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateCollection(ctx context.Context, in entity.CreateCollectionInput) (entity.Collection, error) {
	sql := `
		INSERT INTO guest.collections(customer_id, name, description, is_public)
		VALUES ($1, $2, $3, $4)
		RETURNING 
			id, 
			customer_id,
			name,
			description,
			is_public,
			created_at,
			updated_at
	`
	q := database.Query{
		Name: "collection_repository.CreateCollection",
		Sql:  sql,
	}

	args := []any{
		in.CustomerID,
		in.Name,
		in.Description,
		in.IsPublic,
	}

	row, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return entity.Collection{}, executeSQLError(err)
	}
	defer row.Close()

	var c Collection
	if err := pgxscan.ScanOne(&c, row); err != nil {
		return entity.Collection{}, scanRowError(err)
	}

	return c.ToEntity(), nil
}

func (r *repository) DeleteCollection(ctx context.Context, collectionID string) error {
	sql := `
		DELETE FROM guest.collections
		WHERE id = $1
	`
	q := database.Query{
		Name: "collection_repository.DeleteCollection",
		Sql:  sql,
	}

	cmd, err := r.db.DB().ExecContext(ctx, q, collectionID)
	if err != nil {
		return executeSQLError(err)
	}

	if cmd.RowsAffected() == 0 {
		return apperror.CollectionNotFoundID(collectionID)
	}

	return nil
}

func (r *repository) UpdateCollection(ctx context.Context, in entity.UpdateCollectionInput) (entity.Collection, error) {
	sql := `
		UPDATE guest.collections	
		SET
	`
	args := pgx.NamedArgs{}
	updates := make([]string, 0, 3)

	if in.Name != nil {
		args["name"] = *in.Name
		updates = append(updates, "name=@name")
	}
	if in.Description != nil {
		args["description"] = *in.Description
		updates = append(updates, "description=@description")
	}
	if in.IsPublic != nil {
		args["is_public"] = *in.IsPublic
		updates = append(updates, "is_public=@is_public")
	}

	if len(updates) == 0 {
		return entity.Collection{}, apperror.ErrEmptyUpdate
	}

	sfx := "RETURNING id, customer_id, name, description, is_public, created_at, updated_at"
	sql += fmt.Sprintf(
		`%s, updated_at = now() WHERE
			id=@id 
				AND
			EXISTS (
				SELECT 1 
				FROM guest.collection_collaborators collaborators 
				WHERE collections.id = collaborators.collection_id AND customer_id = @customer_id
			) 
		%s`,
		strings.Join(updates, ", "),
		sfx,
	)

	args["id"] = in.CollectionID
	args["customer_id"] = in.CustomerID

	q := database.Query{
		Name: "collection_repository.UpdateCollection",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, args)
	if err != nil {
		return entity.Collection{}, executeSQLError(err)
	}
	defer row.Close()

	var c Collection
	if err := pgxscan.ScanOne(&c, row); err != nil {
		return entity.Collection{}, scanRowError(err)
	}

	return c.ToEntity(), nil
}

func (r *repository) ListCustomerCollections(
	ctx context.Context,
	customerID string,
	cursorTime time.Time,
	cursorID string,
	limit int,
) ([]entity.Collection, error) {
	sql := `
        SELECT id, customer_id, name, description, is_public, created_at, updated_at
        FROM guest.collections collections
        WHERE 
			customer_id = @customer_id
				OR
			EXISTS (
				SELECT 1 
				FROM guest.collection_collaborators collaborators 
				WHERE collections.id = collaborators.collection_id AND customer_id = @customer_id
			)
    `
	args := pgx.NamedArgs{
		"customer_id": customerID,
		"limit":       limit,
	}

	if !cursorTime.IsZero() && cursorID != "" {
		sql += " AND (created_at, id) < (@cursor_time, @cursor_id)"
		args["cursor_time"] = cursorTime
		args["cursor_id"] = cursorID
	}
	sql += " ORDER BY created_at DESC, id DESC LIMIT @limit"

	q := database.Query{
		Name: "collection_repository.ListCustomerCollections",
		Sql:  sql,
	}

	var cs Collections
	rows, err := r.db.DB().QueryContext(ctx, q, args)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	if err := pgxscan.ScanAll(&cs, rows); err != nil {
		return nil, scanRowsError(err)
	}

	return cs.ToEntities(), nil
}

func (r *repository) GetCollection(ctx context.Context, collectionID string) (entity.Collection, error) {
	sql := `
		SELECT
			id, 
			customer_id,
			name,
			description,
			is_public,
			created_at,
			updated_at
		FROM guest.collections
		WHERE id = $1
	`
	q := database.Query{
		Name: "collection_repository.GetCollection",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, collectionID)
	if err != nil {
		return entity.Collection{}, executeSQLError(err)
	}
	defer row.Close()

	var c Collection
	if err := pgxscan.ScanOne(&c, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Collection{}, apperror.CollectionNotFoundID(collectionID)
		}

		return entity.Collection{}, scanRowError(err)
	}

	return c.ToEntity(), nil
}

func (r *repository) GetCollectionForUpdate(ctx context.Context, collectionID string) (entity.Collection, error) {
	sql := `
		SELECT
			id,
			customer_id,
			name,
			description,
			is_public,
			created_at,
			updated_at
		FROM guest.collections
		WHERE id = $1
		FOR UPDATE
	`
	q := database.Query{
		Name: "collection_repository.GetCollectionForUpdate",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, collectionID)
	if err != nil {
		return entity.Collection{}, executeSQLError(err)
	}
	defer row.Close()

	var c Collection
	if err := pgxscan.ScanOne(&c, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Collection{}, apperror.CollectionNotFoundID(collectionID)
		}

		return entity.Collection{}, scanRowError(err)
	}

	return c.ToEntity(), nil
}

func (r *repository) CountVenues(ctx context.Context, collectionID string) (int, error) {
	sql := `
		SELECT COUNT(*) AS count
		FROM guest.collection_venues
		WHERE collection_id = $1
	`
	q := database.Query{
		Name: "collection_repository.CountVenues",
		Sql:  sql,
	}

	var count int
	if err := r.db.DB().QueryRowContext(ctx, q, collectionID).Scan(&count); err != nil {
		return 0, scanRowError(err)
	}

	return count, nil
}

func (r *repository) GetMaxSortOrder(ctx context.Context, collectionID string) (float64, error) {
	sql := `
		SELECT COALESCE(MAX(sort_order), 0) 
		FROM guest.collection_venues 
		WHERE collection_id = $1
	`
	q := database.Query{
		Name: "collection_repository.GetMaxSortOrder",
		Sql:  sql,
	}

	var maxSortOrder float64
	err := r.db.DB().QueryRowContext(ctx, q, collectionID).Scan(&maxSortOrder)
	if err != nil {
		return 0, executeSQLError(err)
	}

	return maxSortOrder, nil
}

func (r *repository) CheckIfVenueInCollection(ctx context.Context, collectionID string, venueID int64) (bool, error) {
	sql := `
		SELECT EXISTS
		(SELECT 1 FROM guest.collection_venues WHERE collection_id = $1 AND venue_id = $2)
	`
	q := database.Query{
		Name: "collection_repository.CheckIfVenueInCollection",
		Sql:  sql,
	}

	var exists bool
	if err := r.db.DB().QueryRowContext(ctx, q, collectionID, venueID).Scan(&exists); err != nil {
		return false, scanRowError(err)
	}

	return exists, nil
}

func (r *repository) GetCollectionVenue(ctx context.Context, collectionID string, venueID int64) (entity.CollectionVenue, error) {
	sql := `
		SELECT
			collection_id,
			venue_id,
			sort_order,
			added_at
		FROM guest.collection_venues
		WHERE collection_id = $1 AND venue_id = $2
	`
	q := database.Query{
		Name: "collection_repository.GetCollectionVenue",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, collectionID, venueID)
	if err != nil {
		return entity.CollectionVenue{}, executeSQLError(err)
	}
	defer row.Close()

	var cv CollectionVenue
	if err := pgxscan.ScanOne(&cv, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.CollectionVenue{}, apperror.VenueNotFoundInCollection(venueID)
		}

		return entity.CollectionVenue{}, scanRowError(err)
	}

	return cv.ToEntity(), nil
}

func (r *repository) ListCollectionVenues(ctx context.Context, collectionID string) ([]entity.CollectionVenue, error) {
	sql := `
		SELECT
			collection_id,
			venue_id,
			sort_order,
			added_at
		FROM guest.collection_venues
		WHERE collection_id = $1
		ORDER BY sort_order ASC, venue_id ASC
	`
	q := database.Query{
		Name: "collection_repository.ListCollectionVenues",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, collectionID)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var cvs CollectionVenues
	if err := pgxscan.ScanAll(&cvs, rows); err != nil {
		return nil, scanRowsError(err)
	}

	return cvs.ToEntities(), nil
}

func (r *repository) AddVenue(ctx context.Context, collectionID string, venueID int64, sortOrder float64) error {
	sql := `
		INSERT INTO guest.collection_venues (collection_id, venue_id, sort_order)
		VALUES ($1, $2, $3)
	`
	q := database.Query{
		Name: "collection_repository.AddVenue",
		Sql:  sql,
	}

	_, err := r.db.DB().ExecContext(ctx, q, collectionID, venueID, sortOrder)
	if err != nil {
		// handle if venue is already in collection
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				switch pgErr.ConstraintName {
				case constraintCollectionVenuesPkey:
					return apperror.ErrVenueAlreadyInCollection
				}
			}
		}

		return executeSQLError(err)
	}

	return nil
}

func (r *repository) RemoveVenue(ctx context.Context, collectionID string, venueID int64) error {
	sql := `
		DELETE FROM guest.collection_venues 
		WHERE collection_id = $1 AND venue_id = $2
	`
	q := database.Query{
		Name: "collection_repository.RemoveVenue",
		Sql:  sql,
	}

	cmd, err := r.db.DB().ExecContext(ctx, q, collectionID, venueID)
	if err != nil {
		return executeSQLError(err)
	}

	if cmd.RowsAffected() == 0 {
		return apperror.VenueNotFoundInCollection(venueID)
	}

	return nil
}

func (r *repository) UpdateVenueSortOrder(ctx context.Context, collectionID string, venueID int64, sortOrder float64) error {
	sql := `
		UPDATE guest.collection_venues
		SET sort_order = $1
		WHERE collection_id = $2 AND venue_id = $3
	`
	q := database.Query{
		Name: "collection_repository.UpdateVenueSortOrder",
		Sql:  sql,
	}

	cmd, err := r.db.DB().ExecContext(ctx, q, sortOrder, collectionID, venueID)
	if err != nil {
		return executeSQLError(err)
	}

	if cmd.RowsAffected() == 0 {
		return apperror.VenueNotFoundInCollection(venueID)
	}

	return nil
}

func (r *repository) RebalanceCollectionSortOrders(ctx context.Context, collectionID string) error {
	sql := `
		UPDATE guest.collection_venues cv
		SET sort_order = new_orders.new_sort_order
		FROM (
			SELECT venue_id, ROW_NUMBER() OVER (ORDER BY sort_order ASC, venue_id ASC) * 100.0 AS new_sort_order
			FROM guest.collection_venues
			WHERE collection_id = $1
		) new_orders
		WHERE collection_id = $1 AND cv.venue_id = new_orders.venue_id
	`
	q := database.Query{
		Name: "collection_repository.RebalanceCollectionSortOrders",
		Sql:  sql,
	}

	if _, err := r.db.DB().ExecContext(ctx, q, collectionID); err != nil {
		return executeSQLError(err)
	}

	return nil
}

func (r *repository) HasVenuesBetween(ctx context.Context, collectionID string, venueID int64, lower float64, upper float64) (bool, error) {
	sql := `
		SELECT EXISTS (
			SELECT 1 FROM guest.collection_venues
			WHERE collection_id = $1 AND (sort_order > $3 AND sort_order < $4) AND venue_id != $2
		)
	`
	q := database.Query{
		Name: "collection_repository.HasVenuesBetween",
		Sql:  sql,
	}

	var has bool
	if err := r.db.DB().QueryRowContext(ctx, q, collectionID, venueID, lower, upper).Scan(&has); err != nil {
		return false, scanRowError(err)
	}

	return has, nil
}

func (r *repository) CreateCollaborator(ctx context.Context, collectionID string, customerID string) error {
	sql := `
		INSERT INTO guest.collection_collaborators(collection_id, customer_id)
		VALUES($1, $2)
	`

	q := database.Query{
		Name: "collection_repository.CreateCollaborator",
		Sql:  sql,
	}
	args := []any{
		collectionID,
		customerID,
	}

	_, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		// handle if collaborator is already added
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				switch pgErr.ConstraintName {
				case constraintCollectionCollaboratorsPkey:
					return apperror.CustomerAlreadyCollaborator(customerID)
				}
			}
		}

		return executeSQLError(err)
	}

	return nil
}

func (r *repository) DeleteCollaborator(ctx context.Context, collectionID string, customerID string) error {
	sql := `
		DELETE FROM guest.collection_collaborators
		WHERE collection_id = $1 AND customer_id = $2
	`

	q := database.Query{
		Name: "collection_repository.DeleteCollaborator",
		Sql:  sql,
	}
	args := []any{
		collectionID,
		customerID,
	}

	cmd, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return executeSQLError(err)
	}

	if cmd.RowsAffected() == 0 {
		return apperror.CollaboratorNotFound(customerID)
	}

	return nil
}

func (r *repository) CheckIfCollaborator(ctx context.Context, collectionID string, customerID string) (bool, error) {
	sql := `
		SELECT EXISTS
			(SELECT 1 FROM guest.collection_collaborators WHERE collection_id = $1 AND customer_id = $2)
	`

	q := database.Query{
		Name: "collection_repository.CheckIfCollaborator",
		Sql:  sql,
	}
	args := []any{
		collectionID,
		customerID,
	}

	var exists bool
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&exists); err != nil {
		return false, scanRowError(err)
	}

	return exists, nil
}

func (r *repository) CountCollaborators(ctx context.Context, collectionID string) (int, error) {
	sql := `
		SELECT COUNT(*) AS count
		FROM guest.collection_collaborators
		WHERE collection_id = $1
	`
	q := database.Query{
		Name: "collection_repository.CountCollaborators",
		Sql:  sql,
	}

	var count int
	if err := r.db.DB().QueryRowContext(ctx, q, collectionID).Scan(&count); err != nil {
		return 0, scanRowError(err)
	}

	return count, nil
}

func (r *repository) ListCollaborators(ctx context.Context, collectionID string) ([]entity.Collaborator, error) {
	sql := `
		SELECT 
			collaborators.collection_id AS collection_id,
			collaborators.customer_id AS customer_id,
			customers.username AS username,
			customers.avatar_object_key AS avatar_object_key,
			collaborators.added_at AS added_at
		FROM guest.collection_collaborators collaborators
		JOIN guest.customers customers ON collaborators.customer_id = customers.id
		WHERE collection_id = $1
	`
	q := database.Query{
		Name: "collection_repository.ListCollaborators",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, collectionID)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var collaborators Collaborators
	if err := pgxscan.ScanAll(&collaborators, rows); err != nil {
		return nil, scanRowsError(err)
	}

	return collaborators.ToEntities(), nil
}

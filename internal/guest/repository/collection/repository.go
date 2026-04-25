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

	constraintCollectionCollaboratorsCustomerFkey = "collection_collaborators_customer_id_fkey"

	constraintCollectionInvitationsUnique         = "idx_collection_invitations_unique"
	constraintCollectionInvitationsInviteeFkey    = "collection_invitations_invitee_id_fkey"
	constraintCollectionInvitationsCollectionFkey = "collection_invitations_collection_id_fkey"
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
		`%s, updated_at = now() WHERE id=@id %s`,
		strings.Join(updates, ", "),
		sfx,
	)

	args["id"] = in.CollectionID

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
				WHERE collections.id = collaborators.collection_id AND collaborators.customer_id = @customer_id
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

func (r *repository) CreateCollaborator(ctx context.Context, collectionID string, inviteeID string) error {
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
		inviteeID,
	}

	_, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// handle if collaborator is already added
			if pgErr.Code == pgerrcode.UniqueViolation {
				switch pgErr.ConstraintName {
				case constraintCollectionCollaboratorsPkey:
					return apperror.CustomerAlreadyCollaborator(inviteeID)
				}
			}

			// handle if invitee customer exists
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				switch pgErr.ConstraintName {
				case constraintCollectionCollaboratorsCustomerFkey:
					return apperror.InviteeCustomerNotFoundID(inviteeID)
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

func (r *repository) CreateInvitation(ctx context.Context, in entity.InviteCollaboratorInput) (string, error) {
	sql := `
		INSERT INTO guest.collection_invitations(collection_id, status, inviter_id, invitee_id, expires_at)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id
	`
	args := []any{
		in.CollectionID,
		entity.PendingInvitationStatus,
		in.InviterID,
		in.InviteeID,
		in.Expiry,
	}

	q := database.Query{
		Name: "collection_repository.CreateInvitation",
		Sql:  sql,
	}

	var id string
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if pgErr.ConstraintName == constraintCollectionInvitationsUnique {
					return "", apperror.InvitationAlreadySent(in.CollectionID, in.InviteeID)
				}
			case pgerrcode.ForeignKeyViolation:
				// handle if invitee customer exists
				if pgErr.ConstraintName == constraintCollectionInvitationsInviteeFkey {
					return "", apperror.InviteeCustomerNotFoundID(in.InviteeID)
				}

				// handle if collection still exists
				if pgErr.ConstraintName == constraintCollectionInvitationsCollectionFkey {
					return "", apperror.CollectionNotFoundID(in.CollectionID)
				}
			}
		}

		return "", scanRowError(err)
	}

	return id, nil
}

func (r *repository) GetInvitation(ctx context.Context, invitationID string) (entity.Invitation, error) {
	sql := `
		SELECT 
			id,
			collection_id,
			status,
			inviter_id,
			invitee_id,
			expires_at,
			last_sent_at,
			created_at
		FROM guest.collection_invitations
		WHERE id = $1
	`
	q := database.Query{
		Name: "collection_repository.GetInvitation",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, invitationID)
	if err != nil {
		return entity.Invitation{}, executeSQLError(err)
	}
	defer row.Close()

	var invitation Invitation
	if err := pgxscan.ScanOne(&invitation, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Invitation{}, apperror.InvitationNotFoundID(invitationID)
		}

		return entity.Invitation{}, scanRowError(err)
	}

	return invitation.ToEntity(), nil
}

func (r *repository) GetInvitationForUpdate(ctx context.Context, invitationID string) (entity.Invitation, error) {
	sql := `
		SELECT 
			id,
			collection_id,
			status,
			inviter_id,
			invitee_id,
			expires_at,
			last_sent_at,
			created_at
		FROM guest.collection_invitations
		WHERE id = $1
		FOR UPDATE
	`
	q := database.Query{
		Name: "collection_repository.GetInvitationForUpdate",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, invitationID)
	if err != nil {
		return entity.Invitation{}, executeSQLError(err)
	}
	defer row.Close()

	var invitation Invitation
	if err := pgxscan.ScanOne(&invitation, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Invitation{}, apperror.InvitationNotFoundID(invitationID)
		}

		return entity.Invitation{}, scanRowError(err)
	}

	return invitation.ToEntity(), nil
}

func (r *repository) GetInvitationByInvitee(ctx context.Context, collectionID string, inviteeID string) (entity.Invitation, error) {
	sql := `
		SELECT 
			id,
			collection_id,
			status,
			inviter_id,
			invitee_id,
			expires_at,
			last_sent_at,
			created_at
		FROM guest.collection_invitations
		WHERE collection_id = $1 AND invitee_id = $2
	`
	q := database.Query{
		Name: "collection_repository.GetInvitationByInvitee",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, collectionID, inviteeID)
	if err != nil {
		return entity.Invitation{}, executeSQLError(err)
	}
	defer row.Close()

	var invitation Invitation
	if err := pgxscan.ScanOne(&invitation, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Invitation{}, apperror.InvitationNotFoundForInvitee(collectionID, inviteeID)
		}

		return entity.Invitation{}, scanRowError(err)
	}

	return invitation.ToEntity(), nil
}

func (r *repository) UpdateInvitationStatus(ctx context.Context, invitationID string, status entity.InvitationStatus) error {
	sql := `
		UPDATE guest.collection_invitations
		SET status = $1
		WHERE id = $2
	`

	q := database.Query{
		Name: "collection_repository.UpdateInvitationStatus",
		Sql:  sql,
	}

	cmd, err := r.db.DB().ExecContext(ctx, q, status, invitationID)
	if err != nil {
		return executeSQLError(err)
	}

	if cmd.RowsAffected() == 0 {
		return apperror.InvitationNotFoundID(invitationID)
	}

	return nil
}

func (r *repository) RefreshInvitation(ctx context.Context, invitationID string, newExpiry time.Time) error {
	sql := `
		UPDATE guest.collection_invitations
		SET 
			expires_at = $1,
			status = 'pending',
			last_sent_at = now()
		WHERE id = $2
	`
	q := database.Query{
		Name: "collection_repository.RefreshInvitation",
		Sql:  sql,
	}

	cmd, err := r.db.DB().ExecContext(ctx, q, newExpiry, invitationID)
	if err != nil {
		return executeSQLError(err)
	}
	if cmd.RowsAffected() == 0 {
		return apperror.InvitationNotFoundID(invitationID)
	}

	return nil
}

func (r *repository) ListInvitations(ctx context.Context, in entity.ListInvitationsInput) ([]entity.EnrichedInvitation, error) {
	sql := `
		SELECT
			i.id AS id,
			i.status AS status,
			i.created_at AS created_at,
			i.expires_at AS expires_at,
			
			c.id AS collection_id,
			c.name AS collection_name,

			inviter.id AS inviter_id,
			inviter.username AS inviter_username,
			inviter.avatar_object_key AS inviter_avatar_object_key,

			invitee.id AS invitee_id,
			invitee.username AS invitee_username,
			invitee.avatar_object_key AS invitee_avatar_object_key
		FROM guest.collection_invitations i
		JOIN guest.collections c ON i.collection_id = c.id
		JOIN guest.customers inviter ON i.inviter_id = inviter.id
		JOIN guest.customers invitee ON i.invitee_id = invitee.id
	`

	args := pgx.NamedArgs{}
	conditions := make([]string, 0, 5)

	if in.CollectionID != nil {
		args["collection_id"] = *in.CollectionID
		conditions = append(conditions, "c.id=@collection_id")
	}
	if in.InviteeID != nil {
		args["invitee_id"] = *in.InviteeID
		conditions = append(conditions, "invitee.id=@invitee_id")
	}
	if in.InviterID != nil {
		args["inviter_id"] = *in.InviterID
		conditions = append(conditions, "inviter.id=@inviter_id")
	}
	if in.Status != nil {
		args["status"] = *in.Status
		conditions = append(conditions, "i.status=@status")
	}
	if len(in.CursorID) > 0 {
		args["cursor_id"] = in.CursorID
		conditions = append(conditions, "i.id<@cursor_id")
	}

	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}

	sql += " ORDER BY i.id DESC LIMIT @limit"
	args["limit"] = in.Limit

	q := database.Query{
		Name: "collection_repository.ListInvitations",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	eis := make(EnrichedInvitations, 0, in.Limit)
	if err := pgxscan.ScanAll(&eis, rows); err != nil {
		return nil, scanRowsError(err)
	}

	return eis.ToEntities(), nil
}

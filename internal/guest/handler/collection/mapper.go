package collection

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/validator"
)

func listMyCollectionsRequestToInput(req listMyCollectionsRequest, customerID string) (entity.ListCustomerCollectionsInput, error) {
	var cursorTime time.Time
	var cursorID string
	if req.PageToken != "" {
		ct, cid, err := parsePageToken(req.PageToken)
		if err != nil {
			return entity.ListCustomerCollectionsInput{}, apperror.ErrInvalidPageToken
		}
		cursorTime = ct
		cursorID = cid
	}

	limit := req.PageSize
	switch {
	case limit <= 0:
		limit = defaultCollectionsLimit
	case limit > maxCollectionsLimit:
		limit = maxCollectionsLimit
	}

	return entity.ListCustomerCollectionsInput{
		CustomerID: customerID,
		CursorTime: cursorTime,
		CursorID:   cursorID,
		Limit:      limit + 1,
	}, nil
}

func generatePageToken(createdAt time.Time, id string) string {
	timeStr := createdAt.Format(time.RFC3339Nano)
	raw := fmt.Sprintf("%s|%s", timeStr, id)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func parsePageToken(token string) (time.Time, string, error) {
	if token == "" {
		return time.Time{}, "", nil
	}

	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return time.Time{}, "", fmt.Errorf("invalid token encoding: %w", err)
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return time.Time{}, "", fmt.Errorf("invalid token format")
	}

	if err := uuid.Validate(parts[1]); err != nil {
		return time.Time{}, "", fmt.Errorf("invalid cursor id format")
	}

	createdAt, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, "", fmt.Errorf("invalid token time format: %w", err)
	}

	return createdAt, parts[1], nil
}

func listCustomerCollectionsOutputToResponse(out entity.ListCustomerCollectionsOutput) listMyCollectionsResponse {
	var pageToken string
	if out.NextCursorTime != nil && out.NextCursorID != nil {
		pageToken = generatePageToken(*out.NextCursorTime, *out.NextCursorID)
	}

	collections := make([]collectionResponse, 0, len(out.Collections))
	for _, c := range out.Collections {
		collections = append(collections, collectionToResponse(c))
	}

	return listMyCollectionsResponse{
		Collections:   collections,
		NextPageToken: pageToken,
	}
}

func enrichedVenueItemsToResponse(venues []entity.EnrichedVenueItem) listVenuesResponse {
	list := make([]enrichedVenueItemResponse, 0, len(venues))
	for _, v := range venues {
		list = append(list, enrichedVenueItemToResponse(v))
	}

	return listVenuesResponse{
		Venues: list,
	}
}

func reorderVenueRequestToReorderVenue(params reorderVenueParams, req reorderVenueRequest, customerID string) entity.ReorderVenueInput {
	return entity.ReorderVenueInput{
		CollectionID: params.CollectionID,
		VenueID:      params.VenueID,

		CustomerID: customerID,

		PrevVenueID: req.PrevVenueID,
		NextVenueID: req.NextVenueID,
	}
}

func createCollectionRequestToCreateCollection(req createCollectionRequest, customerID string) (entity.CreateCollectionInput, error) {
	var valErrors []validator.ValidationErrorItem

	name := strings.TrimSpace(req.Name)
	if len(name) == 0 {
		valErrors = append(valErrors, validator.ValidationErrorItem{
			Field:   "name",
			Message: "This field must be at least 1 characters long",
		})

		return entity.CreateCollectionInput{}, &validator.ValidationError{Errors: valErrors}
	}

	if req.Description != nil {
		if v := strings.TrimSpace(*req.Description); v == "" {
			req.Description = nil
		} else {
			req.Description = &v
		}
	}

	return entity.CreateCollectionInput{
		CustomerID:  customerID,
		Name:        name,
		Description: req.Description,
		IsPublic:    *req.IsPublic,
	}, nil
}

func updateCollectionRequestToUpdateCollection(body updateCollectionBody, collectionID, customerID string) (entity.UpdateCollectionInput, error) {
	var valErrors []validator.ValidationErrorItem

	if body.Name != nil {
		v := strings.TrimSpace(*body.Name)
		if v == "" {
			valErrors = append(valErrors, validator.ValidationErrorItem{
				Field:   "name",
				Message: "This field must be at least 1 character long",
			})
		} else {
			body.Name = &v
		}
	}
	if body.Description != nil {
		v := strings.TrimSpace(*body.Description)
		body.Description = &v
	}

	if len(valErrors) > 0 {
		return entity.UpdateCollectionInput{}, &validator.ValidationError{Errors: valErrors}
	}

	if body.Name == nil && body.Description == nil && body.IsPublic == nil {
		return entity.UpdateCollectionInput{}, apperror.ErrEmptyUpdate
	}

	return entity.UpdateCollectionInput{
		CollectionID: collectionID,
		CustomerID:   customerID,

		Name:        body.Name,
		Description: body.Description,
		IsPublic:    body.IsPublic,
	}, nil
}

func inviteCollaboratorRequestToInput(collectionID string, inviteeID string, inviterID string) entity.InviteCollaboratorInput {
	return entity.InviteCollaboratorInput{
		CollectionID: collectionID,
		InviterID:    inviterID,
		InviteeID:    inviteeID,
	}
}

func removeCollaboratorRequestToRemoveCollaborator(params removeCollaboratorParams, customerID string) entity.RemoveCollaboratorInput {
	return entity.RemoveCollaboratorInput{
		CollectionID: params.CollectionID,
		CustomerID:   customerID,

		TargetCustomerID: params.TargetCustomerID,
	}
}

func listInvitationsRequestToInput(params listInvitationsParams, callerID string) (entity.ListInvitationsInput, error) {
	var cursorID string
	if len(params.PageToken) > 0 {
		decoded, err := base64.RawURLEncoding.DecodeString(params.PageToken)
		if err != nil {
			return entity.ListInvitationsInput{}, apperror.ErrInvalidPageToken
		}

		cursorID = string(decoded)
	}

	limit := params.PageSize
	switch {
	case limit <= 0:
		limit = defaultInvitationsLimit
	case limit > maxInvitationsLimit:
		limit = maxInvitationsLimit
	}

	var status *entity.InvitationStatus
	if params.Status != nil {
		s := entity.InvitationStatus(*params.Status)
		status = &s
	}

	if params.InviterID != nil && *params.InviterID != callerID {
		return entity.ListInvitationsInput{}, apperror.ErrCannotListOthersOutboundInvites
	}
	if params.InviteeID != nil && *params.InviteeID != callerID {
		return entity.ListInvitationsInput{}, apperror.ErrCannotListOthersInboundInvites
	}

	return entity.ListInvitationsInput{
		CollectionID: params.CollectionID,
		InviterID:    params.InviterID,
		InviteeID:    params.InviteeID,
		CallerID:     callerID,
		Status:       status,
		CursorID:     cursorID,
		Limit:        limit + 1,
	}, nil
}

func (h *handler) listInvitationsOutputToResponse(out entity.ListInvitationsOutput) listInvitationsResponse {
	var pageToken string
	if len(out.NextCursorID) > 0 {
		pageToken = base64.RawURLEncoding.EncodeToString([]byte(out.NextCursorID))
	}

	list := make([]invitationResponse, 0, len(out.Invitations))
	for _, i := range out.Invitations {
		list = append(list, h.invitationToResponse(i))
	}

	return listInvitationsResponse{
		Invitations:   list,
		NextPageToken: pageToken,
	}
}

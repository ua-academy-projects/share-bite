package collection

import (
	"strings"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/validator"
)

func listMyCollectionsRequestToInput(req listMyCollectionsRequest, customerID string) entity.ListCustomerCollectionsInput {
	return entity.ListCustomerCollectionsInput{
		CustomerID: customerID,
		PageSize:   req.PageSize,
		PageToken:  req.PageToken,
	}
}

func listCustomerCollectionsOutputToResponse(out entity.ListCustomerCollectionsOutput) listMyCollectionsResponse {
	collections := make([]collectionResponse, 0, len(out.Collections))
	for _, c := range out.Collections {
		collections = append(collections, collectionToResponse(c))
	}

	return listMyCollectionsResponse{
		Collections:   collections,
		NextPageToken: out.NextPageToken,
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

func reorderVenueRequestToReorderVenue(uri reorderVenueUri, body reorderVenueBody, customerID string) entity.ReorderVenueInput {
	return entity.ReorderVenueInput{
		CollectionID: uri.CollectionID,
		VenueID:      uri.VenueID,

		CustomerID: customerID,

		PrevVenueID: body.PrevVenueID,
		NextVenueID: body.NextVenueID,
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

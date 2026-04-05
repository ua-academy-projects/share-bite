package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/request"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary		Update a collection
// @Description	Updates the details of an existing collection (name, description, or visibility).
// @Description	Only the provided fields will be updated. Fails if the user does not own the collection.
//
// @Tags			collections
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			collectionId	path		string						true	"Collection ID (UUID)"
// @Param			request			body		updateCollectionBody		true	"Collection fields to update"
//
// @Success		200				{object}	updateCollectionResponse	"Collection successfully updated"
// @Failure		400				{object}	response.ErrorResponse		"Validation error (e.g., empty update payload) or bad request"
// @Failure		401				{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		403				{object}	response.ErrorResponse		"Forbidden: Customer profile not found or user does not own this collection"
// @Failure		404				{object}	response.ErrorResponse		"Not Found: Collection not found"
// @Failure		500				{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/collections/{collectionId} [patch]
func (h *handler) updateCollection(c *gin.Context) {
	var uri updateCollectionUri
	if err := request.BindUri(c, &uri); err != nil {
		c.Error(err)
		return
	}

	var body updateCollectionBody
	if err := request.BindJSON(c, &body); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	in, err := updateCollectionRequestToUpdateCollection(body, uri.CollectionID, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	collection, err := h.service.UpdateCollection(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := updateCollectionResponse{Collection: collectionToResponse(collection)}
	c.JSON(http.StatusOK, resp)
}

type updateCollectionUri struct {
	CollectionID string `uri:"collectionId" binding:"required,uuid"`
}

type updateCollectionBody struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description" binding:"omitempty,max=300"`
	IsPublic    *bool   `json:"isPublic" binding:"omitempty"`
}

type updateCollectionResponse struct {
	Collection collectionResponse `json:"collection"`
}

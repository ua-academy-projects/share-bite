package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary		Delete a collection
// @Description	Deletes an existing collection.
// @Description	Fails if the user does not own the collection.
//
// @Tags			collections
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			collectionId	path	string	true	"Collection ID (UUID)"
//
// @Success		204				"Collection successfully deleted"
// @Failure		400				{object}	response.ErrorResponse		"Invalid collection ID format"
// @Failure		401				{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		403				{object}	response.ErrorResponse		"Forbidden: Customer profile not found or user does not own this collection"
// @Failure		404				{object}	response.ErrorResponse		"Not Found: Collection not found"
// @Failure		500				{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/collections/{collectionId} [delete]
func (h *handler) deleteCollection(c *gin.Context) {
	var req deleteCollectionRequest
	if err := request.BindUri(c, &req); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()

	if err := h.service.DeleteCollection(ctx, req.CollectionID, customerID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

type deleteCollectionRequest struct {
	CollectionID string `uri:"collectionId" binding:"required,uuid"`
}

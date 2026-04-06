package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/request"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary		List venues in a collection
// @Description	Retrieves a list of venues belonging to a specific collection, ordered by their sort order.
// @Description	The collection must be public or belong to the authenticated user.
//
// @Tags			collections
// @Accept			json
// @Produce		json
//
// @Param			collectionId	path		string						true	"Collection ID (UUID)"
//
// @Success		200				{object}	listVenuesResponse			"Successfully retrieved the list of venues"
// @Failure		400				{object}	response.ErrorResponse		"Invalid collection ID format"
// @Failure		401				{object}	response.AuthErrorResponse	"Unauthorized: Token was provided but is invalid or expired"
// @Failure		403				{object}	response.ErrorResponse		"Forbidden: Token is valid but customer profile not found"
// @Failure		404				{object}	response.ErrorResponse		"Not Found: Collection does not exist or user cannot access it"
// @Failure		500				{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/collections/{collectionId}/venues [get]
func (h *handler) listVenues(c *gin.Context) {
	var req listVenuesRequest
	if err := request.BindUri(c, &req); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetOptionalCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	venues, err := h.service.ListVenues(ctx, req.CollectionID, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := enrichedVenueItemsToResponse(venues)
	c.JSON(http.StatusOK, resp)
}

type listVenuesRequest struct {
	CollectionID string `uri:"collectionId" binding:"required,uuid"`
}

type listVenuesResponse struct {
	Venues []enrichedVenueItemResponse `json:"venues"`
}

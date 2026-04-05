package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/request"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
)

// @Summary Add a venue to a collection
// @Description Adds a specific venue to an existing collection.
// @Description Fails if the collection is full (limit is 100 venues),
// @Description if the venue is already in the collection, or if the user does not own it.
//
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param collectionId path string true "Collection ID (UUID)"
// @Param venueId path string true "Venue ID (UUID)"
//
// @Success 204 "Venue successfully added to the collection"
// @Failure 400 {object} response.ErrorResponse "Invalid path parameters, or collection is full (limit is 100 venues)"
// @Failure 401 {object} response.AuthErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 403 {object} response.ErrorResponse "Forbidden: Customer profile not found or user does not own this collection"
// @Failure 404 {object} response.ErrorResponse "Not Found: Collection not found"
// @Failure 409 {object} response.ErrorResponse "Conflict: Venue is already in the collection"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
//
// @Router /collections/{collectionId}/venues/{venueId} [post]
func (h *handler) addVenue(c *gin.Context) {
	var req addVenueRequest
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

	if err := h.service.AddVenue(ctx, req.CollectionID, customerID, req.VenueID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

type addVenueRequest struct {
	CollectionID string `uri:"collectionId" binding:"required,uuid"`
	VenueID      string `uri:"venueId" binding:"required,uuid"`
}

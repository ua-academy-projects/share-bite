package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/request"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary Remove a venue from a collection
// @Description Removes a specific venue from an existing collection.
// @Description Fails if the user does not own the collection or if the venue is not in the collection.
//
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param collectionId path string true "Collection ID (UUID)"
// @Param venueId path string true "Venue ID (UUID)"
//
// @Success 204 "Venue successfully removed from the collection"
// @Failure 400 {object} response.ErrorResponse "Invalid path parameters"
// @Failure 401 {object} response.AuthErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 403 {object} response.ErrorResponse "Forbidden: Customer profile not found or user does not own this collection"
// @Failure 404 {object} response.ErrorResponse "Not Found: Collection not found or venue is not in the collection"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
//
// @Router /collections/{collectionId}/venues/{venueId} [delete]
func (h *handler) removeVenue(c *gin.Context) {
	var req removeVenueRequest
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

	if err := h.service.RemoveVenue(ctx, req.CollectionID, customerID, req.VenueID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

type removeVenueRequest struct {
	CollectionID string `uri:"collectionId" binding:"required,uuid"`
	VenueID      string `uri:"venueId" binding:"required,uuid"`
}

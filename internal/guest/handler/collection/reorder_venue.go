package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/request"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary Reorder a venue in a collection
// @Description Updates the position (sort order) of a specific venue within a collection.
// @Description You must provide either the ID of the venue that should precede it,
// @Description the ID of the venue that should follow it, or both.
//
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param collectionId path string true "Collection ID (UUID)"
// @Param venueId path string true "Venue ID (UUID)"
// @Param request body reorderVenueBody true "Reorder details (prevVenueId and/or nextVenueId)"
//
// @Success 204 "Venue successfully reordered"
// @Failure 400 {object} response.ErrorResponse "Validation error (e.g., missing both neighbor IDs) or bad request"
// @Failure 401 {object} response.AuthErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 403 {object} response.ErrorResponse "Forbidden: Customer profile not found or user does not own this collection"
// @Failure 404 {object} response.ErrorResponse "Not Found: Collection not found, target venue not in collection, or neighbor not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
//
// @Router /collections/{collectionId}/venues/{venueId}/reorder [post]
func (h *handler) reorderVenue(c *gin.Context) {
	var uri reorderVenueUri
	if err := request.BindUri(c, &uri); err != nil {
		c.Error(err)
		return
	}

	var body reorderVenueBody
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
	in := reorderVenueRequestToReorderVenue(uri, body, customerID)

	if err := h.service.ReorderVenue(ctx, in); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

type reorderVenueUri struct {
	CollectionID string `uri:"collectionId" binding:"required,uuid"`
	VenueID      string `uri:"venueId" binding:"required,uuid"`
}

type reorderVenueBody struct {
	PrevVenueID *string `json:"prevVenueId" binding:"required_without=NextVenueID,omitempty,uuid"`
	NextVenueID *string `json:"nextVenueId" binding:"required_without=PrevVenueID,omitempty,uuid"`
}

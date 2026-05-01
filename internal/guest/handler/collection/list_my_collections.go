package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

const (
	maxCollectionsLimit     = 100
	defaultCollectionsLimit = 20
)

// @Summary		List current user's collections
// @Description	Retrieves a paginated list of collections belonging to the authenticated user,
// @Description	ordered by creation date descending.
//
// @Tags			collections
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			pageSize	query		int							false	"Number of items to return (default is 20, max is 100)"
// @Param			pageToken	query		string						false	"Pagination token returned from a previous request"
//
// @Success		200			{object}	listMyCollectionsResponse	"Successfully retrieved the list of collections"
// @Failure		400			{object}	response.ErrorResponse		"Invalid pagination parameters"
// @Failure		401			{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		403			{object}	response.ErrorResponse		"Forbidden: Customer profile not found"
// @Failure		500			{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/collections/me [get]
func (h *handler) listMyCollections(c *gin.Context) {
	var req listMyCollectionsRequest
	if err := request.BindQuery(c, &req); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	in, err := listMyCollectionsRequestToInput(req, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	out, err := h.service.ListCustomerCollections(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := listCustomerCollectionsOutputToResponse(out)
	c.JSON(http.StatusOK, resp)
}

type listMyCollectionsRequest struct {
	PageSize  int    `form:"pageSize,default=20" binding:"gte=1,lte=100"`
	PageToken string `form:"pageToken"`
}

type listMyCollectionsResponse struct {
	Collections   []collectionResponse `json:"collections"`
	NextPageToken string               `json:"nextPageToken,omitempty"`
}

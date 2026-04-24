package follow

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"net/http"
)

// @Summary		Get my following
// @Description	Returns a paginated list of users the current customer is following.
//
// @Tags			follow
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			pageSize	query	int		false	"Page size (default 20)"
// @Param			pageToken	query	string	false	"Cursor for pagination"
//
// @Success		200	{object}	dto.ListCustomersResponse
// @Failure		400	{object}	response.ErrorResponse		"Invalid request"
// @Failure		401	{object}	response.AuthErrorResponse	"Unauthorized"
// @Failure		500	{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/customers/me/following [get]
func (h *handler) listMyFollowing(c *gin.Context) {
	var req dto.ListFollowingRequest
	if err := request.BindQuery(c, &req); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	customerID, err := httpctx.GetCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	in := listFollowingRequestToInput(req, customerID, &customerID)
	out, err := h.service.ListFollowing(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, h.listCustomersResponse(out.Customers, out.NextPageToken))

}

// @Summary		Get following
// @Description	Returns a paginated list of users the customer is following.
// @Description	If the profile is private, only the owner can view following.
//
// @Tags			follow
// @Accept			json
// @Produce		json
//
// @Param			id			path	string	true	"Customer ID (UUID)"
// @Param			pageSize	query	int		false	"Page size (default 20)"
// @Param			pageToken	query	string	false	"Cursor for pagination"
//
// @Success		200	{object}	dto.ListCustomersResponse
// @Failure		400	{object}	response.ErrorResponse		"Invalid request or page token"
// @Failure		403	{object}	response.ErrorResponse		"Following list is private"
// @Failure		404	{object}	response.ErrorResponse		"Customer not found"
// @Failure		500	{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/customers/{id}/following [get]
func (h *handler) listFollowing(c *gin.Context) {
	var req dto.ListFollowingRequest
	if err := request.BindQuery(c, &req); err != nil {
		c.Error(err)
		return
	}

	targetCustomerID := c.Param("id")
	if targetCustomerID == "" {
		c.Error(apperror.ErrInvalidParam)
		return
	}

	var requesterCustomerID *string
	if id, err := httpctx.GetCustomerID(c); err == nil {
		requesterCustomerID = &id
	}

	in := listFollowingRequestToInput(req, targetCustomerID, requesterCustomerID)

	out, err := h.service.ListFollowing(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, h.listCustomersResponse(out.Customers, out.NextPageToken))
}

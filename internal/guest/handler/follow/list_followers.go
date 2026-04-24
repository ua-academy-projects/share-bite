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

// @Summary		Get followers
// @Description	Returns a paginated list of followers for a specific customer.
// @Description	If the profile is private, only the owner can view followers.
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
// @Failure		403	{object}	response.ErrorResponse		"Followers list is private"
// @Failure		404	{object}	response.ErrorResponse		"Customer not found"
// @Failure		500	{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/customers/{id}/followers [get]
func (h *handler) listFollowers(c *gin.Context) {
	var req dto.ListFollowersRequest
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

	in := listFollowersRequestToInput(req, targetCustomerID, requesterCustomerID)

	out, err := h.service.ListFollowers(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := h.listCustomersResponse(out.Customers, out.NextPageToken)
	c.JSON(http.StatusOK, resp)
}

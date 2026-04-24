package follow

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"net/http"
)

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
	c.JSON(http.StatusOK, listFollowingOutputToResponse(out))

}

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

	c.JSON(http.StatusOK, listFollowingOutputToResponse(out))
}

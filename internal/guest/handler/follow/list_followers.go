package follow

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"net/http"
)

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

	resp := listFollowersOutputToResponse(out)
	c.JSON(http.StatusOK, resp)
}

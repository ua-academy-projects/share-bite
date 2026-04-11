package post

import (
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

func (h *handler) get(c *gin.Context) {
	var req getRequest
	if err := request.BindUri(c, &req); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	customerID := getOptionalCustomerID(c, h.customerService)
	post, err := h.service.Get(ctx, req.ID, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := getResponse{Post: postToResponse(post, h.storage)}
	c.JSON(http.StatusOK, resp)
}

type getRequest struct {
	ID string `uri:"id" binding:"required,numeric"`
}

type getResponse struct {
	Post postResponse `json:"post"`
}

func getOptionalCustomerID(c *gin.Context, custSvc customerService) string {
	userID, err := httpctx.GetUserID(c)
	if err == nil && userID != "" {
		customer, err := custSvc.GetByUserID(c.Request.Context(), userID)
		if err == nil {
			return customer.ID
		}
	}
	return ""
}

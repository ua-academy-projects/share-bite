package follow

import (
	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"net/http"
)

// @Summary		Unfollow a user
// @Description	Removes a follow relationship between the current user and another customer.
//
// @Tags			follow
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			id	path	string	true	"Target customer ID (UUID)"
//
// @Success		204	"Successfully unfollowed"
// @Failure		400	{object}	response.ErrorResponse		"Invalid request"
// @Failure		401	{object}	response.AuthErrorResponse	"Unauthorized"
// @Failure		404	{object}	response.ErrorResponse		"Follow relationship not found"
// @Failure		500	{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/customers/{id}/follow [delete]
func (h *handler) unfollow(c *gin.Context) {
	customerID, err := httpctx.GetCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	targetCustomerID := c.Param("id")
	if targetCustomerID == "" {
		c.Error(apperror.ErrInvalidParam)
		return
	}
	if err := h.service.Unfollow(c.Request.Context(), customerID, targetCustomerID); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

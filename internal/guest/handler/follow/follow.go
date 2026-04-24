package follow

import (
	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"net/http"
)

// @Summary		Follow a user
// @Description	Creates a follow relationship between the current user and another customer.
// @Description	Fails if already following, target does not exist, or trying to follow yourself.
//
// @Tags			follow
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			id	path	string	true	"Target customer ID (UUID)"
//
// @Success		201	{object}	entity.CustomerFollow	"Successfully followed"
// @Failure		400	{object}	response.ErrorResponse		"Invalid request"
// @Failure		401	{object}	response.AuthErrorResponse	"Unauthorized"
// @Failure		404	{object}	response.ErrorResponse		"Customer not found"
// @Failure		409	{object}	response.ErrorResponse		"Already following"
// @Failure		500	{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/customers/{id}/follow [post]
func (h *handler) follow(c *gin.Context) {
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
	follow, err := h.service.Follow(c.Request.Context(), customerID, targetCustomerID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, follow)
}

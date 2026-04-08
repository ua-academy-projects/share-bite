package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary		Create a new collection
// @Description	Creates a new collection for the authenticated user.
//
// @Tags			collections
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			request	body		createCollectionRequest		true	"Collection details"
//
// @Success		201		{object}	createCollectionResponse	"Collection successfully created"
// @Failure		400		{object}	response.ErrorResponse		"Validation error (e.g., name is empty or too long)"
// @Failure		401		{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		403		{object}	response.ErrorResponse		"Forbidden: Customer profile not found or insufficient permissions"
// @Failure		500		{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/collections [post]
func (h *handler) createCollection(c *gin.Context) {
	var req createCollectionRequest
	if err := request.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	in, err := createCollectionRequestToCreateCollection(req, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	collection, err := h.service.CreateCollection(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := createCollectionResponse{Collection: collectionToResponse(collection)}
	c.JSON(http.StatusCreated, resp)
}

type createCollectionRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Description *string `json:"description" binding:"omitempty,max=300"`
	IsPublic    *bool   `json:"isPublic" binding:"required"`
}

type createCollectionResponse struct {
	Collection collectionResponse `json:"collection"`
}

package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/request"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary Get a collection by ID
// @Description Retrieves details of a specific collection.
// @Description Returns the collection if it is public or if it belongs to the authenticated user.
// @Description You can pass a token optionally to access your private collections.
//
// @Tags collections
// @Accept json
// @Produce json
//
// @Param collectionId path string true "Collection ID (UUID)"
//
// @Success 200 {object} getCollectionResponse "Successfully retrieved the collection"
// @Failure 400 {object} response.ErrorResponse "Invalid collection ID format"
// @Failure 401 {object} response.AuthErrorResponse "Unauthorized: Token was provided but is invalid or expired"
// @Failure 403 {object} response.AuthErrorResponse "Forbidden: Token is valid but customer profile not found"
// @Failure 404 {object} response.ErrorResponse "Not Found: Collection does not exist, is private, or does not belong to the user"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
//
// @Router /collections/{collectionId} [get]
func (h *handler) getCollection(c *gin.Context) {
	var req getCollectionRequest
	if err := request.BindUri(c, &req); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetOptionalCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	collection, err := h.service.GetCollection(ctx, req.CollectionID, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := getCollectionResponse{Collection: collectionToResponse(collection)}
	c.JSON(http.StatusOK, resp)
}

type getCollectionRequest struct {
	CollectionID string `uri:"collectionId" binding:"required,uuid"`
}

type getCollectionResponse struct {
	Collection collectionResponse `json:"collection"`
}

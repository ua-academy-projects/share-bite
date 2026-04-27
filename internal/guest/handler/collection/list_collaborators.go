package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary		List collaborators in a collection
// @Description	Retrieves a list of collaborators for a specific collection.
// @Description	The collection must be public or the authenticated user must be the owner or a collaborator.
// @Description	Token is optional for public collections.
//
// @Tags			collections
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			collectionId	path		string						true	"Collection ID (UUID)"
//
// @Success		200				{object}	listCollaboratorsResponse	"Successfully retrieved the list of collaborators"
// @Failure		400				{object}	response.ErrorResponse		"Invalid collection ID format"
// @Failure		401				{object}	response.AuthErrorResponse	"Unauthorized: Token was provided but is invalid or expired"
// @Failure		403				{object}	response.ErrorResponse		"Forbidden: Token is valid but customer profile not found, or collection is private"
// @Failure		404				{object}	response.ErrorResponse		"Not Found: Collection does not exist"
// @Failure		500				{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/collections/{collectionId}/collaborators [get]
func (h *handler) listCollaborators(c *gin.Context) {
	var params listCollaboratorsParams
	if err := request.BindUri(c, &params); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetOptionalCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	collaborators, err := h.service.ListCollaborators(ctx, params.CollectionID, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := h.collaboratorsToListResponse(collaborators)
	c.JSON(http.StatusOK, resp)
}

type listCollaboratorsParams struct {
	CollectionID string `uri:"collectionId" binding:"required,uuid"`
}

type listCollaboratorsResponse struct {
	Collaborators []collaboratorResponse `json:"collaborators"`
}

func (h *handler) collaboratorsToListResponse(collaborators []entity.Collaborator) listCollaboratorsResponse {
	list := make([]collaboratorResponse, 0, len(collaborators))
	for _, c := range collaborators {
		list = append(list, h.collaboratorToResponse(c))
	}

	return listCollaboratorsResponse{Collaborators: list}
}

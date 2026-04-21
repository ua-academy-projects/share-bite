package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

func (h *handler) listCollaborators(c *gin.Context) {
	var uri listCollaboratorsUri
	if err := request.BindUri(c, &uri); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetOptionalCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	collaborators, err := h.service.ListCollaborators(ctx, uri.CollectionID, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := h.collaboratorsToListResponse(collaborators)
	c.JSON(http.StatusOK, resp)
}

type listCollaboratorsUri struct {
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

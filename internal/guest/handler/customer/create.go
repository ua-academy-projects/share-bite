package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/request"
)

func (h *handler) create(c *gin.Context) {
	var req createRequest
	if err := request.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	in := createRequestToCreateCustomer(req, userID)

	customerID, err := h.service.Create(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := createResponse{CustomerID: customerID}
	c.JSON(http.StatusCreated, resp)
}

type createRequest struct {
	UserName  string `json:"userName" binding:"required,min=3,max=30"`
	FirstName string `json:"firstName" binding:"required,min=2,max=50"`
	LastName  string `json:"lastName" binding:"required,min=2,max=50"`

	Bio *string `json:"bio" binding:"omitempty,max=500"`
}

type createResponse struct {
	CustomerID string `json:"customerId"`
}

func createRequestToCreateCustomer(req createRequest, userID string) entity.CreateCustomer {
	return entity.CreateCustomer{
		UserID:    userID,
		UserName:  req.UserName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
	}
}

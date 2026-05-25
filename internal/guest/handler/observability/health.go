package observability

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, healthCheckResponse{Status: "OK"})
}

type healthCheckResponse struct {
	Status string `json:"status"`
}

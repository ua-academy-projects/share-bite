package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// Check godoc
// @Summary      Infrastructure health check
// @Description  Verifies that the auth service is up and running for Kubernetes probes.
// @Tags         Infrastructure
// @Produce      json
// @Success      200      {object}  map[string]string "Returns status healthy"
// @Router       /health [get]
func (h *Handler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

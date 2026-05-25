package observability

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/config"
)

const (
	// TODO: replace it with dynamic value
	appVersion = "v1.0.0"
)

func (h *handler) info(c *gin.Context) {
	resp := infoResponse{
		Version: appVersion,
		Env:     config.Config().App.Stage(),
	}
	c.JSON(http.StatusOK, resp)
}

type infoResponse struct {
	Version string `json:"version"`
	Env     string `json:"environment"`
}

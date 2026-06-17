package observability

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/pkg/version"
)

func (h *handler) info(c *gin.Context) {
	resp := infoResponse{
		Version:    version.Version,
		CommitHash: version.CommitHash,
		BuildTime:  version.BuildTime,
		Env:        config.Config().App.Stage(),
	}

	c.JSON(http.StatusOK, resp)
}

type infoResponse struct {
	Version    string `json:"version"`
	CommitHash string `json:"commit"`
	BuildTime  string `json:"buildTime"`
	Env        string `json:"environment"`
}

package observability

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/config"
)

const (
	connectedStatus    = "connected"
	disconnectedStatus = "disconnected"
)

func (h *handler) statusCheck(c *gin.Context) {
	var (
		ctx = c.Request.Context()

		dbStatus    = connectedStatus
		redisStatus = connectedStatus
	)

	if err := h.db.DB().Ping(ctx); err != nil {
		dbStatus = disconnectedStatus
	}
	if err := h.redis.Ping(ctx).Err(); err != nil {
		redisStatus = disconnectedStatus
	}

	resp := statusCheckResponse{
		App:   config.Config().App.Name(),
		DB:    dbStatus,
		Redis: redisStatus,
	}
	c.JSON(http.StatusOK, resp)
}

type statusCheckResponse struct {
	App   string `json:"app"`
	DB    string `json:"database"`
	Redis string `json:"redis"`
}

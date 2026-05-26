package observability

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

const (
	connectedStatus    = "connected"
	disconnectedStatus = "disconnected"
)

func (h *handler) statusCheck(c *gin.Context) {
	var (
		ctx        = c.Request.Context()
		statusCode = http.StatusOK

		dbStatus    = connectedStatus
		redisStatus = connectedStatus
	)

	if err := h.db.DB().Ping(ctx); err != nil {
		logger.ErrorKV(ctx, "database ping failed", "error", err)

		dbStatus = disconnectedStatus
		statusCode = http.StatusServiceUnavailable
	}
	if err := h.redis.Ping(ctx).Err(); err != nil {
		logger.ErrorKV(ctx, "redis ping failed", "error", err)

		redisStatus = disconnectedStatus
		statusCode = http.StatusServiceUnavailable
	}

	resp := statusCheckResponse{
		App:   config.Config().App.Name(),
		DB:    dbStatus,
		Redis: redisStatus,
	}
	c.JSON(statusCode, resp)
}

type statusCheckResponse struct {
	App   string `json:"app"`
	DB    string `json:"database"`
	Redis string `json:"redis"`
}

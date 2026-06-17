package observability

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

const (
	adminUserRole     = "admin"
	moderatorUserRole = "moderator"
)

type handler struct {
	db    database.Client
	redis *redis.Client
}

func RegisterHandlers(
	r *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	db database.Client,
	redis *redis.Client,
) {
	h := &handler{
		db:    db,
		redis: redis,
	}

	r.GET("/health", h.healthCheck)

	protected := r.Group("/").Use(authMiddleware, middleware.RequireRoles(adminUserRole, moderatorUserRole))

	protected.GET("/status", h.statusCheck)
	protected.GET("/info", h.info)
}

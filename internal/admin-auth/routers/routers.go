package routers

import (
	"github.com/gin-gonic/gin"
	authhttp "github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"
)

func SetupRouter(
	r *gin.RouterGroup,
	authHandler *authhttp.Handler,
	limiter gin.HandlerFunc,
) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/recover-access", limiter, authHandler.RecoverAccess)
		authGroup.POST("/reset-password", limiter, authHandler.ResetPassword)
	}
}

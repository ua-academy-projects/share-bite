package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	authhttp "github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/ghAuth"
)

func SetupRouter(r *gin.RouterGroup, authHandler *authhttp.Handler, authMiddleware gin.HandlerFunc, limiter gin.HandlerFunc, gh *ghAuth.Handler) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	{
		authGroup := r.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/refresh", authHandler.Refresh)
			authGroup.POST("/oauth/:provider/callback", authHandler.OAuthCallback)

			protectedUserGroup := r.Group("/user").Use(authMiddleware)
			{
				protectedUserGroup.POST("/link/:provider", authHandler.OAuthLinkAccount)
			}
			authGroup.POST("/recover-access", limiter, authHandler.RecoverAccess)
			authGroup.POST("/reset-password", limiter, authHandler.ResetPassword)
			
			authGroup.GET("/github", gin.WrapF(gh.Login))
			authGroup.GET("/github/callback", gin.WrapF(gh.Callback))
			authGroup.GET("/github/success", gin.WrapF(gh.Success))
		}
	}
}

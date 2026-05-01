package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/ua-academy-projects/share-bite/docs/api/admin-auth"
	adminhttp "github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/admin"
	"github.com/ua-academy-projects/share-bite/internal/middleware"

	authhttp "github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/provider/github"
)

func SetupRouter(r *gin.RouterGroup, authHandler *authhttp.Handler, adminHandler *adminhttp.Handler,  gh github.Handler,authMiddleware gin.HandlerFunc, limiter gin.HandlerFunc) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/oauth/:provider/callback", authHandler.OAuthCallback)
		authGroup.POST("/recover-access", limiter, authHandler.RecoverAccess)
		authGroup.POST("/reset-password", limiter, authHandler.ResetPassword)

	}

	protectedUserGroup := r.Group("/user").Use(authMiddleware)
	{
		protectedUserGroup.POST("/logout", authHandler.Logout)
		protectedUserGroup.POST("/link/:provider", authHandler.OAuthLinkAccount)
		protectedUserGroup.POST("/sessions/revoke-all", authHandler.RevokeAllSessions)
	}
	adminGroup := r.Group("/admin").Use(authMiddleware)
	{
		adminGroup.GET("/users", middleware.RequireRoles("admin", "moderator"), adminHandler.GetUsersList)
		adminGroup.GET("/users/:id", middleware.RequireRoles("admin", "moderator"), adminHandler.GetUserDetails)
		adminGroup.PATCH("/users/:id/role", middleware.RequireRoles("admin"), adminHandler.ChangeUserRole)
		protectedUserGroup := r.Group("/user").Use(authMiddleware)
		{
			protectedUserGroup.POST("/logout", authHandler.Logout)
			protectedUserGroup.POST("/link/:provider", authHandler.OAuthLinkAccount)
			protectedUserGroup.POST("/sessions/revoke-all", authHandler.RevokeAllSessions)

			authGroup.GET("/github", gin.WrapF(gh.Login))
			authGroup.GET("/github/callback", gin.WrapF(gh.Callback))
			authGroup.GET("/github/success", gin.WrapF(gh.Success))
		}
	}
}

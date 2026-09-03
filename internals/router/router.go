package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hareshkhan01/PollyRoute/internals/handlers"
	"github.com/hareshkhan01/PollyRoute/internals/middleware"
)

func SetupRouter(
	routeHandler *handlers.RouteHandler,
	authHandlers *handlers.AuthHandlers,
	jwtSecret string,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")

	publicAuth := api.Group("/auth")
	{
		publicAuth.POST("/login", authHandlers.Login)
		publicAuth.POST("/register", authHandlers.Register)
		publicAuth.POST("/refresh", authHandlers.Refresh)
	}

	protectedAuth := api.Group("/auth")
	protectedAuth.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protectedAuth.POST("/logout", authHandlers.Logout)
	}

	protectedApi := api.Group("/routes")

	protectedApi.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protectedApi.POST("/analyze", routeHandler.AnalyzedRoutes)
	}

	return router
}

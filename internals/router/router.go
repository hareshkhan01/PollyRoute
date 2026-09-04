package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hareshkhan01/PollyRoute/internals/handlers"
	"github.com/hareshkhan01/PollyRoute/internals/middleware"
)

func SetupRouter(
	routeHandler *handlers.RouteHandler,
	authHandlers *handlers.AuthHandlers,
	jwtSecret string,
	ctx context.Context,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")

	publicAuth := api.Group("/auth")
	{
		publicAuth.POST("/login", authHandlers.Login(ctx))
		publicAuth.POST("/register", authHandlers.Register(ctx))
		publicAuth.POST("/refresh", authHandlers.Refresh(ctx))
	}

	protectedAuth := api.Group("/auth")
	protectedAuth.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protectedAuth.POST("/logout", authHandlers.Logout(ctx))
	}

	protectedApi := api.Group("/routes")

	protectedApi.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protectedApi.POST("/analyze", routeHandler.AnalyzedRoutes)
	}

	return router
}

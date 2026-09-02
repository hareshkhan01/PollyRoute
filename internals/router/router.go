package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hareshkhan01/PollyRoute/internals/handlers"
)

func SetupRouter(
	routeHandler *handlers.RouteHandler,
	authHandlers *handlers.AuthHandlers,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	auth := api.Group("/auth")

	api.POST("/routes/analyze", routeHandler.AnalyzedRoutes)

	auth.POST("/login", authHandlers.Login)
	auth.POST("/register", authHandlers.Register)
	auth.POST("/refresh", authHandlers.Refresh)
	auth.POST("/logout", authHandlers.Logout)

	return router
}

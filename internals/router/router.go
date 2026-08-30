package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hareshkhan01/PollyRoute/internals/handlers"
)

func SetupRouter(routeHandler *handlers.RouteHandler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	api.POST("/routes/analyze", routeHandler.AnalyzedRoutes)

	return router
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hareshkhan01/PollyRoute/internals/application"
	"github.com/hareshkhan01/PollyRoute/internals/domain"
	"github.com/hareshkhan01/PollyRoute/internals/request"
	"github.com/hareshkhan01/PollyRoute/internals/response"
)

type RouteHandler struct {
	routeAnalysis *application.RouteAnalyzeService
}

func NewRouteHandler(routeAnalysis *application.RouteAnalyzeService) *RouteHandler {
	return &RouteHandler{
		routeAnalysis: routeAnalysis,
	}
}
func (h *RouteHandler) AnalyzedRoutes(c *gin.Context) {
	var req request.AnalyzeRouteRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	origin := domain.Coordinate{
		Lat: req.Origin.Lat,
		Lng: req.Origin.Lng,
	}

	destination := domain.Coordinate{
		Lat: req.Destination.Lat,
		Lng: req.Destination.Lng,
	}

	analyzedRoutes, err := h.routeAnalysis.AnalyzeRoutes(
		c.Request.Context(),
		&origin,
		&destination,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	apiResponse := response.ToAnalysisRouteResponse(analyzedRoutes)

	c.JSON(http.StatusOK, apiResponse)

}

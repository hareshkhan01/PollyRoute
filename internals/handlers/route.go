package handlers

import (
	"github.com/hareshkhan01/PollyRoute/internals/application"
)
type RouteHandler struct{
	routeAnalysis *application.RouteAnalyzeService
}
func NewRouteHandler(routeAnalysis *application.RouteAnalyzeService) RouteHandler{
	return &RouteHandler{
		routeAnalysis: routeAnalysis,
	}
}
func (h *RouteHandler) 
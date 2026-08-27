package segmentation

import (
	// "github.com/hareshkhan01/PollyRoute/pkg/geo"
	"github.com/hareshkhan01/PollyRoute/internals/domain"
)

// This  function is used to divide the route into 5km segment using haversine mathematical formula
func RouteSegmentationService(coordinates []domain.Coordinate) []domain.RouteSegment {
	segments := make([]domain.RouteSegment, 0)

}

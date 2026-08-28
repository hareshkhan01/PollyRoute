package segmentation

import (
	"github.com/hareshkhan01/PollyRoute/internals/domain"
	"github.com/hareshkhan01/PollyRoute/pkg/geo"
)

// This  function is used to divide the route into 5km segment using haversine mathematical formula
func RouteSegmentationService(coordinates []domain.Coordinate) []domain.RouteSegment {
	segments := make([]domain.RouteSegment, 0)

	var currDistance float64
	var currCoordinates []domain.Coordinate
	index := 0
	var midPoint domain.Coordinate
	isMidPointSet := false
	i := 1
	for ; i < len(coordinates); i++ {
		distance := geo.Haversine(
			coordinates[i-1],
			coordinates[i],
		)
		currCoordinates = append(currCoordinates, coordinates[i-1])
		currDistance += distance
		if currDistance >= 2500.0 && !isMidPointSet {
			midPoint = coordinates[i]
			isMidPointSet = true
		}
		if currDistance >= 5000.0 {
			currCoordinates = append(currCoordinates, coordinates[i])
			segments = append(segments, domain.RouteSegment{
				Index:       index,
				Coordinates: currCoordinates,
				Distance:    currDistance,
				Midpoint:    midPoint,
			})
			currCoordinates = nil
			currDistance = 0
			midPoint = domain.Coordinate{}
			isMidPointSet = false
			index++
		}
	}
	if len(currCoordinates) > 0 {
		if !isMidPointSet {
			midPoint = currCoordinates[len(currCoordinates)/2]
		}
		currCoordinates = append(currCoordinates, coordinates[i-1]) // add the last remaining coordinate
		segments = append(segments, domain.RouteSegment{
			Index:       index,
			Coordinates: currCoordinates,
			Distance:    currDistance,
			Midpoint:    midPoint,
		})
	}
	return segments
}

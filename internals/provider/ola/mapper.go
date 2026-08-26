package ola

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hareshkhan01/PollyRoute/internals/domain"
	"github.com/hareshkhan01/PollyRoute/pkg/polyline"
)

func MapResponse(response Response) (domain.Routes, error) {
	if response.Status != "SUCCESS" {
		return domain.Routes{}, fmt.Errorf("OLA Direction Request/Response failed status %v", response.Status)
	}

	routes := make([]domain.Route, 0, len(response.Routes))

	for _, route := range response.Routes {
		mappedRoute, err := mapRoute(route)
		if err != nil {
			return domain.Routes{}, err
		}
		routes = append(routes, mappedRoute)
	}
	return domain.Routes{
		Routes: routes,
	}, nil
}
func mapRoute(route Route) (domain.Route, error) {

	if len(route.Legs) <= 0 {
		return domain.Route{}, fmt.Errorf("OLA route contains no legs")
	}

	leg := route.Legs[0]

	coordiantes, err := polyline.Decode(route.Polyline)
	if err != nil {
		return domain.Route{}, fmt.Errorf(
			"decode OLA overview polyline: %w",
			err,
		)
	}
	trafficRanges, err := mapTrafficAdvisory(route.TraverAdvisory)
	if err != nil {
		return domain.Route{}, fmt.Errorf("Failed to map traffic advisory %w", err)
	}
	isSuspicious := hasSuspiciousManeuver(leg.Steps)

	return domain.Route{
		TotalDistance: leg.Distance,
		TotalDuration: leg.Duration,
		Coordinates:   coordiantes,
		TrafficRanges: trafficRanges,
		IsSuspicious:  isSuspicious,
		Copyrights:    route.Copyrights,
	}, nil

}
func mapTrafficAdvisory(trafficAdvisory string) ([]domain.TrafficRange, error) {
	if len(trafficAdvisory) <= 0 {
		return nil, nil
	}
	entries := strings.Split(trafficAdvisory, "|")
	ranges := make([]domain.TrafficRange, 0, len(entries))

	for _, entry := range entries {
		parts := strings.Split(entry, ",")

		if len(parts) != 3 {
			return nil, fmt.Errorf("Inavid travel_advisory/traffic entry.")
		}

		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf(
				"invalid traffic start index %q: %w",
				parts[0],
				err,
			)
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf(
				"invalid traffic end index %q: %w",
				parts[1],
				err,
			)
		}
		congestion, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf(
				"invalid traffic congestion index %q: %w",
				parts[2],
				err,
			)
		}
		ranges = append(ranges, domain.TrafficRange{
			Start:      start,
			End:        end,
			Congestion: congestion,
		})
	}
	return ranges, nil
}

func hasSuspiciousManeuver(steps []Step) bool {
	for _, step := range steps {
		if step.Maneuver == "u-turn" {
			return true
		}
	}
	return false
}

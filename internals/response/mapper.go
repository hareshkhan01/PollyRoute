package response

import (
	"github.com/hareshkhan01/PollyRoute/internals/domain"
)

func ToAnalysisRouteReponse(routes []domain.AnalyzedRoute) AnalysisRouteReponse {
	response := AnalysisRouteReponse{
		Routes: make([]RouteResponse, 0, len(routes)),
	}

	for i, route := range routes {
		response.Routes = append(response.Routes, toRouteResponse(route, i))
	}
	return response
}

func toRouteResponse(route domain.AnalyzedRoute, id int) RouteResponse {
	return RouteResponse{
		ID:       id + 1,
		Polyline: route.Route.Polyline,
		Distance: route.Route.TotalDistance,
		Duration: route.Route.TotalDuration,
		Scores: ScoreResponse{
			PollutionScore: route.AqiScore,
			WeatherScore:   route.WeatherScore,
			TrafficScore:   route.TrafficScore,
			BalancedScore:  route.BalancedScore,
		},
		Traffic:      toTrafficResponse(route.Route.TrafficRanges),
		Segments:     toSegmentsResponse(route.Segments),
		IsSuspicious: route.Route.IsSuspicious,
		Copyright:    route.Route.Copyrights,
	}
}

func toTrafficResponse(trafficRanges []domain.TrafficRange) []TrafficResponse {
	trafficResponse := make([]TrafficResponse, 0)

	for _, trafficRange := range trafficRanges {
		if trafficRange.Congestion >= 3 {
			trafficResponse = append(trafficResponse, TrafficResponse{
				Start:      trafficRange.Start,
				End:        trafficRange.End,
				Congestion: trafficRange.Congestion,
			})
		}
	}
	return trafficResponse
}

func toSegmentsResponse(segments []domain.RouteSegment) []SegmentResponse {
	segmentsResponse := make([]SegmentResponse, 0, len(segments))

	for _, segment := range segments {
		response := SegmentResponse{
			Index:    segment.Index,
			Distance: segment.Distance,
			Midpoint: ResponseCoordinate{
				Lat: segment.Midpoint.Lat,
				Lng: segment.Midpoint.Lng,
			},
		}

		if segment.AqiAndPollutants != nil {
			response.AqiAndPollutants = AQIAndPollutantsResponse{
				Aqi:               float64(segment.AqiAndPollutants.Aqi),
				DominantPollutant: segment.AqiAndPollutants.DominantPollutant,
			}
		}

		if segment.Weather != nil {
			response.Weather = WeatherResponse{
				Temperature:             float64(segment.Weather.Temperature),
				ThunderstormProbability: float64(segment.Weather.ThunderstromProbability),
			}
		}

		segmentsResponse = append(segmentsResponse, response)
	}

	return segmentsResponse
}

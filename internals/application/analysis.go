package application

import (
	"context"

	"github.com/hareshkhan01/PollyRoute/internals/domain"
	"github.com/hareshkhan01/PollyRoute/internals/service/aqi"
	"github.com/hareshkhan01/PollyRoute/internals/service/route"
	"github.com/hareshkhan01/PollyRoute/internals/service/scoring"
	"github.com/hareshkhan01/PollyRoute/internals/service/segmentation"
	"github.com/hareshkhan01/PollyRoute/internals/service/weather"
)

type RouteAnalyzeService struct {
	routeDirectionService *route.DirectionService
	weatherService        *weather.WeatherService
	aqiService            *aqi.AqiService
	segmentationService   *segmentation.SegmentationService
	scoreService          *scoring.ScoringService
}

func NewRouteAnalyzeService(
	routeDirectionService *route.DirectionService,
	weatherService *weather.WeatherService,
	aqiService *aqi.AqiService,
	segmentationService *segmentation.SegmentationService,
	scoreService *scoring.ScoringService,
) *RouteAnalyzeService {
	return &RouteAnalyzeService{
		routeDirectionService: routeDirectionService,
		weatherService:        weatherService,
		aqiService:            aqiService,
		segmentationService:   segmentationService,
		scoreService:          scoreService,
	}
}

func (r *RouteAnalyzeService) AnalyzeRoutes(
	ctx context.Context,
	origin *domain.Coordinate,
	destination *domain.Coordinate,
) ([]domain.AnalyzedRoute, error) {
	// Segment the route into smaller segments based on traffic ranges
	routes, err := r.routeDirectionService.GetDirection(ctx, origin, destination)
	if err != nil {
		return nil, err
	}

	analyzedRoutes := make([]domain.AnalyzedRoute, 0, len(routes.Routes))

	for _, route := range routes.Routes {
		segments := r.segmentationService.SegmentatRoute(route.Coordinates)
		for i := range segments {
			segment := &segments[i]
			aqiResponse, err := r.aqiService.GetAqi(
				ctx,
				segment.Midpoint,
			)
			if err != nil {
				return nil, err
			}
			weatherResponse, err := r.weatherService.GetWeather(
				ctx,
				segment.Midpoint,
			)
			if err != nil {
				return nil, err
			}
			segment.AqiAndPollutants = aqiResponse
			segment.Weather = weatherResponse
		}
		pollutionScore := r.scoreService.CalculatePollutionScore(segments)
		weatherScore := r.scoreService.CalculateWeatherScore(segments)
		trafficScore := r.scoreService.CalculateTrafficScore(route.TrafficRanges)
		balancedScore := r.scoreService.CalculateBalancedScore(weatherScore, pollutionScore, trafficScore)
		analyzedRoutes = append(analyzedRoutes, domain.AnalyzedRoute{
			Route:         route,
			Segments:      segments,
			AqiScore:      pollutionScore,
			WeatherScore:  weatherScore,
			TrafficScore:  trafficScore,
			BalancedScore: balancedScore,
		})
	}
	return analyzedRoutes, nil
}

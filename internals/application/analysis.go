package application

import (
	"context"
	"sync"

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
			weatherCh := make(chan *domain.Weather, 1)
			aqiCh := make(chan *domain.AQIAndPollutants, 1)

			var wg sync.WaitGroup
			wg.Add(2)

			var weatherError error
			var aqiError error

			go func() {
				defer wg.Done()
				weatherCh <- func() *domain.Weather {
					weather, err := r.weatherService.GetWeather(
						ctx, segment.Midpoint,
					)
					weatherError = err
					return weather
				}()
			}()

			go func() {
				defer wg.Done()
				aqiCh <- func() *domain.AQIAndPollutants {
					aqi, err := r.aqiService.GetAqi(
						ctx,
						segment.Midpoint,
					)
					aqiError = err
					return aqi
				}()
			}()
			wg.Wait()

			if weatherError != nil {
				return nil, weatherError
			}

			if aqiError != nil {
				return nil, aqiError
			}

			segment.Weather = <-weatherCh
			segment.AqiAndPollutants = <-aqiCh

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

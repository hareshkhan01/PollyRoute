package application

import (
	"context"
	"fmt"
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
		var wg sync.WaitGroup
		for i := range segments {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				segment := &segments[i]

				var (
					weatherResponse *domain.Weather
					aqiResponse     *domain.AQIAndPollutants

					weatherErr error
					aqiErr     error
				)

				var inner sync.WaitGroup
				inner.Add(2)

				go func() {
					defer inner.Done()

					weatherResponse, weatherErr = r.weatherService.GetWeather(
						ctx, segment.Midpoint,
					)
				}()

				go func() {
					defer inner.Done()
					aqiResponse, aqiErr = r.aqiService.GetAqi(
						ctx,
						segment.Midpoint,
					)
				}()
				inner.Wait()

				if weatherErr != nil {
					fmt.Println(weatherErr)
					return
				}

				if aqiErr != nil {
					fmt.Println(aqiErr)
					return
				}

				segment.Weather = weatherResponse
				segment.AqiAndPollutants = aqiResponse
			}(i)

		}
		wg.Wait()
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

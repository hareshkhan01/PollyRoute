package weather

import (
	"context"
	"fmt"

	"github.com/hareshkhan01/PollyRoute/internals/domain"
	"github.com/hareshkhan01/PollyRoute/internals/provider/google/weather"
)

type WeatherService struct {
	weatherClient *weather.Client
}

func NewWeatherService(client *weather.Client) *WeatherService {
	return &WeatherService{
		weatherClient: client,
	}
}

func (w *WeatherService) GetWeather(
	ctx context.Context,
	coordianate domain.Coordinate,
) (*domain.Weather, error) {

	response, err := w.weatherClient.Weather(ctx, coordianate)
	if err != nil {
		return nil, fmt.Errorf("fetch weather response: %w", err)
	}
	mappedResponse, err := weather.MapResponse(response)
	if err != nil {
		return nil, fmt.Errorf("map weather response: %w", err)
	}
	return mappedResponse, nil
}

package aqi

import (
	"context"
	"fmt"

	"github.com/hareshkhan01/PollyRoute/internals/domain"
	"github.com/hareshkhan01/PollyRoute/internals/provider/google/aqi"
)

type AqiService struct {
	aqiClient *aqi.Client
}

func NewAqiService(client *aqi.Client) *AqiService {
	return &AqiService{
		aqiClient: client,
	}
}
func (a *AqiService) GetAqi(
	ctx context.Context,
	coordinate domain.Coordinate,
) (*domain.AQIAndPollutants, error) {
	response, err := a.aqiClient.Aqi(
		ctx,
		coordinate,
	)

	if err != nil {
		return nil, fmt.Errorf("Aqi Response: %w", err)
	}
	mappedResponse, err := aqi.MapResponse(response)
	if err != nil {
		return nil, fmt.Errorf("Map Aqi Response: %w", err)
	}
	return mappedResponse, nil
}

package route

import (
	"context"
	"fmt"

	"github.com/hareshkhan01/PollyRoute/internals/domain"
	"github.com/hareshkhan01/PollyRoute/internals/provider/ola"
)

type DirectionService struct {
	olaClient *ola.Client
}

func NewDirectionService(client *ola.Client) *DirectionService {
	return &DirectionService{
		olaClient: client,
	}
}
func (o *DirectionService) GetDirection(ctx context.Context, origin *domain.Coordinate, destination *domain.Coordinate) (*domain.Routes, error) {

	originString := fmt.Sprintf("%f,%f", origin.Lat, origin.Lng)
	destinationString := fmt.Sprintf("%f,%f", destination.Lat, destination.Lng)
	response, err := o.olaClient.Directions(ctx, originString, destinationString, true)

	if err != nil {
		return nil, fmt.Errorf("fetch directions from OLA: %w", err)
	}

	mappedResponse, err := ola.MapResponse(response)
	if err != nil {
		return nil, fmt.Errorf("map OLA response: %w", err)
	}
	return mappedResponse, nil
}

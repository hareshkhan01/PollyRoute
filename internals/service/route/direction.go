package route

import (
	"context"
	"fmt"
	"log"

	"github.com/hareshkhan01/PollyRoute/internals/config"
	"github.com/hareshkhan01/PollyRoute/internals/domain"
	"github.com/hareshkhan01/PollyRoute/internals/provider/ola"
)

func DirectionService(origin *domain.Coordinate, destination *domain.Coordinate) *domain.Routes {
	config, err := config.Load()

	if err != nil {
		log.Fatalf("Can not  Load the Config File %w", err)
	}

	olaClient, err := ola.NewClient(
		config.OLA_API_KEY,
		config.OLA_DIRECTIONS_URL,
	)
	if err != nil {
		log.Fatalf("Can not create Ola Client %w", err)
	}
	ctx := context.Background()
	originString := fmt.Sprintf("%s,%s", origin.Lat, origin.Lng)
	destinationString := fmt.Sprintf("%s,%s", destination.Lat, destination.Lng)
	response, err := olaClient.Directions(&ctx, originString, destinationString, true)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

}

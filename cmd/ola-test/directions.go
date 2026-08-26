package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hareshkhan01/PollyRoute/internals/config"
	"github.com/hareshkhan01/PollyRoute/internals/provider/ola"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := ola.NewClient(config.OLA_API_KEY, config.OLA_DIRECTIONS_URL)
	if err != nil {
		log.Fatalf("Failed to create OLA Client: %v", err)
	}
	ctx := context.Background()
	result, err := client.Directions(&ctx, "22.186032,87.980629", "22.595622,88.263088", true)
	if err != nil {
		log.Fatalf("Failed to get directions from OLA: %v", err)
	}
	fmt.Println("Status: ", result.Status)
	fmt.Println("Routes: ", len(result.Routes))

	for i, route := range result.Routes {
		if len(route.Legs) == 0 {
			continue
		}
		leg := route.Legs[0]
		fmt.Printf("\nRoute %d\n", i+1)
		fmt.Printf("Distance: %d meters, Duration: %d seconds\n", leg.Distance, leg.Duration)
		for j, step := range leg.Steps {
			fmt.Printf("Step %d: Maneuver: %s\n", j+1, step.Maneuver)
		}
		// fmt.Println("Overview Polyline: ", route.OverviewPolyline)
		fmt.Println("Copyrights: ", route.Copyrights)
		fmt.Println("-----------------------------")
		// fmt.Printf("\nRoute %d\n", i+1)
		// fmt.Printf("Distance: %d meters\n", leg.Distance)
		// fmt.Printf("Duration: %d seconds\n", leg.Duration)
		// fmt.Printf("Polyline length: %d\n", len(route.OverviewPolyline))
		// fmt.Printf("Copyrights: %s\n", route.CopyRights)
	}
}

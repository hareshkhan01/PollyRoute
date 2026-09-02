package main

import (
	"log"

	"github.com/hareshkhan01/PollyRoute/internals/application"
	"github.com/hareshkhan01/PollyRoute/internals/config"
	"github.com/hareshkhan01/PollyRoute/internals/handlers"
	"github.com/hareshkhan01/PollyRoute/internals/provider/google/aqi"
	googleWeather "github.com/hareshkhan01/PollyRoute/internals/provider/google/weather"
	"github.com/hareshkhan01/PollyRoute/internals/provider/ola"
	"github.com/hareshkhan01/PollyRoute/internals/router"
	aqiService "github.com/hareshkhan01/PollyRoute/internals/service/aqi"
	routeService "github.com/hareshkhan01/PollyRoute/internals/service/route"
	"github.com/hareshkhan01/PollyRoute/internals/service/scoring"
	"github.com/hareshkhan01/PollyRoute/internals/service/segmentation"
	weatherService "github.com/hareshkhan01/PollyRoute/internals/service/weather"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// Providers
	olaClient, err := ola.NewClient(
		cfg.OLA_API_KEY,
		cfg.OLA_DIRECTIONS_URL,
	)
	if err != nil {
		log.Fatalf("create OLA client: %v", err)
	}

	weatherClient, err := googleWeather.NewClient(
		cfg.GOOGLE_API_KEY,
		cfg.GOOGLE_WEATHER_URL,
	)
	if err != nil {
		log.Fatalf("create weather client: %v", err)
	}

	aqiClient, err := aqi.NewClient(
		cfg.GOOGLE_API_KEY,
		cfg.GOOGLE_AQI_URL,
	)
	if err != nil {
		log.Fatalf("create AQI client: %v", err)
	}

	// Services
	directionService := routeService.NewDirectionService(olaClient)
	weatherSvc := weatherService.NewWeatherService(weatherClient)
	aqiSvc := aqiService.NewAqiService(aqiClient)
	segmentationSvc := segmentation.NewSegmentationService()
	scoringSvc := scoring.NewScoringService()

	// Application service
	routeAnalysisService := application.NewRouteAnalyzeService(
		directionService,
		weatherSvc,
		aqiSvc,
		segmentationSvc,
		scoringSvc,
	)

	// Handler
	routeHandler := handlers.NewRouteHandler(
		routeAnalysisService,
	)

	// Router
	r := router.SetupRouter(routeHandler)
	port := cfg.PORT
	// Server
	log.Println("PollyRoute server running on :", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("start server: %v", err)
	}
}

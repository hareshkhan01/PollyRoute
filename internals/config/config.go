package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT               string
	OLA_API_KEY        string
	OLA_DIRECTIONS_URL string
	GOOGLE_API_KEY     string
	GOOGLE_WEATHER_URL string
	GOOGLE_AQI_URL     string
	DATABASE_URL       string
	JWT_SECRET         string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Warning: .env not found")
	}

	port := os.Getenv("PORT")
	olaApiKey := os.Getenv("OLA_API_KEY")
	olaDirectionsUrl := os.Getenv("OLA_DIRECTIONS_URL")
	googleApiKey := os.Getenv("GOOGLE_API_KEY")
	googleWeatherUrl := os.Getenv("GOOGLE_WEATHER_URL")
	googleAqiUrl := os.Getenv("GOOGLE_AQI_URL")
	databseUrl := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	if port == "" {
		return nil, fmt.Errorf("Failed to load PORT number from .env")
	}
	if olaApiKey == "" {
		return nil, fmt.Errorf("Failed to load OLA_API_KEY from .env")
	}
	if olaDirectionsUrl == "" {
		return nil, fmt.Errorf("Failed to load OLA_DIRECTIONS_URL from .env")
	}
	if googleApiKey == "" {
		return nil, fmt.Errorf("Failed to load GOOGLE_API_KEY from .env")
	}
	if googleWeatherUrl == "" {
		return nil, fmt.Errorf("Failed to load GOOGLE_WEATHER_URLq from .env")
	}
	if googleAqiUrl == "" {
		return nil, fmt.Errorf("Failed to load GOOGLE_AQI_URL from .env")
	}
	if databseUrl == "" {
		return nil, fmt.Errorf("Failed to load DATABASE_URL from .env")
	}
	if jwtSecret == "" {
		return nil, fmt.Errorf("Failed to load JWT_SECRET from .env")
	}
	return &Config{
		OLA_API_KEY:        olaApiKey,
		OLA_DIRECTIONS_URL: olaDirectionsUrl,
		GOOGLE_API_KEY:     googleApiKey,
		GOOGLE_WEATHER_URL: googleWeatherUrl,
		GOOGLE_AQI_URL:     googleAqiUrl,
		PORT:               port,
		DATABASE_URL:       databseUrl,
		JWT_SECRET:         jwtSecret,
	}, nil

}

package scoring

import (
	"github.com/hareshkhan01/PollyRoute/internals/domain"
)

type ScoringService struct {
}

func NewScoringService() *ScoringService {
	return &ScoringService{}
}

func (s *ScoringService) CalculatePollutionScore(segments []domain.RouteSegment) float64 {
	if len(segments) == 0 {
		return 0.0
	}
	weightedAqi := 0.0
	totalDistance := 0.0
	for _, segment := range segments {
		if segment.AqiAndPollutants == nil {
			continue
		}
		aqi := float64(segment.AqiAndPollutants.Aqi)
		distance := float64(segment.Distance)

		weightedAqi += aqi * distance
		totalDistance += distance
	}

	if totalDistance == 0 {
		return 0
	}
	return weightedAqi / totalDistance
}

func (s *ScoringService) CalculateTrafficScore(trafficRanges []domain.TrafficRange) float64 {
	congestion := 0.0
	for _, trafficRange := range trafficRanges {
		congestion += float64(trafficRange.Congestion)
	}

	return congestion
}
func (s *ScoringService) CalculateWeatherScore(segments []domain.RouteSegment) float64 {
	if len(segments) == 0 {
		return 0
	}
	weightedWeatherScore := 0.0
	totalDistance := 0.0
	for _, segment := range segments {
		if segment.Weather == nil {
			continue
		}
		segmentScore := segmentWeatherSecore(segment.Weather)
		distance := segment.Distance

		weightedWeatherScore += segmentScore * distance
		totalDistance += distance
	}
	if totalDistance == 0 {
		return 0
	}
	return weightedWeatherScore / totalDistance
}

func segmentWeatherSecore(weather *domain.Weather) float64 {

	if weather == nil {
		return 0
	}

	conditionScore := calculateWeatherConditionScore(weather.WeatherConditionType)
	rainScore := calculateRainScore(weather.PrecipitationQpf)
	visibilityScore := calculateVisibilityScore(weather.Visibility)
	windScore := calculateWindScore(weather.WindSpeed, weather.WindGust)
	temperatureScore := calculateTemperatureScore(weather.Temperature)
	thunderScore := float64(weather.ThunderstromProbability)

	return rainScore*0.30 +
		visibilityScore*0.25 +
		thunderScore*0.20 +
		windScore*0.15 +
		temperatureScore*0.05 +
		conditionScore*0.05

}
func calculateTemperatureScore(temperature float32) float64 {
	switch {
	case temperature >= 18 && temperature <= 30:
		return 0
	case temperature <= 35:
		return 20
	case temperature <= 40:
		return 50
	case temperature <= 45:
		return 80
	case temperature > 45:
		return 100
	default:
		return 30
	}
}
func calculateWindScore(speed, gust float32) float64 {
	maxWind := speed

	if gust > speed {
		maxWind = gust
	}

	switch {
	case maxWind < 20:
		return 0
	case maxWind < 30:
		return 20
	case maxWind < 40:
		return 40
	case maxWind < 50:
		return 60
	case maxWind < 60:
		return 80
	default:
		return 100
	}
}
func calculateVisibilityScore(visibility float32) float64 {
	switch {
	case visibility > 10:
		return 0
	case visibility > 8:
		return 10
	case visibility > 6:
		return 20
	case visibility > 4:
		return 40
	case visibility > 2:
		return 70
	default:
		return 100
	}
}
func calculateRainScore(qpf float32) float64 {
	switch {
	case qpf <= 0:
		return 0
	case qpf <= 1:
		return 10
	case qpf <= 2.5:
		return 25
	case qpf <= 5:
		return 50
	case qpf <= 10:
		return 75
	default:
		return 100
	}
}

func calculateWeatherConditionScore(condition string) float64 {
	switch condition {

	case "CLEAR",
		"MOSTLY_CLEAR",
		"PARTLY_CLOUDY",
		"MOSTLY_CLOUDY",
		"CLOUDY":
		return 0

	case "LIGHT_RAIN",
		"LIGHT_RAIN_SHOWERS",
		"CHANCE_OF_SHOWERS",
		"SCATTERED_SHOWERS":
		return 20

	case "RAIN",
		"RAIN_SHOWERS",
		"LIGHT_TO_MODERATE_RAIN":
		return 40

	case "HEAVY_RAIN",
		"HEAVY_RAIN_SHOWERS",
		"MODERATE_TO_HEAVY_RAIN",
		"RAIN_PERIODICALLY_HEAVY":
		return 70

	case "THUNDERSTORM",
		"THUNDERSHOWER",
		"LIGHT_THUNDERSTORM_RAIN",
		"SCATTERED_THUNDERSTORMS",
		"HEAVY_THUNDERSTORM":
		return 100

	case "WINDY":
		return 40

	case "WIND_AND_RAIN":
		return 70

	case "HAIL",
		"HAIL_SHOWERS":
		return 100

	default:
		return 0
	}
}

func (s *ScoringService) CalculateBalancedScore(weatherScore, pollutionScore, trafficScore float64) float64 {
	return (weatherScore + pollutionScore + trafficScore) / 3
}

package aqi

/*
dateTime
indexes[].code
indexes[].aqi
indexes[].category
indexes[].dominantPollutant

pollutants[].code
pollutants[].concentration.value
pollutants[].concentration.units
*/
/*
type AQIResponse struct {
	DateTime   string      `json:"dateTime"`
	AqiIndexes []AqiIndex  `json:"indexes"`
	Pollutants []Pollutant `json:"pollutants"`
}

type AqiIndex struct {
	Code              string  `json:"code"`
	Aqi               float32 `json:"aqi"`
	Category          string  `json:"category"`
	DominantPollutant string  `json:"dominantPollutant"`
}

type Pollutant struct {
	Code          string `json:"code"`
	Concentration struct {
		Value float32 `json:"value"`
		Units string  `json:"units"`
	} `json:"concentration"`
}
*/

/*
type AQIAndPollutants struct {
	DateTime          time.Time
	AqiIndexCode      string
	Aqi               float32
	Category          string
	DominantPollutant string
	Pollutants        []Pollutant
}

type Pollutant struct {
	Code               string
	ConcentrationValue float32
	ConcentrationUnit  string
}
*/

import (
	"fmt"

	"github.com/hareshkhan01/PollyRoute/internals/domain"

	"github.com/hareshkhan01/PollyRoute/pkg/datetime"
)

func MapResponse(response *AQIResponse) (*domain.AQIAndPollutants, error) {
	if response == nil {
		return nil, fmt.Errorf("Aqi Reponse Can not be null")
	}
	if len(response.AqiIndexes) == 0 {
		return nil, fmt.Errorf("AQI Indexes are empty!")
	}
	dateTime := response.DateTime
	// Aqi
	aqi := response.AqiIndexes[0].Aqi
	aqiIndexCode := response.AqiIndexes[0].Code
	category := response.AqiIndexes[0].Category
	dominantPollutant := response.AqiIndexes[0].DominantPollutant
	// pollutants
	pollutants := make([]domain.Pollutant, 0, len(response.Pollutants))

	for _, pollutant := range response.Pollutants {
		pollutants = append(pollutants, domain.Pollutant{
			Code:               pollutant.Code,
			ConcentrationValue: pollutant.Concentration.Value,
			ConcentrationUnit:  pollutant.Concentration.Units,
		})
	}
	dateTimeObj, err := datetime.StringToDateTime(dateTime)
	if err != nil {
		return nil, fmt.Errorf("Cannot convert date string to object! %w", err)
	}
	return &domain.AQIAndPollutants{
		DateTime:          dateTimeObj,
		AqiIndexCode:      aqiIndexCode,
		Aqi:               aqi,
		Category:          category,
		DominantPollutant: dominantPollutant,
		Pollutants:        pollutants,
	}, nil
}

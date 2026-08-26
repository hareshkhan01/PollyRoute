package domain

import (
	"time"
)

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

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

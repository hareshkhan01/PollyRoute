package weather

/*
temperature
relativeHumidity
wind.speed
wind.direction
precipitation.qpf
airPressure.meanSeaLevelMillibars
weatherCondition.type

wind.gust
visibility
cloudCover
thunderstormProbability
*/
type WeatherResponse struct {
	Temperature      Temperature   `json:"temperature"`
	Visibility       Visibility    `json:"visibility"`
	RelativeHumidity float32       `json:"relativeHumidity"`
	Wind             Wind          `json:"wind"`
	Precipitation    Precipitation `json:"precipitation"`
	AirPressure      Airpressure   `json:"airPressure"`
	WeatherCondition struct {
		Type string `json:"type"`
	} `json:"weatherCondition"`
	CloudCover              float32 `json:"cloudCover"`
	ThunderstormProbability float32 `json:"thunderstormProbability"`
}

type Temperature struct {
	Unit    string  `json:"unit"`
	Degrees float32 `json:"degrees"`
}
type Visibility struct {
	Unit     string  `json:"unit"`
	Distance float32 `json:"distance"`
}

type Airpressure struct {
	MeanSeaLevelMillibars float32 `json:"meanSeaLevelMillibars"`
}

type Wind struct {
	Direction struct {
		Cardinal string  `json:"cardinal"`
		Degrees  float32 `json:"degrees"`
	} `json:"direction"`

	Speed struct {
		Unit  string  `json:"unit"`
		Value float32 `json:"value"`
	} `json:"speed"`

	Gust struct {
		Unit  string  `json:"unit"`
		Value float32 `json:"value"`
	} `json:"gust"`
}

type Precipitation struct {
	Qpf struct {
		Unit     string  `json:"unit"`
		Quantity float32 `json:"quantity"`
	} `json:"qpf"`
}

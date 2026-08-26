package domain

type Weather struct {
	Temperature             float32
	Visibility              float32
	RelativeHumidity        float32
	WindDirection           float32
	WindSpeed               float32
	WindGust                float32
	PrecipitationQpf        float32
	AirPressure             float32
	WeatherConditionType    string
	CloudCover              float32
	ThunderstromProbability float32
}

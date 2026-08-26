package weather

func MapResponse(WeatherResponse *response)(*domain.Weather,error){
	if response==nil{
		return nil,fmt.Errorf("Response is Null!")
	}
	weather:=domain.Weather{
		Temperature: response.Temperature.Degrees,
		Visibility: response.Visibility.Distance,
		RelativeHumidity: response.RelativeHumidity,
		WindDirection: response.Wind.Direction.Degrees,
		WindSpeed: response.Wind.Speed.Value,
		WindGust: response.Wind.Gust.Value,
		PrecipitationQpf: response.Precipitation.Qpf.Quantity,
		AirPressure: response.AirPressure.MeanSeaLevelMillibars,
		WeatherConditionType: response.WeatherCondition.Type,
		CloudCover: response.CloudCover,
		ThunderstromProbability: response.ThunderstormProbability,
	}
	return &weather,nil
}

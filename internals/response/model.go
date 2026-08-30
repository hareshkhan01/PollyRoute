package response

type AnalysisRouteResponse struct {
	Routes []RouteResponse `json:"routes"`
}

type RouteResponse struct {
	ID           int               `json:"id"`
	Polyline     string            `json:"polyline"`
	Distance     int64             `json:"distance"`
	Duration     int64             `json:"duration"`
	Scores       ScoreResponse     `json:"scores"`
	Traffic      []TrafficResponse `json:"traffic"`
	Segments     []SegmentResponse `json:"segments"`
	IsSuspicious bool              `json:"isSuspicious"`
	Copyright    string            `json:"copyright"`
}

type ScoreResponse struct {
	PollutionScore float64 `json:"pollution"`
	WeatherScore   float64 `json:"weather"`
	TrafficScore   float64 `json:"traffic"`
	BalancedScore  float64 `json:"balanced"`
}

type TrafficResponse struct {
	Start      int `json:"start"`
	End        int `json:"end"`
	Congestion int `json:"congestion"`
}

type SegmentResponse struct {
	Index            int                      `json:"index"`
	Distance         float64                  `json:"distance"`
	Midpoint         ResponseCoordinate       `json:"midpoint"`
	AqiAndPollutants AQIAndPollutantsResponse `json:"aqiAndPollutants"`
	Weather          WeatherResponse          `json:"weather"`
}

type ResponseCoordinate struct {
	Lat float64
	Lng float64
}

type AQIAndPollutantsResponse struct {
	Aqi               float64 `json:"aqi"`
	DominantPollutant string  `json:"dominantPollutant"`
}

type WeatherResponse struct {
	Temperature             float64 `json:"temperature"`
	ThunderstormProbability float64 `json:"thunderstormProbability"`
}

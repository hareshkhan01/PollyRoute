package ola

type Response struct {
	Status string  `json:"status"`
	Routes []Route `json:"routes"`
}

type Route struct {
	Legs           []Leg  `json:"legs"`
	Polyline       string `json:"overview_polyline"`
	TravelAdvisory string `json:"travel_advisory"`
	Copyrights     string `json:"copyrights"`
}

type Leg struct {
	Steps    []Step `json:"steps"`
	Distance int64  `json:"distance"`
	Duration int64  `json:"duration"`
}

type Step struct {
	Maneuver string `json:"maneuver"`
}

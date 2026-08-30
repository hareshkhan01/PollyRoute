package request

type AnalyzeRouteRequest struct {
	Origin      CoordianteRequest `json:"origin"`
	Destinition CoordianteRequest `json:"destination"`
}

type CoordianteRequest struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

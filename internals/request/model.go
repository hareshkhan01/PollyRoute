package request

type AnalyzeRouteRequest struct {
	Origin      CoordinateRequest `json:"origin" binding:"required"`
	Destination CoordinateRequest `json:"destination" binding:"required"`
}

type CoordinateRequest struct {
	Lat float64 `json:"latitude" binding:"required,gte=-90,lte=90"`
	Lng float64 `json:"longitude" binding:"required,gte=-180,lte=180"`
}

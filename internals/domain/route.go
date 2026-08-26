package domain

type Routes struct {
	Routes []Route
}

type Route struct {
	TotalDistance int64
	TotalDuration int64
	Coordinates   []Coordinate
	TrafficRanges []TrafficRange
	IsSuspicious  bool
	Copyrights    string
}

type Coordinate struct {
	Lat float64
	Lng float64
}

type TrafficRange struct {
	Start      int
	End        int
	Congestion int
}

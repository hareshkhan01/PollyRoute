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

type RouteSegment struct {
	Index       int
	Coordinates []Coordinate
	Distance    float64
	Midpoint    Coordinate
}

/*
Example RouteSgement
Segement S1(C1 to C5)
Respresentation
Index 0
Coordinates[]{C1,C2,C3,C4,C5}
Distance 5300
Midpoint - C4
*/

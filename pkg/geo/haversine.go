package geo

import (
	"math"

	"github.com/hareshkhan01/PollyRoute/internals/domain"
)

// const R = 6173.0 // Earth Radius In KM
const R = 6_371_000 // Earth Radius In Meters as OLA API returns distance in meter
func Haversine(origin domain.Coordinate, destination domain.Coordinate) float64 {

	// Convert Latitude and Longitude Degrees to Raidus

	// Latitude
	phi1 := origin.Lat * math.Pi / 180
	phi2 := destination.Lat * math.Pi / 180

	// Longitude
	lambda1 := origin.Lng * math.Pi / 180
	lambda2 := destination.Lng * math.Pi / 180

	deltaPhi := phi2 - phi1
	deltaLambda := lambda2 - lambda1

	a := math.Pow(math.Sin(deltaPhi/2.0), 2.0) + math.Cos(phi1)*math.Cos(phi2)*math.Pow(math.Sin(deltaLambda/2.0), 2.0)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c

	return d
}

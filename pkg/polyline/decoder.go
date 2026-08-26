package polyline

import (
	"fmt"

	"github.com/hareshkhan01/PollyRoute/internals/domain"
)

func Decode(encoded string) ([]domain.Coordinate, error) {
	var lat int
	var lng int

	var coordinates []domain.Coordinate

	for i := 0; i < len(encoded); {

		dLat, next, err := decodeValue(encoded, i)
		if err != nil {
			return nil, fmt.Errorf("Failed to Decode Lat. %w", err)
		}
		i = next
		lat += dLat

		dLng, next, err := decodeValue(encoded, i)
		if err != nil {
			return nil, fmt.Errorf("Failed to Decode Lng. %w", err)
		}
		i = next
		lng += dLng

		coordinates = append(coordinates,
			domain.Coordinate{
				Lat: float64(lat) / 1e5,
				Lng: float64(lng) / 1e5,
			},
		)

	}
	return coordinates, nil
}

func decodeValue(encoded string, index int) (int, int, error) {
	var result int
	var shift uint

	for {
		if index >= len(encoded) {
			return 0, index, fmt.Errorf("Unexpected end of polyline/encoded ")
		}
		b := int(encoded[index]) - 63
		index++
		result |= (b & 0x1f) << shift
		if b < 0x20 {
			break
		}
		shift += 5
		if shift > 30 {
			return 0, index, fmt.Errorf("Invalid polyline value can't be greater than 31 bit")
		}
	}
	if result&1 != 0 {
		return ^(result >> 1), index, nil
	}

	return result >> 1, index, nil

}

package util

import (
	"math"

	"github.com/rafihayne/ch/pkg/graph"
)

// Copy pasted from https://github.com/umahmood/haversine
// MIT lisenced

const (
	// earthRadiusMi = 3958 // radius of the earth in miles.
	earthRaidusKm = 6371 // radius of the earth in kilometers.
)

// degreesToRadians converts from degrees to radians.
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// Distance calculates the shortest path between two coordinates on the surface
// of the Earth. This function returns distance in meters
func Haversine(lhs, rhs graph.NodeValue) float64 {
	lat1 := degreesToRadians(lhs.Y)
	lon1 := degreesToRadians(lhs.X)
	lat2 := degreesToRadians(rhs.Y)
	lon2 := degreesToRadians(rhs.X)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// mi = c * earthRadiusMi
	km := c * earthRaidusKm

	// Convert to meters
	return km * 1000
}

func Euclidean(lhs, rhs graph.NodeValue) float64 {
	return math.Sqrt(math.Pow(lhs.X-rhs.X, 2) + math.Pow(lhs.Y-rhs.Y, 2))
}

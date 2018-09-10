package geo

import (
	"math"
	"strconv"

	"googlemaps.github.io/maps"
)

const EARTH_RADIUS = 6371000 // m
var DEGREES_TO_RADIANS = math.Pi / 180.0
var RADIANS_TO_DEGREES = 180.0 / math.Pi

// LatLng represents a latitude,longitude pair.
type LatLng = maps.LatLng

func Coordinate(coord float64) string {
	return strconv.FormatFloat(coord, 'f', -1, 64)
}

// ParseLatLng will parse a string representation of a 'lat,lng' pair.
func ParseLatLng(s string) (LatLng, error) {
	return maps.ParseLatLng(s)
}

// ParseLatLngs parses a string of | separated 'lat,lng' pairs.
func ParseLatLngs(s string) ([]LatLng, error) {
	return maps.ParseLatLngList(s)
}

// Polyline represents a list of lat,lng points encoded as a byte array.
// See: https://developers.google.com/maps/documentation/utilities/polylinealgorithm
type Polyline = maps.Polyline

// DecodePolyline converts a polyline encoded string to an array of LatLng objects
func DecodePolyline(s string) ([]LatLng, error) {
	return maps.DecodePolyline(s)
}

// EncodePolyline returns a new encoded Polyline from the given points.
func EncodePolyline(lls []LatLng) string {
	return maps.Encode(lls)
}

// Center returns a LatLng object representing the center of the given point.
//func Center(lls []LatLng) LatLng {}

// CenterOfMinimumDistance is an alias for Center
//var CenterOfMinimumDistance = Center

// Average returns a LatLng object representing the average of the given points.
func Average(lls []LatLng) LatLng {
	if len(lls) == 0 {
		return LatLng{}
	} else if len(lls) == 1 {
		return lls[0]
	}

	x := 0.0
	y := 0.0
	z := 0.0

	for _, ll := range lls {
		lat := ll.Lat * DEGREES_TO_RADIANS
		lng := ll.Lng * DEGREES_TO_RADIANS

		x += math.Cos(lat) * math.Cos(lng)
		y += math.Cos(lat) * math.Sin(lng)
		z += math.Sin(lat)
	}

	tot := float64(len(lls))
	x /= tot
	y /= tot
	z /= tot

	lng := math.Atan2(y, x)
	lat := math.Atan2(z, math.Sqrt(x*x+y*y))

	return LatLng{Lat: lat * RADIANS_TO_DEGREES, Lng: lng * RADIANS_TO_DEGREES}
}

// GeographicMidpoint is an alias for Average.
var GeographicMidpoint = Average

// Centroid is an alias for Average.
var Centroid = Average

// Distance calculates the Haversine distance between two points in metres.
// See: http://www.movable-type.co.uk/scripts/latlong.html
func Distance(p1, p2 LatLng) float64 {
	dLat := (p2.Lat - p1.Lat) * DEGREES_TO_RADIANS
	dLon := (p2.Lng - p1.Lng) * DEGREES_TO_RADIANS

	lat1 := p1.Lat * DEGREES_TO_RADIANS
	lat2 := p2.Lat * DEGREES_TO_RADIANS

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}

// HaversineDistance is an alias for Distance.
var HaversineDistance = Distance

// GreatCircleDistance is an alias for Distance.
var GreatCircleDistance = Distance

// Bearing computes the bearing/direction to travel from p1 to p2 in degrees.
// See: http://www.movable-type.co.uk/scripts/latlong.html
func Bearing(p1, p2 LatLng) float64 {
	dLon := (p2.Lng - p1.Lng) * DEGREES_TO_RADIANS

	lat1 := p1.Lat * DEGREES_TO_RADIANS
	lat2 := p2.Lat * DEGREES_TO_RADIANS

	y := math.Sin(dLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) -
		math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLon)
	brng := math.Atan2(y, x) * RADIANS_TO_DEGREES

	return normalizeBearing(brng)
}

// Direction is an alias for Bearing.
var Direction = Bearing

// AverageBearing computers the mean bearing for the set of points pts.
// See: https://en.wikipedia.org/wiki/Mean_of_circular_quantities
func AverageBearing(pts []LatLng) float64 {
	if len(pts) <= 1 {
		return 0.0
	}

	x := 0.0
	y := 0.0
	pt := pts[0]
	for i := 1; i < len(pts); i++ {
		a := Bearing(pt, pts[i])

		x += math.Cos(a * DEGREES_TO_RADIANS)
		y += math.Sin(a * DEGREES_TO_RADIANS)

		pt = pts[i]
	}
	return normalizeBearing(math.Atan2(y, x) * RADIANS_TO_DEGREES)
}

// AverageDirection is an alias for AverageBearing.
var AverageDirection = AverageBearing

func normalizeBearing(b float64) float64 {
	return b + math.Ceil(-b/360)*360
}

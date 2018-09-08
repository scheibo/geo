package geo

import (
	"googlemaps.github.io/maps"
)

// LatLng represents a latitude,longitude pair.
type LatLng = maps.LatLng

// ParseLatLng will parse a string representation of a 'lat,lng' pair.
func ParseLatLng(s string) (LatLng, error) {
	return maps.ParseLatLng(s)
}

// ParseLatLngs parses a string of | separated 'lat,lng' pairs.
func ParseLatLngs(s string) ([]LatLng, error) {
	return maps.ParseLatLngList(s)
}

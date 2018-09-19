package geo

import (
	"context"
	"os"

	"googlemaps.github.io/maps"
)

// MAX_LOCATIONS_PER_REQUEST is the maximum number of 'lat,lng' pairs the Google
// Maps API allows in a single request.
const MAX_LOCATIONS_PER_REQUEST = 512

// Client may be used to access methods which rely upon the Google Maps API.
type Client struct {
	client *maps.Client
}

// NewClient returns a new client that can access the Google Maps API.
func NewClient(key ...string) (*Client, error) {
	k := os.Getenv("GOOGLE_MAPS_API_KEY")
	if len(key) > 0 && key[0] != "" {
		k = key[0]
	}

	c, err := maps.NewClient(maps.WithAPIKey(k))
	if err != nil {
		return nil, err
	}

	return &Client{client: c}, nil
}

// Elevation determines the elevation for each 'lat,lng' pair in ll.
func (c *Client) Elevation(lls []LatLng) ([]LatLngEle, error) {
	var lles []LatLngEle

	per := MAX_LOCATIONS_PER_REQUEST
	num := len(lls) / per

	var i int
	var err error
	for ; i < num; i++ {
		lles, err = c.fillElevation(lles, lls[i*per:(i+1)*per])
		if err != nil {
			return []LatLngEle{}, err
		}
	}

	if len(lls)%per != 0 {
		lles, err = c.fillElevation(lles, lls[i*per:])
		if err != nil {
			return []LatLngEle{}, err
		}
	}

	return lles, nil
}

func (c *Client) fillElevation(lles []LatLngEle, lls []LatLng) ([]LatLngEle, error) {
	r := maps.ElevationRequest{Locations: toMaps(lls)}
	res, err := c.client.Elevation(context.Background(), &r)
	if err != nil {
		return []LatLngEle{}, err
	}
	return fromMaps(lles, res), nil
}

func toMaps(lls []LatLng) []maps.LatLng {
	var mll []maps.LatLng
	for _, ll := range lls {
		mll = append(mll, maps.LatLng{Lat: ll.Lat, Lng: ll.Lng})
	}
	return mll
}

func fromMaps(lles []LatLngEle, res []maps.ElevationResult) []LatLngEle {
	for _, r := range res {
		lles = append(lles, LatLngEle{Lat: r.Location.Lat, Lng: r.Location.Lng, Ele: r.Elevation})
	}
	return lles
}

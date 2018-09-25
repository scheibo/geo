package geo

import (
	"bytes"
	"io"
	"math"
	"strconv"
	"strings"
)

const EARTH_RADIUS = 6371008.8 // m
var DEGREES_TO_RADIANS = math.Pi / 180.0
var RADIANS_TO_DEGREES = 180.0 / math.Pi

// LatLng represents a 'latitude,longitude' pair.
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// LatLngEle represents a 'latitude,longitude,elevation' triple.
type LatLngEle struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	Ele float64 `json:"ele",omitempty`
}

func Coordinate(coord float64) string {
	return strconv.FormatFloat(coord, 'f', -1, 64)
}

// ParseLatLng will parse a string representation of a 'lat,lng' pair.
func ParseLatLng(s string) (LatLng, error) {
	ll := strings.Split(s, ",")
	lat, err := strconv.ParseFloat(ll[0], 64)
	if err != nil {
		return LatLng{}, err
	}
	lng, err := strconv.ParseFloat(ll[1], 64)
	if err != nil {
		return LatLng{}, err
	}
	return LatLng{Lat: lat, Lng: lng}, nil
}

// ParseLatLngEle will parse a string representation of a 'lat,lng,ele' triples.
func ParseLatLngEle(s string) (LatLngEle, error) {
	lle := strings.Split(s, ",")
	lat, err := strconv.ParseFloat(lle[0], 64)
	if err != nil {
		return LatLngEle{}, err
	}
	lng, err := strconv.ParseFloat(lle[1], 64)
	if err != nil {
		return LatLngEle{}, err
	}
	ele, err := strconv.ParseFloat(lle[2], 64)
	if err != nil {
		return LatLngEle{}, err
	}
	return LatLngEle{Lat: lat, Lng: lng, Ele: ele}, nil
}

// ParseLatLngs parses a string of | separated 'lat,lng' pairs.
func ParseLatLngs(s string) ([]LatLng, error) {
	result := []LatLng{}

	ls := strings.Split(s, "|")
	for _, l := range ls {
		ll, err := ParseLatLng(l)
		if err != nil {
			return []LatLng{}, err
		}
		result = append(result, ll)
	}
	return result, nil
}

// ParseLatLngEles parses a string of | separated 'lat,lng,ele' triples.
func ParseLatLngEles(s string) ([]LatLngEle, error) {
	result := []LatLngEle{}

	ls := strings.Split(s, "|")
	for _, l := range ls {
		lle, err := ParseLatLngEle(l)
		if err != nil {
			return []LatLngEle{}, err
		}
		result = append(result, lle)
	}
	return result, nil
}

func (ll *LatLng) String() string {
	return Coordinate(ll.Lat) + "," + Coordinate(ll.Lng)
}

func (ll *LatLng) latLngEle() LatLngEle {
	return LatLngEle{Lat: ll.Lat, Lng: ll.Lng, Ele: math.Inf(-1)}
}

func (lle *LatLngEle) String() string {
	return Coordinate(lle.Lat) + "," +
		Coordinate(lle.Lng) + "," +
		Coordinate(lle.Ele)
}

func (lle *LatLngEle) LatLng() LatLng {
	return LatLng{Lat: lle.Lat, Lng: lle.Lng}
}

// Polyline represents a list of lat,lng points encoded as a byte array.
// See: https://developers.google.com/maps/documentation/utilities/polylinealgorithm
type Polyline struct {
	Points string `json:"points"`
}

// ZPolyline represents a list of lat,lng,ele points encoded as a byte array,
// extending the standard polyline encoding algorithm.
type ZPolyline struct {
	Points string `json:"points"`
}

// DecodePolyline converts a polyline encoded string to an array of LatLng objects
func DecodePolyline(s string) ([]LatLng, error) {
	p := &Polyline{
		Points: s,
	}
	return p.Decode()
}

// DecodeZPolyline converts a z-polyline encoded string to an array of LatLngEle objects
func DecodeZPolyline(s string) ([]LatLngEle, error) {
	p := &ZPolyline{
		Points: s,
	}
	return p.Decode()
}

// Decode converts this encoded Polyline to an array of LatLng objects.
func (p *Polyline) Decode() ([]LatLng, error) {
	input := bytes.NewBufferString(p.Points)

	var lat, lng int64
	lls := make([]LatLng, 0, len(p.Points)/2)
	for {
		dlat, _ := decodeInt(input)
		dlng, err := decodeInt(input)
		if err == io.EOF {
			return lls, nil
		}
		if err != nil {
			return nil, err
		}

		lat, lng = lat+dlat, lng+dlng
		lls = append(lls, LatLng{
			Lat: float64(lat) * 1e-5,
			Lng: float64(lng) * 1e-5,
		})
	}
}

// Decode converts this encoded ZPolyline to an array of LatLngEle objects.
func (p *ZPolyline) Decode() ([]LatLngEle, error) {
	input := bytes.NewBufferString(p.Points)

	var lat, lng, ele int64
	lles := make([]LatLngEle, 0, len(p.Points)/3)
	for {
		dlat, _ := decodeInt(input)
		dlng, _ := decodeInt(input)
		dele, err := decodeInt(input)
		if err == io.EOF {
			return lles, nil
		}
		if err != nil {
			return nil, err
		}

		lat, lng, ele = lat+dlat, lng+dlng, ele+dele
		lles = append(lles, LatLngEle{
			Lat: float64(lat) * 1e-5,
			Lng: float64(lng) * 1e-5,
			Ele: float64(ele) * 1e-5,
		})
	}
}

// EncodePolyline returns a new encoded Polyline from the given LatLng points.
func EncodePolyline(lls []LatLng) string {
	var prevLat, prevLng int64

	out := new(bytes.Buffer)
	out.Grow(len(lls) * 4)

	for _, ll := range lls {
		lat := int64(ll.Lat * 1e5)
		lng := int64(ll.Lng * 1e5)

		encodeInt(lat-prevLat, out)
		encodeInt(lng-prevLng, out)

		prevLat, prevLng = lat, lng
	}

	return out.String()
}

// EncodeZPolyline returns a new encoded ZPolyline from the given LatLngEle points.
func EncodeZPolyline(lles []LatLngEle) string {
	var prevLat, prevLng, prevEle int64

	out := new(bytes.Buffer)
	out.Grow(len(lles) * 4)

	for _, lle := range lles {
		lat := int64(lle.Lat * 1e5)
		lng := int64(lle.Lng * 1e5)
		ele := int64(lle.Ele * 1e5)

		encodeInt(lat-prevLat, out)
		encodeInt(lng-prevLng, out)
		encodeInt(ele-prevEle, out)

		prevLat, prevLng, prevEle = lat, lng, ele
	}

	return out.String()
}

// decodeInt reads an encoded int64 from the passed io.ByteReader.
func decodeInt(r io.ByteReader) (int64, error) {
	result := int64(0)
	var shift uint8

	for {
		raw, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		b := raw - 63
		result += int64(b&0x1f) << shift
		shift += 5

		if b < 0x20 {
			bit := result & 1
			result >>= 1
			if bit != 0 {
				result = ^result
			}
			return result, nil
		}
	}
}

// encodeInt writes an encoded int64 to the passed io.ByteWriter.
func encodeInt(v int64, w io.ByteWriter) {
	if v < 0 {
		v = ^(v << 1)
	} else {
		v <<= 1
	}
	for v >= 0x20 {
		w.WriteByte((0x20 | (byte(v) & 0x1f)) + 63)
		v >>= 5
	}
	w.WriteByte(byte(v) + 63)
}

// Center returns a LatLng object representing the center of the given point.
//func Center(lls []LatLng) LatLng {}

// CenterOfMinimumDistance is an alias for Center
//var CenterOfMinimumDistance = Center

// Average returns a LatLng object representing the average of the given points.
func Average(lls []LatLng) LatLng {
	lles := make([]LatLngEle, len(lls))
	for i := 0; i < len(lls); i++ {
		lles[i] = lls[i].latLngEle()
	}
	avg := AverageZ(lles)
	return avg.LatLng()
}

// GeographicMidpoint is an alias for Average.
var GeographicMidpoint = Average

// Centroid is an alias for Average.
var Centroid = Average

// AverageZ returns a LatLngEle object representing the average of the given points.
func AverageZ(lles []LatLngEle) LatLngEle {
	if len(lles) == 0 {
		return LatLngEle{}
	} else if len(lles) == 1 {
		return lles[0]
	}

	x := 0.0
	y := 0.0
	z := 0.0

	e := 0.0

	for _, lle := range lles {
		lat := lle.Lat * DEGREES_TO_RADIANS
		lng := lle.Lng * DEGREES_TO_RADIANS

		x += math.Cos(lat) * math.Cos(lng)
		y += math.Cos(lat) * math.Sin(lng)
		z += math.Sin(lat)

		e += lle.Ele
	}

	tot := float64(len(lles))
	x /= tot
	y /= tot
	z /= tot
	e /= tot

	lng := math.Atan2(y, x)
	lat := math.Atan2(z, math.Sqrt(x*x+y*y))

	return LatLngEle{Lat: lat * RADIANS_TO_DEGREES, Lng: lng * RADIANS_TO_DEGREES, Ele: e}
}

// GeographicMidpointZ is an alias for AverageZ.
var GeographicMidpointZ = AverageZ

// CentroidZ is an alias for AverageZ.
var CentroidZ = AverageZ

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

// AverageBearing computes the mean bearing for the set of points pts.
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

// AverageBearingZ computes the mean bearing for the set of points pts.
func AverageBearingZ(lles []LatLngEle) float64 {
	return AverageBearing(LatLngs(lles))
}

// LatLngs converts lles from LatLngEles to LatLngs.
func LatLngs(lles []LatLngEle) []LatLng {
	lls := make([]LatLng, len(lles))
	for i := 0; i < len(lles); i++ {
		lls[i] = lles[i].LatLng()
	}
	return lls
}

// AverageDirection is an alias for AverageBearing.
var AverageDirectionZ = AverageBearingZ

func normalizeBearing(b float64) float64 {
	return b + math.Ceil(-b/360)*360
}

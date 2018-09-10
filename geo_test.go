package geo

import (
	"math"
	"testing"
)

const eps = 0.0001

const poly = "kbhcF`_ciV|Bt@nBDrA_@dG}CfASvA@~Af@`Av@rEnF~@|@`@PtALb@C|@[t@w@v@_BlG{KhE}GbAqA`A}A~AyBhAsAvAyArAcArGsE|CmB~DkCjDoCnCoCvAyA~@iAbAcAtA_BT]d@kAb@wC^_DXsCh@}IZwAtBiEbEkHv@qAbAoAjCaC~FcDvA{@`CmApDqB|@g@j@c@r@cAh@wA`@uBBa@EuAo@aGIyA_@}DEy@J_Cf@kBlEeI|@yB`@yBJmBJiGReBXmAn@}Ax@wAxCkEpAwBAKc@ZWB"

func TestParseLatLng(t *testing.T) {
	tests := []struct {
		s        string
		expected LatLng
	}{
		{"12.345678,56.789012", LatLng{12.345678, 56.789012}},
	}
	for _, tt := range tests {
		actual, err := ParseLatLng(tt.s)
		if err != nil {
			t.Errorf("ParseLatLng(%s): got %+v, want %+v", tt.s, err, tt.expected)
		}

		if !actual.AlmostEqual(&tt.expected, eps) {
			t.Errorf("ParseLatLng(%s): got %+v, want %+v", tt.s, actual, tt.expected)
		}
	}
}

func TestParseLatLngs(t *testing.T) {
	tests := []struct {
		s        string
		expected []LatLng
	}{
		{"12.34,56.78|14.89,123.89", []LatLng{{12.34, 56.78}, {14.89, 123.89}}},
	}
	for _, tt := range tests {
		actual, err := ParseLatLngs(tt.s)
		if err != nil {
			t.Errorf("ParseLatLngs(%s): got %+v, want %+v", tt.s, err, tt.expected)
		}

		if len(actual) != len(tt.expected) {
			t.Errorf("ParseLatLngs(%s): got %+v, want %+v", tt.s, actual, tt.expected)
		}

		for i, ll := range tt.expected {
			if !actual[i].AlmostEqual(&ll, eps) {
				t.Errorf("ParseLatLngs(%s): got %+v, want %+v", tt.s, actual, tt.expected)
			}
		}
	}
}

func TestDecodePolyline(t *testing.T) {
	tests := []struct {
		s          string
		start, end LatLng
	}{
		{poly, LatLng{37.4021463, -122.2451293}, LatLng{37.3721483, -122.2083962}},
	}
	for _, tt := range tests {
		actual, err := DecodePolyline(tt.s)
		if err != nil {
			t.Errorf("DecodePolyline(%s): got %+v", tt.s, err)
		}

		if len(actual) < 2 {
			t.Errorf("DecodePolyline(%s): got %+v, want len >= 2", tt.s, actual)
		}

		if !actual[0].AlmostEqual(&tt.start, eps) || !actual[len(actual)-1].AlmostEqual(&tt.end, eps) {
			t.Errorf("DecodePolyline(%s): got %+v, want (%+v, ..., %+v)", tt.s, actual, tt.start, tt.end)
		}
	}
}

func TestEncodePolyline(t *testing.T) {
	tests := []struct {
		s string
	}{
		{poly},
	}
	for _, tt := range tests {
		dec, err := DecodePolyline(tt.s)
		if err != nil {
			t.Errorf("DecodePolyline(%s): got %+v", tt.s, err)
		}
		actual := EncodePolyline(dec)
		if actual != tt.s {
			t.Errorf("EncodePolyline(%+v): got %s, want %s", dec, actual, poly)
		}
	}
}

func TestBearing(t *testing.T) {
	tests := []struct {
		a, b     LatLng
		expected float64
	}{
		{LatLng{37.4021463, -122.2451293}, LatLng{37.3721483, -122.2083962}, 135.77},
	}
	for _, tt := range tests {
		actual := Bearing(tt.a, tt.b)
		if !Eqf(actual, tt.expected, eps) {
			t.Errorf("Bearing(%+v, %+v): got %.2f, want %.2f", tt.a, tt.b, actual, tt.expected)
		}
	}
}

func TestAverageBearing(t *testing.T) {
	tests := []struct {
		s        string
		expected float64
	}{
		{poly, 134.43},
	}
	for _, tt := range tests {
		dec, err := DecodePolyline(tt.s)
		if err != nil {
			t.Errorf("DecodePolyline(%s): got %+v", tt.s, err)
		}
		actual := AverageBearing(dec)
		if !Eqf(actual, tt.expected, eps) {
			t.Errorf("AverageBearing(%+v): got %.2f, want %.2f", dec, actual, tt.expected)
		}
	}
}

// Eqf returns true when floats a and b are equal to within some small epsilon eps.
func Eqf(a, b float64, eps ...float64) bool {
	e := 1e-3
	if len(eps) > 0 {
		e = eps[0]
	}
	// min is the smallest normal value possible
	const min = float64(2.2250738585072014E-308) // 1 / 2**(1022)

	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)

	if a == b {
		return true
	} else if a == b || b == 0 || diff < min {
		// a or b is zero or both are extremely close to it relative error is less meaningful here
		return diff < (e * min)
	} else {
		// use relative error
		return diff/(absA+absB) < e
	}
}

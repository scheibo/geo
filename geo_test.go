package geo

import "testing"

const eps = 0.0001

func TestParseLatLng(t *testing.T) {
	tests := []struct {
		s        string
		expected LatLng
	}{
		{"12.345678,56.789012", LatLng{Lat: 12.345678, Lng: 56.789012}},
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
		{"12.34,56.78|14.89,123.89", []LatLng{{Lat: 12.34, Lng: 56.78}, {Lat: 14.89, Lng: 123.89}}},
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

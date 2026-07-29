package repository

import "testing"

func TestParseMeliShipmentID(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want int64
	}{
		{"guide id numerico", "47626453039", 47626453039},
		{"con espacios", "  47626453039 ", 47626453039},
		{"tracking con prefijo", "MEL47626453039FMDOF01", 0},
		{"vacio", "", 0},
		{"cero", "0", 0},
		{"negativo", "-5", 0},
		{"no numerico", "ABC", 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseMeliShipmentID(tc.raw); got != tc.want {
				t.Fatalf("parseMeliShipmentID(%q) = %d, want %d", tc.raw, got, tc.want)
			}
		})
	}
}

package usecasemessaging

import (
	"testing"
)

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name        string
		phone       string
		wantErr     bool
		errContains string
	}{
		// Casos válidos
		{
			name:    "número colombiano con prefijo +",
			phone:   "+573001234567",
			wantErr: false,
		},
		{
			name:    "número colombiano con prefijo 00",
			phone:   "00573001234567",
			wantErr: false,
		},
		{
			name:    "número colombiano sin prefijo",
			phone:   "573001234567",
			wantErr: false,
		},
		{
			name:    "número USA con prefijo +",
			phone:   "+13055551234",
			wantErr: false,
		},
		{
			name:    "número España con prefijo +",
			phone:   "+34612345678",
			wantErr: false,
		},
		{
			name:    "número con espacios es limpiado correctamente",
			phone:   "+57 300 123 4567",
			wantErr: false,
		},
		{
			name:    "número Colombia con guiones",
			phone:   "+57-300-123-4567",
			wantErr: false,
		},

		// Casos de error
		{
			name:        "número vacío",
			phone:       "",
			wantErr:     true,
			errContains: "vacío",
		},
		{
			name:        "número sin código de país reconocido",
			phone:       "123456789",
			wantErr:     true,
			errContains: "código de país",
		},
		{
			name:        "número colombiano demasiado corto",
			phone:       "+5730012",
			wantErr:     true,
			errContains: "dígitos",
		},
		{
			name:        "número colombiano demasiado largo",
			phone:       "+5730012345678901",
			wantErr:     true,
			errContains: "dígitos",
		},
		{
			name:        "prefijo + sin dígitos adicionales",
			phone:       "+",
			wantErr:     true,
			errContains: "código de país",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneNumber(tt.phone)

			if tt.wantErr && err == nil {
				t.Errorf("ValidatePhoneNumber(%q) esperaba error, no obtuvo ninguno", tt.phone)
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ValidatePhoneNumber(%q) no esperaba error, obtuvo: %v", tt.phone, err)
				return
			}
			if tt.wantErr && tt.errContains != "" && err != nil {
				if !containsSubstring(err.Error(), tt.errContains) {
					t.Errorf("ValidatePhoneNumber(%q) error = %q, se esperaba que contuviera %q",
						tt.phone, err.Error(), tt.errContains)
				}
			}
		})
	}
}

func TestExtractCountryCodeAndNumber(t *testing.T) {
	tests := []struct {
		name            string
		phone           string
		wantCountryCode string
		wantNumber      string
	}{
		{
			name:            "Colombia (57)",
			phone:           "573001234567",
			wantCountryCode: "57",
			wantNumber:      "3001234567",
		},
		{
			name:            "USA (1)",
			phone:           "13055551234",
			wantCountryCode: "1",
			wantNumber:      "3055551234",
		},
		{
			name:            "España (34)",
			phone:           "34612345678",
			wantCountryCode: "34",
			wantNumber:      "612345678",
		},
		{
			name:            "código de país desconocido",
			phone:           "999123456",
			wantCountryCode: "",
			wantNumber:      "999123456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotNumber := extractCountryCodeAndNumber(tt.phone)

			if gotCode != tt.wantCountryCode {
				t.Errorf("extractCountryCodeAndNumber(%q) countryCode = %q, quería %q",
					tt.phone, gotCode, tt.wantCountryCode)
			}
			if gotNumber != tt.wantNumber {
				t.Errorf("extractCountryCodeAndNumber(%q) number = %q, quería %q",
					tt.phone, gotNumber, tt.wantNumber)
			}
		})
	}
}

func TestGetCountryCodeLength(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		wantMin     int
		wantMax     int
	}{
		{name: "Colombia 57", countryCode: "57", wantMin: 10, wantMax: 10},
		{name: "USA 1", countryCode: "1", wantMin: 10, wantMax: 10},
		{name: "España 34", countryCode: "34", wantMin: 9, wantMax: 9},
		{name: "Brasil 55", countryCode: "55", wantMin: 10, wantMax: 11},
		{name: "código desconocido usa defaults", countryCode: "999", wantMin: 7, wantMax: 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := getCountryCodeLength(tt.countryCode)

			if gotMin != tt.wantMin {
				t.Errorf("getCountryCodeLength(%q) min = %d, quería %d", tt.countryCode, gotMin, tt.wantMin)
			}
			if gotMax != tt.wantMax {
				t.Errorf("getCountryCodeLength(%q) max = %d, quería %d", tt.countryCode, gotMax, tt.wantMax)
			}
		})
	}
}

// containsSubstring es un helper para verificar subcadenas en mensajes de error
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || searchSubstring(s, substr))
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

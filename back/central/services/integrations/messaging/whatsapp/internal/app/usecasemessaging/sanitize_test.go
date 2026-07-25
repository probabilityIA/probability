package usecasemessaging

import (
	"strings"
	"testing"
)

func TestSanitizeTemplateParam(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"vacio", "", ""},
		{"sin cambios", "Casco AGV x2", "Casco AGV x2"},
		{"salto de linea", "- Casco AGV\n  Cant: 2", "- Casco AGV Cant: 2"},
		{"tabs", "Casco\tAGV", "Casco AGV"},
		{"crlf", "Casco\r\nAGV", "Casco AGV"},
		{"espacios multiples", "Casco      AGV", "Casco AGV"},
		{"bordes", "  Casco AGV \n", "Casco AGV"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SanitizeTemplateParam(tc.input)
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
			if strings.ContainsAny(got, "\n\r\t") {
				t.Fatalf("resultado conserva caracteres prohibidos por Meta: %q", got)
			}
			if strings.Contains(got, "    ") {
				t.Fatalf("resultado conserva 4+ espacios consecutivos: %q", got)
			}
		})
	}
}

func TestSanitizeTemplateVariablesNil(t *testing.T) {
	if SanitizeTemplateVariables(nil) != nil {
		t.Fatal("nil debe seguir siendo nil")
	}
}

func TestSanitizeTemplateVariables(t *testing.T) {
	in := map[string]string{"1": "Juan", "2": "- Casco\n  Cant: 2"}
	got := SanitizeTemplateVariables(in)

	if got["2"] != "- Casco Cant: 2" {
		t.Fatalf("got %q", got["2"])
	}
	if in["2"] != "- Casco\n  Cant: 2" {
		t.Fatal("no debe mutar el mapa original")
	}
}

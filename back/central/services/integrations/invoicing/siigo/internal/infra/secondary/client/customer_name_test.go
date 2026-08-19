package client

import (
	"strings"
	"testing"
)

func TestPersonSiempreMandaExactamenteDosElementos(t *testing.T) {
	casos := []string{
		"Naldy Geraldine Alvarez Perez",
		"Juan",
		"Juan Perez",
		"Maria del Carmen Gomez",
		"",
		"   ",
		"Ana Maria de los Santos Rodriguez Vargas",
	}

	for _, caso := range casos {
		got := buildSiigoCustomerName(caso, "person")
		if len(got) != 2 {
			t.Errorf("nombre %q -> %d elementos, Siigo exige 2 para person", caso, len(got))
		}
		for i, parte := range got {
			if strings.TrimSpace(parte) == "" {
				t.Errorf("nombre %q -> elemento %d vacio: Siigo responde invalid_array", caso, i)
			}
		}
	}
}

func TestCompanySiempreMandaUnSoloElemento(t *testing.T) {
	casos := []string{"Distribuciones El Sol SAS", "Acme", ""}

	for _, caso := range casos {
		got := buildSiigoCustomerName(caso, "company")
		if len(got) != 1 {
			t.Errorf("empresa %q -> %d elementos, Siigo exige 1 para company", caso, len(got))
		}
		if strings.TrimSpace(got[0]) == "" {
			t.Errorf("empresa %q -> elemento vacio", caso)
		}
	}
}

func TestSeQuitanAcentosYCaracteresEspeciales(t *testing.T) {
	got := buildSiigoCustomerName("Naldy Geraldiná Álvarez Pérez", "person")
	juntos := strings.Join(got, " ")

	for _, r := range juntos {
		if r > 127 {
			t.Fatalf("quedo un caracter no ASCII (%q) en %q: Siigo no permite caracteres especiales", r, juntos)
		}
	}
	if !strings.Contains(juntos, "Alvarez") {
		t.Errorf("se esperaba Alvarez sin tilde, quedo %q", juntos)
	}
}

func TestSeLimpianSimbolosSinPerderLasPalabras(t *testing.T) {
	got := buildSiigoCustomerName("Juan-Carlos O'Brien & Cia.", "person")
	juntos := strings.Join(got, " ")

	if strings.ContainsAny(juntos, "-'&.") {
		t.Errorf("quedaron simbolos en %q", juntos)
	}
	if !strings.Contains(juntos, "Juan") || !strings.Contains(juntos, "Carlos") {
		t.Errorf("se perdieron palabras al limpiar: %q", juntos)
	}
}

func TestCadaCampoRespetaElMaximoDeCienCaracteres(t *testing.T) {
	largo := strings.Repeat("Alejandro ", 30)

	for _, tipo := range []string{"person", "company"} {
		for _, parte := range buildSiigoCustomerName(largo, tipo) {
			if len(parte) > 100 {
				t.Errorf("%s: campo de %d caracteres, el maximo de Siigo es 100", tipo, len(parte))
			}
		}
	}
}

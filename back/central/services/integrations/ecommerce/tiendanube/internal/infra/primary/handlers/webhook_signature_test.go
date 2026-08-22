package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func firmar(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func TestFirmaValida(t *testing.T) {
	body := []byte(`{"store_id":8126740,"event":"order/created","id":2051703380}`)
	secret := "un-client-secret-de-prueba"

	casos := []struct {
		nombre   string
		body     []byte
		secret   string
		firma    string
		esperado bool
	}{
		{"firma correcta", body, secret, firmar(body, secret), true},
		{"firma en mayusculas", body, secret, "ABCDEF", false},
		{"secret distinto", body, secret, firmar(body, "otro-secret"), false},
		{"body alterado", []byte(`{"store_id":8126740,"event":"order/created","id":999}`), secret, firmar(body, secret), false},
		{"firma vacia", body, secret, "", false},
		{"secret vacio", body, "", firmar(body, secret), false},
		{"firma no hexadecimal", body, secret, "no-es-hex", false},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			if got := firmaValida(caso.body, caso.secret, caso.firma); got != caso.esperado {
				t.Fatalf("firmaValida = %v, se esperaba %v", got, caso.esperado)
			}
		})
	}
}

func TestFirmaValidaAceptaEspaciosEnElHeader(t *testing.T) {
	body := []byte(`{"event":"order/paid"}`)
	secret := "secreto"
	if !firmaValida(body, secret, "  "+firmar(body, secret)+"  ") {
		t.Fatal("la firma con espacios alrededor debia aceptarse")
	}
}

func TestStoreIDDeIntegracion(t *testing.T) {
	casos := []struct {
		nombre   string
		storeID  string
		config   map[string]interface{}
		esperado string
	}{
		{"columna store_id", "8126740", nil, "8126740"},
		{"desde config", "", map[string]interface{}{"store_id": "8126740"}, "8126740"},
		{"config numerica", "", map[string]interface{}{"store_id": float64(8126740)}, "8.12674e+06"},
		{"sin dato", "", nil, ""},
		{"config sin la llave", "", map[string]interface{}{"otra": "cosa"}, ""},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			if got := storeIDDeIntegracion(caso.storeID, caso.config); got != caso.esperado {
				t.Fatalf("storeIDDeIntegracion = %q, se esperaba %q", got, caso.esperado)
			}
		})
	}
}

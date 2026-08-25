package client

import (
	"strings"
	"testing"
)

func TestMapEnvioClickError_NoDevuelveElNombreDelProveedor(t *testing.T) {
	entradas := []string{
		"EnvioClick devolvio 500",
		"dial tcp api.envioclickpro.com.co:443 connect: connection refused",
		"envio click no proceso la solicitud",
		"https://api.envioclickpro.com.co/api/v2 timeout",
	}

	for _, in := range entradas {
		got := mapEnvioClickError(in)
		if strings.Contains(strings.ToLower(got), "envio") && strings.Contains(strings.ToLower(got), "click") {
			t.Fatalf("entrada %q expone el proveedor: %q", in, got)
		}
		if strings.Contains(strings.ToLower(got), "envioclick") || strings.Contains(strings.ToLower(got), "envioclik") {
			t.Fatalf("entrada %q expone el proveedor: %q", in, got)
		}
	}
}

func TestMapEnvioClickError_ConservaElDetalleUtil(t *testing.T) {
	got := mapEnvioClickError("order 8891 not found")

	if !strings.Contains(got, "order 8891 not found") {
		t.Fatalf("esperaba conservar el detalle, obtuve: %q", got)
	}
}

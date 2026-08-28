package client

import (
	"strings"
	"testing"
)

func TestSiigoErrorMessage(t *testing.T) {
	casos := []struct {
		nombre   string
		body     string
		contiene string
	}{
		{
			nombre:   "documento inactivo",
			body:     `{"Status":400,"Errors":[{"Code":"parameter_inactive","Message":"The id is inactive: 30606","Params":["document.id"],"Detail":"x"}]}`,
			contiene: "tipo de documento configurado (id 30606) esta inactivo en Siigo",
		},
		{
			nombre:   "producto inactivo",
			body:     `{"Status":400,"Errors":[{"Code":"parameter_inactive","Message":"The code is inactive: BN116-2XL","Params":["items[1].code"]}]}`,
			contiene: "producto con codigo BN116-2XL (item 1 de la factura) esta inactivo",
		},
		{
			nombre:   "documento sin electronica",
			body:     `{"Status":400,"Errors":[{"Code":"document_settings","Message":"The send cannot be used, you must verify the document settings","Params":["stamp.send"]}]}`,
			contiene: "no esta habilitado para facturacion electronica",
		},
		{
			nombre:   "servicio caido",
			body:     `{"Status":400,"Errors":[{"Code":"documents_service","Message":"The Documents service is currently unavailable. Please try in a few minutes"}]}`,
			contiene: "no esta disponible en este momento",
		},
		{
			nombre:   "vendedor inexistente",
			body:     `{"Status":400,"Errors":[{"Code":"invalid_reference","Message":"The seller doesn't exist: 0","Params":["seller"]}]}`,
			contiene: "no existe en Siigo (seller)",
		},
		{
			nombre:   "codigo desconocido",
			body:     `{"Status":400,"Errors":[{"Code":"otro","Message":"algo paso","Params":["campo.x"]}]}`,
			contiene: "algo paso (campo.x)",
		},
		{
			nombre:   "cuerpo no json",
			body:     `<html>502</html>`,
			contiene: "codigo 400",
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got := siigoErrorMessage([]byte(c.body), 400, "la factura")
			if !strings.Contains(got, c.contiene) {
				t.Fatalf("mensaje = %q, se esperaba que contuviera %q", got, c.contiene)
			}
		})
	}
}

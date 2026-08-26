package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func TestElCuerpoRealDeProduccionSeReconoce(t *testing.T) {
	real := []byte(`{"failed_shipments":[{"shipment_id":"47846252129","order_id":2000018097840680,"message":"Shipment 47846252129 status is shipped","error_code":"shipment_not_printable"}]}`)

	if err := classifyLabelBody(real); !errors.Is(err, domain.ErrLabelAlreadyShipped) {
		t.Fatalf("esperaba ErrLabelAlreadyShipped, obtuve %v", err)
	}
}

func TestOtrosFallosDelCanalNoSeConfundenConDespachado(t *testing.T) {
	otro := []byte(`{"failed_shipments":[{"shipment_id":"1","message":"Shipment 1 does not exist","error_code":"not_found"}]}`)

	if err := classifyLabelBody(otro); !errors.Is(err, domain.ErrLabelNotAvailable) {
		t.Fatalf("esperaba ErrLabelNotAvailable, obtuve %v", err)
	}
}

func TestUnPdfNoSeClasificaComoError(t *testing.T) {
	if err := classifyLabelBody([]byte("%PDF-1.4 contenido binario")); err != nil {
		t.Fatalf("un pdf no debia clasificarse como error: %v", err)
	}
}

func TestSoloUnPdfRealPasaComoPdf(t *testing.T) {
	casos := map[string]bool{
		"%PDF-1.4\n1 0 obj":               true,
		"\n  %PDF-1.7":                    true,
		`{"failed_shipments":[]}`:         false,
		"<html><body>error</body></html>": false,
		"":                                false,
	}

	for cuerpo, esperado := range casos {
		if got := looksLikePDF([]byte(cuerpo)); got != esperado {
			t.Fatalf("looksLikePDF(%q) = %v, esperado %v", cuerpo, got, esperado)
		}
	}
}

func TestSeDetectaJsonDisfrazadoDePdf(t *testing.T) {
	if !looksLikeJSON([]byte(`{"failed_shipments":[]}`)) {
		t.Fatal("un objeto json debia detectarse")
	}
	if !looksLikeJSON([]byte(`  [1,2]`)) {
		t.Fatal("un arreglo json debia detectarse")
	}
	if looksLikeJSON([]byte("%PDF-1.4")) {
		t.Fatal("un pdf no es json")
	}
}

func clienteContra(t *testing.T, status int, contentType string, body string) domain.IMeliClient {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return New().WithBaseURL(srv.URL)
}

func TestGetShipmentLabelRechazaJsonDisfrazadoDePdf(t *testing.T) {
	cli := clienteContra(t, http.StatusOK, "application/json", `{"failed_shipments":[{"shipment_id":"1","message":"Shipment 1 status is shipped"}]}`)

	label, err := cli.GetShipmentLabel(context.Background(), "token", 1, "pdf")

	if label != nil {
		t.Fatalf("no debia devolver un archivo: %d bytes con content-type %q", len(label.Content), label.ContentType)
	}
	if !errors.Is(err, domain.ErrLabelAlreadyShipped) {
		t.Fatalf("esperaba ErrLabelAlreadyShipped, obtuve %v", err)
	}
}

func TestGetShipmentLabelRechazaUnCuerpoQueNoEsPdf(t *testing.T) {
	cli := clienteContra(t, http.StatusOK, "application/pdf", "<html><body>Service Unavailable</body></html>")

	label, err := cli.GetShipmentLabel(context.Background(), "token", 1, "pdf")

	if label != nil {
		t.Fatal("un html no debia entregarse al navegador como pdf")
	}
	if !errors.Is(err, domain.ErrLabelNotPDF) {
		t.Fatalf("esperaba ErrLabelNotPDF, obtuve %v", err)
	}
}

func TestGetShipmentLabelDevuelveElPdfCuandoEsValido(t *testing.T) {
	cli := clienteContra(t, http.StatusOK, "application/pdf", "%PDF-1.4\ncontenido")

	label, err := cli.GetShipmentLabel(context.Background(), "token", 47846252129, "pdf")

	if err != nil {
		t.Fatalf("un pdf valido no debia fallar: %v", err)
	}
	if label.ContentType != "application/pdf" {
		t.Fatalf("content-type inesperado: %q", label.ContentType)
	}
	if label.Filename != "meli-47846252129.pdf" {
		t.Fatalf("nombre de archivo inesperado: %q", label.Filename)
	}
}

func TestGetShipmentLabelTraduceEl400DeEnvioDespachado(t *testing.T) {
	cli := clienteContra(t, http.StatusBadRequest, "application/json",
		`{"failed_shipments":[{"shipment_id":"47846252129","message":"Shipment 47846252129 status is shipped","error_code":"shipment_not_printable"}]}`)

	_, err := cli.GetShipmentLabel(context.Background(), "token", 47846252129, "pdf")

	if !errors.Is(err, domain.ErrLabelAlreadyShipped) {
		t.Fatalf("esperaba ErrLabelAlreadyShipped, obtuve %v", err)
	}
	if strings.Contains(err.Error(), "47846252129") {
		t.Fatalf("el error no debe arrastrar ids internos del canal: %q", err)
	}
}

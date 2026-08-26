package client

import (
	"errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

func requestValido() domain.QuoteRequest {
	return domain.QuoteRequest{
		Description: "Compra en linea",
		Origin: domain.Address{
			Address: "Carrera 10 # 20-30", Suburb: "Centro", FirstName: "Bodega",
			Phone: "3001234567", Email: "bodega@negocio.com",
		},
		Destination: domain.Address{
			Address: "Calle 5 # 6-7", Suburb: "Poblado", FirstName: "Cliente",
			Phone: "3009876543", Email: "cliente@correo.com",
		},
	}
}

func TestDireccionesRealesQueLaTransportadoraRechazaba(t *testing.T) {
	reales := []string{
		"Calle 114#42-300 Torre 5 apartamento 902 alameda del rio conjunto tucan",
		"Calle 114 42 300 Torre 5 apartamento 902 alameda del rio conjunto tucan",
		"CRA 27 f #108-46 las orquideas Casa primer piso la frente del parque",
		"CRA 27 f #108-46 Casa primer piso la frente del parque",
	}

	for _, direccion := range reales {
		t.Run(direccion[:20], func(t *testing.T) {
			req := requestValido()
			req.Destination.Address = direccion

			if err := normalizeQuoteRequest(&req); err != nil {
				t.Fatalf("no debia rechazarse: %v", err)
			}

			largo := len([]rune(req.Destination.Address))
			if largo > maxAddressLen {
				t.Fatalf("quedo en %d caracteres, el limite es %d: %q", largo, maxAddressLen, req.Destination.Address)
			}
			if largo < minTextLen {
				t.Fatalf("quedo demasiado corta: %q", req.Destination.Address)
			}
			if !strings.HasPrefix(direccion, req.Destination.Address) {
				t.Fatalf("se conserva el inicio de la direccion, que es lo que ubica al mensajero: %q", req.Destination.Address)
			}
		})
	}
}

func TestNoCortaAMitadDePalabra(t *testing.T) {
	req := requestValido()
	req.Destination.Address = "Calle 114#42-300 Torre 5 apartamento 902 alameda del rio"

	if err := normalizeQuoteRequest(&req); err != nil {
		t.Fatalf("no debia rechazarse: %v", err)
	}

	if strings.HasSuffix(req.Destination.Address, " ") {
		t.Fatalf("quedo un espacio al final: %q", req.Destination.Address)
	}
	if req.Destination.Address != strings.TrimSpace(req.Destination.Address) {
		t.Fatalf("quedaron espacios sueltos: %q", req.Destination.Address)
	}
}

func TestBarrioYDescripcionDentroDeLimite(t *testing.T) {
	req := requestValido()
	req.Destination.Suburb = "Un barrio con un nombre larguisimo que no cabe de ninguna manera"
	req.Description = "Una descripcion demasiado larga para el limite"

	if err := normalizeQuoteRequest(&req); err != nil {
		t.Fatalf("no debia rechazarse: %v", err)
	}

	if len([]rune(req.Destination.Suburb)) > maxSuburbLen {
		t.Fatalf("barrio de %d caracteres, limite %d", len([]rune(req.Destination.Suburb)), maxSuburbLen)
	}
	if len([]rune(req.Description)) > maxDescriptionLen {
		t.Fatalf("descripcion de %d caracteres, limite %d", len([]rune(req.Description)), maxDescriptionLen)
	}
}

func TestBarrioVacioSeCompletaConLaDireccion(t *testing.T) {
	req := requestValido()
	req.Destination.Suburb = ""

	if err := normalizeQuoteRequest(&req); err != nil {
		t.Fatalf("no debia rechazarse: %v", err)
	}

	if len([]rune(req.Destination.Suburb)) < minTextLen {
		t.Fatalf("el barrio vacio tambien rompia la cotizacion, debia completarse: %q", req.Destination.Suburb)
	}
}

func TestTelefonoSeLimpiaYRecorta(t *testing.T) {
	casos := map[string]string{
		"+57 300 123 4567": "3001234567",
		"(604) 444-5566":   "6044445566",
		"300-123-4567":     "3001234567",
		"57 3001234567":    "3001234567",
	}

	for entrada, esperado := range casos {
		req := requestValido()
		req.Destination.Phone = entrada

		if err := normalizeQuoteRequest(&req); err != nil {
			t.Fatalf("no debia rechazarse: %v", err)
		}
		if req.Destination.Phone != esperado {
			t.Fatalf("telefono %q quedo en %q, esperado %q", entrada, req.Destination.Phone, esperado)
		}
	}
}

func TestDireccionVaciaSeRechazaAntesDeLlamar(t *testing.T) {
	req := requestValido()
	req.Destination.Address = "  "

	err := normalizeQuoteRequest(&req)
	if !errors.Is(err, ErrDestinationAddressMissing) {
		t.Fatalf("esperaba ErrDestinationAddressMissing, obtuve %v", err)
	}

	req = requestValido()
	req.Origin.Address = ""
	if err := normalizeQuoteRequest(&req); !errors.Is(err, ErrOriginAddressMissing) {
		t.Fatalf("esperaba ErrOriginAddressMissing, obtuve %v", err)
	}
}

func TestElErrorNoNombraAlProveedor(t *testing.T) {
	for _, err := range []error{ErrDestinationAddressMissing, ErrOriginAddressMissing} {
		if strings.Contains(strings.ToLower(err.Error()), "envioclick") {
			t.Fatalf("el error nombra al proveedor: %q", err)
		}
	}
}

func TestUnRequestValidoNoSeToca(t *testing.T) {
	req := requestValido()
	original := req

	if err := normalizeQuoteRequest(&req); err != nil {
		t.Fatalf("no debia rechazarse: %v", err)
	}

	if req.Destination.Address != original.Destination.Address ||
		req.Destination.Suburb != original.Destination.Suburb ||
		req.Origin.Address != original.Origin.Address ||
		req.Description != original.Description {
		t.Fatalf("un request que ya cumplia los limites no debia modificarse: %+v", req)
	}
}

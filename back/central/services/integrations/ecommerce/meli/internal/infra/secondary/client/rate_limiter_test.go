package client

import (
	"context"
	"testing"
	"time"
)

func TestCadaCuentaTieneSuPropiaCuota(t *testing.T) {
	base := New().(*MeliClient)

	a := base.ForAccount("254").(*MeliClient)
	b := base.ForAccount("221").(*MeliClient)

	if a.limiter == b.limiter {
		t.Fatal("el limite de MercadoLibre es por vendedor: dos cuentas no pueden compartir cuota")
	}
}

func TestLaMismaCuentaReusaSuCuota(t *testing.T) {
	base := New().(*MeliClient)

	a := base.ForAccount("254").(*MeliClient)
	b := base.ForAccount("254").(*MeliClient)

	if a.limiter != b.limiter {
		t.Fatal("la misma cuenta debe compartir su propio limitador, no crear uno nuevo por llamada")
	}
}

func TestSinCuentaSeQuedaConElLimitadorBase(t *testing.T) {
	base := New().(*MeliClient)

	if got := base.ForAccount(""); got != domainClient(base) {
		t.Fatal("sin clave no hay a que cuenta cobrarle, se usa el mismo cliente")
	}
}

func domainClient(c *MeliClient) interface{} { return c }

func TestElModoPruebasConservaLaCuotaDeLaCuenta(t *testing.T) {
	base := New().(*MeliClient)

	cuenta := base.ForAccount("254").(*MeliClient)
	pruebas := cuenta.WithBaseURL("https://api.pruebas").(*MeliClient)

	if pruebas.limiter != cuenta.limiter {
		t.Fatal("cambiar de base_url no debe reiniciar la cuota de la cuenta")
	}
	if pruebas.limiters != cuenta.limiters {
		t.Fatal("el registro de limitadores debe seguir siendo el mismo")
	}
}

func TestElLimitadorDejaPasarDentroDeLaCuota(t *testing.T) {
	r := newRateLimiter(600)
	ctx := context.Background()

	inicio := time.Now()
	for i := 0; i < 50; i++ {
		if err := r.Wait(ctx); err != nil {
			t.Fatalf("no debia bloquear dentro de la cuota: %v", err)
		}
	}
	if time.Since(inicio) > 200*time.Millisecond {
		t.Fatal("50 llamadas con cuota de 600/min no deben esperar")
	}
}

func TestElLimitadorFrenaAlAgotarLaCuota(t *testing.T) {
	r := newRateLimiter(60)
	ctx := context.Background()

	for i := 0; i < 60; i++ {
		if err := r.Wait(ctx); err != nil {
			t.Fatal(err)
		}
	}

	ctxCorto, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	if err := r.Wait(ctxCorto); err == nil {
		t.Fatal("agotada la cuota debe esperar, y respetar la cancelacion del contexto")
	}
}

func TestElLimitadorPorDefectoUsaLaCuotaDocumentada(t *testing.T) {
	r := newRateLimiter(0)
	if r.capacity != defaultRequestsPerMinute {
		t.Fatalf("un valor invalido debe caer al por defecto, obtuve %v", r.capacity)
	}
	if defaultRequestsPerMinute > 1500 {
		t.Fatalf("MercadoLibre permite 1500 por minuto por vendedor, no podemos configurar %v", defaultRequestsPerMinute)
	}
}

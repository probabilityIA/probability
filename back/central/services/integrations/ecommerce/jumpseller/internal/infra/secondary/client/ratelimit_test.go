package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

func TestIntervalForLimit(t *testing.T) {
	if got := intervalForLimit(60); got != time.Second {
		t.Fatalf("60 por minuto = %v, se esperaba 1s", got)
	}
	if got := intervalForLimit(120); got != 500*time.Millisecond {
		t.Fatalf("120 por minuto = %v, se esperaba 500ms", got)
	}
	if got := intervalForLimit(0); got != intervalForLimit(assumedRateLimitPerMinute) {
		t.Fatal("limite invalido debe caer al ritmo conservador")
	}
	if got := intervalForLimit(100000); got != minRequestInterval {
		t.Fatalf("un limite enorme no debe bajar de %v, dio %v", minRequestInterval, got)
	}
}

func TestPacerSeAutoAjustaConElHeader(t *testing.T) {
	p := newStorePacer()
	inicial := p.interval

	p.observeLimit("120")
	if p.interval != 500*time.Millisecond {
		t.Fatalf("interval = %v, se esperaba 500ms tras leer el header", p.interval)
	}
	if p.interval == inicial {
		t.Fatal("el pacer debe ajustarse al limite real de la tienda")
	}

	p.observeLimit("")
	if p.interval != 500*time.Millisecond {
		t.Fatal("un header ausente no debe cambiar el ritmo")
	}

	p.observeLimit("basura")
	if p.interval != 500*time.Millisecond {
		t.Fatal("un header invalido no debe cambiar el ritmo")
	}
}

func TestPacerBackOffTieneTecho(t *testing.T) {
	p := newStorePacer()
	for i := 0; i < 20; i++ {
		p.backOff()
	}
	if p.interval != maxRequestInterval {
		t.Fatalf("interval = %v, el backoff debe topar en %v", p.interval, maxRequestInterval)
	}
}

func TestPacerEspaciaLasPeticiones(t *testing.T) {
	p := newStorePacer()
	p.observeLimit("600")

	inicio := time.Now()
	for i := 0; i < 3; i++ {
		if err := p.wait(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	transcurrido := time.Since(inicio)

	if transcurrido < 150*time.Millisecond {
		t.Fatalf("3 peticiones a 600/min tardaron %v: no se esta espaciando", transcurrido)
	}
}

func TestPacerRespetaCancelacion(t *testing.T) {
	p := newStorePacer()
	_ = p.wait(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := p.wait(ctx); err == nil {
		t.Fatal("con el contexto cancelado wait debe devolver error, no dormir")
	}
}

func TestReintentaTrasRateLimitYConservaElCuerpo(t *testing.T) {
	var llamadas int32
	var cuerpos []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		cuerpos = append(cuerpos, string(body))

		w.Header().Set(domain.RateLimitHeader, "600")
		w.Header().Set("Content-Type", "application/json")

		if atomic.AddInt32(&llamadas, 1) == 1 {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, _ = w.Write([]byte(`{"product":{"id":100}}`))
	}))
	defer server.Close()

	client := New()
	cred := domain.Credential{APIKey: "login", APISecret: "token", BaseURL: server.URL}

	err := client.SetProductStock(context.Background(), cred, 100, 7)
	if err != nil {
		t.Fatalf("tras un 403 de rate limit deberia reintentar y salir bien: %v", err)
	}
	if llamadas != 2 {
		t.Fatalf("llamadas = %d, se esperaban 2 (una fallida + un reintento)", llamadas)
	}
	if len(cuerpos) != 2 || cuerpos[0] == "" || cuerpos[0] != cuerpos[1] {
		t.Fatalf("el reintento debe reenviar el MISMO cuerpo, no uno vacio: %q vs %q", cuerpos[0], cuerpos[1])
	}
}

func TestForbiddenSinHeaderEsCredencialInvalida(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := New()
	_, err := client.GetStoreInfo(context.Background(), domain.Credential{APIKey: "x", APISecret: "y", BaseURL: server.URL})
	if err != domain.ErrInvalidCredentials {
		t.Fatalf("err = %v: un 403 sin header de rate limit es credencial invalida, no rate limit", err)
	}
}

func TestSeAgotanLosReintentos(t *testing.T) {
	var llamadas int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&llamadas, 1)
		w.Header().Set(domain.RateLimitHeader, "600")
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := New()
	_, err := client.GetStoreInfo(context.Background(), domain.Credential{APIKey: "z", APISecret: "y", BaseURL: server.URL})
	if err == nil {
		t.Fatal("si el rate limit nunca cede debe fallar, no colgarse")
	}
	if llamadas != maxRateLimitRetries+1 {
		t.Fatalf("llamadas = %d, se esperaban %d intentos", llamadas, maxRateLimitRetries+1)
	}
}

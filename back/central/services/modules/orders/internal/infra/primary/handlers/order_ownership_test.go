package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type ownershipCase struct {
	name    string
	method  string
	body    string
	handler func(h *Handlers, c *gin.Context)
}

func ownershipCases() []ownershipCase {
	return []ownershipCase{
		{"raw", http.MethodGet, "", func(h *Handlers, c *gin.Context) { h.GetOrderRaw(c) }},
		{"history", http.MethodGet, "", func(h *Handlers, c *gin.Context) { h.GetOrderHistory(c) }},
		{"update", http.MethodPut, `{}`, func(h *Handlers, c *gin.Context) { h.UpdateOrder(c) }},
		{"delete", http.MethodDelete, "", func(h *Handlers, c *gin.Context) { h.DeleteOrder(c) }},
		{"change-status", http.MethodPut, `{"status":"confirmed"}`, func(h *Handlers, c *gin.Context) { h.ChangeStatus(c) }},
	}
}

func ejecutarConNegocio(t *testing.T, tc ownershipCase, stub *orderUseCaseStub, businessIDToken uint) int {
	t.Helper()
	gin.SetMode(gin.TestMode)

	h := &Handlers{orderCRUD: stub}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(tc.method, "/orders/abc", strings.NewReader(tc.body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	c.Set("business_id", businessIDToken)

	tc.handler(h, c)
	return w.Code
}

func TestOrdenDeOtroNegocio_RechazaEnTodosLosEndpointsPorID(t *testing.T) {
	for _, tc := range ownershipCases() {
		t.Run(tc.name, func(t *testing.T) {
			code := ejecutarConNegocio(t, tc, &orderUseCaseStub{orden: ordenDe(46)}, 26)
			if code != http.StatusForbidden {
				t.Fatalf("se esperaba 403, llego %d", code)
			}
		})
	}
}

func TestOrdenInexistente_DevuelveNoEncontrada(t *testing.T) {
	for _, tc := range ownershipCases() {
		t.Run(tc.name, func(t *testing.T) {
			stub := &orderUseCaseStub{err: errOrdenNoEncontrada{}}
			code := ejecutarConNegocio(t, tc, stub, 26)
			if code != http.StatusNotFound {
				t.Fatalf("se esperaba 404, llego %d", code)
			}
		})
	}
}

func TestSuperAdmin_NoSeValidaPropiedad(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &Handlers{orderCRUD: &orderUseCaseStub{orden: ordenDe(46)}}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/orders/abc/history", nil)
	c.Set("business_id", uint(0))

	if !h.ensureOrderOwnership(c, "abc") {
		t.Fatalf("el super admin debe pasar sin validar propiedad")
	}
}

func TestMismoNegocio_PasaLaValidacion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &Handlers{orderCRUD: &orderUseCaseStub{orden: ordenDe(26)}}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/orders/abc/history", nil)
	c.Set("business_id", uint(26))

	if !h.ensureOrderOwnership(c, "abc") {
		t.Fatalf("la orden del mismo negocio debe pasar, respuesta %d", w.Code)
	}
}

type errOrdenNoEncontrada struct{}

func (errOrdenNoEncontrada) Error() string { return "order not found" }

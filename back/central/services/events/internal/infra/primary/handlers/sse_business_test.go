package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func contextoSSE(businessID interface{}, query string, param string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/notify/sse/order-notify"+query, nil)
	if param != "" {
		c.Params = gin.Params{{Key: "businessID", Value: param}}
	}
	if businessID != nil {
		c.Set("business_id", businessID)
	}
	return c, w
}

func TestSSE_UsuarioNormalNoPuedeEscucharOtroNegocio(t *testing.T) {
	h := &SSEHandler{}
	c, _ := contextoSSE(uint(26), "?business_id=46", "")

	got, ok := h.resolveSSEBusinessID(c)
	if !ok || got != 26 {
		t.Fatalf("se esperaba 26, llego %d ok=%v", got, ok)
	}
}

func TestSSE_UsuarioNormalConPathDeOtroNegocio(t *testing.T) {
	h := &SSEHandler{}
	c, _ := contextoSSE(uint(26), "", "46")

	got, ok := h.resolveSSEBusinessID(c)
	if !ok || got != 26 {
		t.Fatalf("se esperaba 26, llego %d ok=%v", got, ok)
	}
}

func TestSSE_SuperAdminEligeNegocio(t *testing.T) {
	h := &SSEHandler{}
	c, _ := contextoSSE(uint(0), "?business_id=46", "")

	got, ok := h.resolveSSEBusinessID(c)
	if !ok || got != 46 {
		t.Fatalf("se esperaba 46, llego %d ok=%v", got, ok)
	}
}

func TestSSE_SinSesionEs401(t *testing.T) {
	h := &SSEHandler{}
	c, w := contextoSSE(nil, "?business_id=46", "")

	if _, ok := h.resolveSSEBusinessID(c); ok {
		t.Fatalf("no debia resolver")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, llego %d", w.Code)
	}
}

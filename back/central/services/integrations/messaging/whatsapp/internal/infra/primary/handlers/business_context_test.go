package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func contextoConNegocio(businessID interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	if businessID != nil {
		c.Set("business_id", businessID)
	}
	return c, w
}

func TestResolveBusinessID_UsuarioNormalIgnoraElDelBody(t *testing.T) {
	c, _ := contextoConNegocio(uint(26))

	got, ok := resolveBusinessID(c, 46)
	if !ok || got != 26 {
		t.Fatalf("se esperaba el negocio del token (26), llego %d ok=%v", got, ok)
	}
}

func TestResolveBusinessID_SuperAdminUsaElDelBody(t *testing.T) {
	c, _ := contextoConNegocio(uint(0))

	got, ok := resolveBusinessID(c, 46)
	if !ok || got != 46 {
		t.Fatalf("se esperaba 46, llego %d ok=%v", got, ok)
	}
}

func TestResolveBusinessID_SuperAdminSinNegocioEs400(t *testing.T) {
	c, w := contextoConNegocio(uint(0))

	if _, ok := resolveBusinessID(c, 0); ok {
		t.Fatalf("no debia resolver")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, llego %d", w.Code)
	}
}

func TestResolveBusinessID_SinContextoEs401(t *testing.T) {
	c, w := contextoConNegocio(nil)

	if _, ok := resolveBusinessID(c, 26); ok {
		t.Fatalf("no debia resolver")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, llego %d", w.Code)
	}
}

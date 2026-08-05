package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func contextoDePrueba() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

func TestGetAuthInfo_SinDatos(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetAuthInfo(c)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestGetAuthInfo_TipoIncorrecto_NoRompe(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("auth_info", "esto no es un AuthInfo")

	got, ok := GetAuthInfo(c)

	assert.False(t, ok, "un valor de otro tipo no se interpreta como autenticacion valida")
	assert.Nil(t, got)
}

func TestGetAuthInfo_Presente(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("auth_info", &domain.AuthInfo{UserID: 7, Email: "ana@test.com"})

	got, ok := GetAuthInfo(c)

	require.True(t, ok)
	assert.Equal(t, uint(7), got.UserID)
}

func TestGetAuthType(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetAuthType(c)
	assert.False(t, ok)
	assert.Equal(t, domain.AuthTypeUnknown, got)

	c.Set("auth_type", domain.AuthTypeJWT)
	got, ok = GetAuthType(c)
	require.True(t, ok)
	assert.Equal(t, domain.AuthTypeJWT, got)
}

func TestGetAuthType_TipoIncorrecto(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("auth_type", 42)

	got, ok := GetAuthType(c)

	assert.False(t, ok)
	assert.Equal(t, domain.AuthTypeUnknown, got)
}

func TestGetAPIKey_SoloConAutenticacionPorAPIKey(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("auth_info", &domain.AuthInfo{Type: domain.AuthTypeJWT, APIKey: "secreto"})

	got, ok := GetAPIKey(c)

	assert.False(t, ok, "una sesion JWT no expone una API key aunque el campo venga lleno")
	assert.Equal(t, "", got)
}

func TestGetAPIKey_ConAutenticacionPorAPIKey(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("auth_info", &domain.AuthInfo{Type: domain.AuthTypeAPIKey, APIKey: "pk_live_abc"})

	got, ok := GetAPIKey(c)

	require.True(t, ok)
	assert.Equal(t, "pk_live_abc", got)
}

func TestGetAPIKey_SinAutenticacion(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetAPIKey(c)

	assert.False(t, ok)
	assert.Equal(t, "", got)
}

func TestGetUserID(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetUserID(c)
	assert.False(t, ok)
	assert.Equal(t, uint(0), got)

	c.Set("user_id", uint(42))
	got, ok = GetUserID(c)
	require.True(t, ok)
	assert.Equal(t, uint(42), got)
}

func TestGetUserID_TipoIncorrecto(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("user_id", "42")

	got, ok := GetUserID(c)

	assert.False(t, ok, "un user_id que no es uint no se acepta")
	assert.Equal(t, uint(0), got)
}

func TestGetUserEmail(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetUserEmail(c)
	assert.False(t, ok)
	assert.Equal(t, "", got)

	c.Set("user_email", "ana@test.com")
	got, ok = GetUserEmail(c)
	require.True(t, ok)
	assert.Equal(t, "ana@test.com", got)
}

func TestGetUserEmail_TipoIncorrecto(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("user_email", 42)

	got, ok := GetUserEmail(c)

	assert.False(t, ok)
	assert.Equal(t, "", got)
}

func TestGetUserRoles(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetUserRoles(c)
	assert.False(t, ok)
	assert.Nil(t, got)

	c.Set("user_roles", []string{"admin", "vendedor"})
	got, ok = GetUserRoles(c)
	require.True(t, ok)
	assert.Equal(t, []string{"admin", "vendedor"}, got)
}

func TestGetUserRoles_TipoIncorrecto(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("user_roles", "admin")

	got, ok := GetUserRoles(c)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestGetJWTClaims(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetJWTClaims(c)
	assert.False(t, ok)
	assert.Nil(t, got)

	c.Set("jwt_claims", &domain.JWTClaims{UserID: 7})
	got, ok = GetJWTClaims(c)
	require.True(t, ok)
	assert.Equal(t, uint(7), got.UserID)
}

func TestGetJWTClaims_TipoIncorrecto(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("jwt_claims", "claims")

	got, ok := GetJWTClaims(c)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestGetBusinessID_SaleDelAuthInfo(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetBusinessID(c)
	assert.False(t, ok)
	assert.Equal(t, uint(0), got)

	c.Set("auth_info", &domain.AuthInfo{BusinessID: 36})
	got, ok = GetBusinessID(c)
	require.True(t, ok)
	assert.Equal(t, uint(36), got)
}

func TestGetBusinessIDFromContext(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetBusinessIDFromContext(c)
	assert.False(t, ok, "sin business en contexto no hay negocio que usar")
	assert.Equal(t, uint(0), got)

	c.Set("business_id", uint(36))
	got, ok = GetBusinessIDFromContext(c)
	require.True(t, ok)
	assert.Equal(t, uint(36), got)
}

func TestGetBusinessIDFromContext_SuperAdminEsCeroPresente(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("business_id", uint(0))

	got, ok := GetBusinessIDFromContext(c)

	require.True(t, ok, "el cero del super admin esta presente, no ausente")
	assert.Equal(t, uint(0), got)
}

func TestGetBusinessIDFromContext_TipoIncorrecto(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("business_id", 36)

	got, ok := GetBusinessIDFromContext(c)

	assert.False(t, ok, "un business_id int en vez de uint no se acepta")
	assert.Equal(t, uint(0), got)
}

func TestGetBusinessTokenClaims(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetBusinessTokenClaims(c)
	assert.False(t, ok)
	assert.Nil(t, got)

	c.Set("business_token_claims", &domain.BusinessTokenClaims{BusinessID: 36})
	got, ok = GetBusinessTokenClaims(c)
	require.True(t, ok)
	assert.Equal(t, uint(36), got.BusinessID)
}

func TestGetBusinessTypeID(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetBusinessTypeID(c)
	assert.False(t, ok)
	assert.Equal(t, uint(0), got)

	c.Set("business_type_id", uint(3))
	got, ok = GetBusinessTypeID(c)
	require.True(t, ok)
	assert.Equal(t, uint(3), got)
}

func TestGetRoleID(t *testing.T) {
	c, _ := contextoDePrueba()

	got, ok := GetRoleID(c)
	assert.False(t, ok)
	assert.Equal(t, uint(0), got)

	c.Set("role_id", uint(5))
	got, ok = GetRoleID(c)
	require.True(t, ok)
	assert.Equal(t, uint(5), got)
}

func TestIsSuperAdmin_SoloConBusinessCero(t *testing.T) {
	casos := []struct {
		nombre string
		valor  any
		set    bool
		want   bool
	}{
		{nombre: "sin business en contexto", set: false, want: false},
		{nombre: "business cero es super admin", valor: uint(0), set: true, want: true},
		{nombre: "business real no es super admin", valor: uint(36), set: true, want: false},
		{nombre: "tipo incorrecto no es super admin", valor: 0, set: true, want: false},
		{nombre: "string cero no es super admin", valor: "0", set: true, want: false},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			c, _ := contextoDePrueba()
			if tc.set {
				c.Set("business_id", tc.valor)
			}

			assert.Equal(t, tc.want, IsSuperAdmin(c))
		})
	}
}

func TestGetScope_PrefiereElDelAuthInfo(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("auth_info", &domain.AuthInfo{Scope: "platform", BusinessID: 36})

	assert.Equal(t, "platform", GetScope(c),
		"el scope explicito del auth_info tiene prioridad sobre el business_id")
}

func TestGetScope_SinScope_LoDeduceDelBusiness(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("auth_info", &domain.AuthInfo{BusinessID: 0})

	assert.Equal(t, "platform", GetScope(c))

	c2, _ := contextoDePrueba()
	c2.Set("auth_info", &domain.AuthInfo{BusinessID: 36})

	assert.Equal(t, "business", GetScope(c2))
}

func TestGetScope_SinAutenticacion_CaeABusiness(t *testing.T) {
	c, _ := contextoDePrueba()

	assert.Equal(t, "business", GetScope(c),
		"sin datos se asume el scope menos privilegiado")
}

func TestRequireSuperAdmin_SinSerSuperAdmin_Rechaza(t *testing.T) {
	c, w := contextoDePrueba()
	c.Set("business_id", uint(36))

	RequireSuperAdmin()(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
	assert.Contains(t, w.Body.String(), "SUPER_ADMIN_REQUIRED")
}

func TestRequireSuperAdmin_SinAutenticacion_Rechaza(t *testing.T) {
	c, w := contextoDePrueba()

	RequireSuperAdmin()(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
}

func TestRequireSuperAdmin_SiendoSuperAdmin_Permite(t *testing.T) {
	c, w := contextoDePrueba()
	c.Set("business_id", uint(0))

	RequireSuperAdmin()(c)

	assert.False(t, c.IsAborted())
	assert.NotEqual(t, http.StatusForbidden, w.Code)
}

func TestRequireRole_SinRoles_Rechaza(t *testing.T) {
	c, w := contextoDePrueba()

	RequireRole("admin")(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
}

func TestRequireRole_SinElRol_Rechaza(t *testing.T) {
	c, w := contextoDePrueba()
	c.Set("user_roles", []string{"vendedor", "operario"})

	RequireRole("admin")(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
}

func TestRequireRole_ConElRol_Permite(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("user_roles", []string{"vendedor", "admin"})

	RequireRole("admin")(c)

	assert.False(t, c.IsAborted())
}

func TestRequireRole_EsSensibleAMayusculas(t *testing.T) {
	c, w := contextoDePrueba()
	c.Set("user_roles", []string{"Admin"})

	RequireRole("admin")(c)

	assert.Equal(t, http.StatusForbidden, w.Code,
		"la comparacion de rol es exacta, no normaliza mayusculas")
}

func TestRequireAnyRole_SinRoles_Rechaza(t *testing.T) {
	c, w := contextoDePrueba()

	RequireAnyRole("admin", "supervisor")(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
}

func TestRequireAnyRole_ConUnoDeLosRoles_Permite(t *testing.T) {
	c, _ := contextoDePrueba()
	c.Set("user_roles", []string{"supervisor"})

	RequireAnyRole("admin", "supervisor")(c)

	assert.False(t, c.IsAborted())
}

func TestRequireAnyRole_SinNingunRolCoincidente_Rechaza(t *testing.T) {
	c, w := contextoDePrueba()
	c.Set("user_roles", []string{"vendedor"})

	RequireAnyRole("admin", "supervisor")(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
}

func TestRequireAnyRole_SinRolesRequeridos_Rechaza(t *testing.T) {
	c, w := contextoDePrueba()
	c.Set("user_roles", []string{"admin"})

	RequireAnyRole()(c)

	assert.Equal(t, http.StatusForbidden, w.Code,
		"sin roles requeridos ninguno coincide, se rechaza")
}

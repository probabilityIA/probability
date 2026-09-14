package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/authz"
	"github.com/secamc93/probability/back/central/shared/log"
)

type fakeUC struct {
	access *entities.Access
	err    error
}

func (f *fakeUC) GetEffectiveAccess(context.Context, uint, uint, uint) (*entities.Access, error) {
	return f.access, f.err
}
func (f *fakeUC) Can(access *entities.Access, resource, action string) bool {
	if access.IsSuper {
		return true
	}
	want := authz.PermissionCode(resource, action)
	for _, p := range access.Permissions {
		if p == want {
			return true
		}
	}
	return false
}
func (f *fakeUC) Navigation(*entities.Access) []entities.NavItem { return nil }
func (f *fakeUC) SetModuleAccess(ports.IModuleAccess)            {}

var testTable = authz.NewPolicyTable(map[string]authz.RoutePolicy{
	"POST /auth/login": authz.Public(),
	"GET /tickets":     authz.SuperAdmin(),
	"GET /me":          authz.Authenticated(),
}, []authz.PrefixRule{{Prefix: "/orders", Policy: authz.Perm("orders")}})

func runGuard(t *testing.T, g *Guard, method, path string, businessID uint, authenticated bool) (int, bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	g.authenticate = func(c *gin.Context) bool {
		if !authenticated {
			return false
		}
		c.Set("user_id", uint(8))
		c.Set("business_id", businessID)
		return true
	}

	reached := false
	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(g.Handler())
	api.Handle(method, path, func(c *gin.Context) {
		reached = true
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(method, "/api/v1"+path, nil))
	return w.Code, reached
}

func newGuard(uc *fakeUC, mode string) *Guard {
	return &Guard{uc: uc, table: testTable, log: log.New(), mode: mode, enforceModules: map[string]bool{}, apiPrefix: "/api/v1"}
}

func adminAccess(perms ...string) *fakeUC {
	return &fakeUC{access: &entities.Access{BusinessID: 26, SubscriptionStatus: "active", Permissions: perms}}
}

func TestGuard_Enforce_SinPermisoEs403(t *testing.T) {
	code, reached := runGuard(t, newGuard(adminAccess("orders.read"), ModeEnforce), http.MethodDelete, "/orders", 26, true)
	if code != http.StatusForbidden || reached {
		t.Fatalf("se esperaba 403 sin llegar al handler, llego %d reached=%v", code, reached)
	}
}

func TestGuard_Enforce_ConPermisoPasa(t *testing.T) {
	code, reached := runGuard(t, newGuard(adminAccess("orders.read"), ModeEnforce), http.MethodGet, "/orders", 26, true)
	if code != http.StatusOK || !reached {
		t.Fatalf("se esperaba 200, llego %d", code)
	}
}

func TestGuard_Auditoria_NoBloquea(t *testing.T) {
	code, reached := runGuard(t, newGuard(adminAccess(), ModeAudit), http.MethodDelete, "/orders", 26, true)
	if code != http.StatusOK || !reached {
		t.Fatalf("auditoria no debe bloquear, llego %d", code)
	}
}

func TestGuard_Auditoria_ModuloEnEnforceSiBloquea(t *testing.T) {
	g := newGuard(adminAccess(), ModeAudit)
	g.enforceModules["orders"] = true
	code, _ := runGuard(t, g, http.MethodDelete, "/orders", 26, true)
	if code != http.StatusForbidden {
		t.Fatalf("el modulo orders esta en enforce, se esperaba 403, llego %d", code)
	}
}

func TestGuard_Enforce_RutaNoDeclaradaSeNiega(t *testing.T) {
	code, reached := runGuard(t, newGuard(adminAccess(), ModeEnforce), http.MethodGet, "/sin-politica", 26, true)
	if code != http.StatusForbidden || reached {
		t.Fatalf("ruta sin politica debe negarse, llego %d", code)
	}
}

func TestGuard_Publica_NoExigeSesion(t *testing.T) {
	code, reached := runGuard(t, newGuard(adminAccess(), ModeEnforce), http.MethodPost, "/auth/login", 0, false)
	if code != http.StatusOK || !reached {
		t.Fatalf("la ruta publica debe pasar, llego %d", code)
	}
}

func TestGuard_Enforce_SinSesionEs401(t *testing.T) {
	code, _ := runGuard(t, newGuard(adminAccess(), ModeEnforce), http.MethodGet, "/me", 0, false)
	if code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, llego %d", code)
	}
}

func TestGuard_Enforce_SuperAdminOnly(t *testing.T) {
	code, _ := runGuard(t, newGuard(adminAccess(), ModeEnforce), http.MethodGet, "/tickets", 26, true)
	if code != http.StatusForbidden {
		t.Fatalf("un negocio no entra a tickets, llego %d", code)
	}
	code, reached := runGuard(t, newGuard(&fakeUC{access: &entities.Access{IsSuper: true}}, ModeEnforce), http.MethodGet, "/tickets", 0, true)
	if code != http.StatusOK || !reached {
		t.Fatalf("el super admin entra a tickets, llego %d", code)
	}
}

func TestGuard_Enforce_CualquieraDeVariosPermisos(t *testing.T) {
	table := authz.NewPolicyTable(map[string]authz.RoutePolicy{
		"GET /integrations": authz.AnyPermission("integrations.read", "orders.read"),
	}, nil)
	g := newGuard(adminAccess("orders.read"), ModeEnforce)
	g.table = table
	code, reached := runGuard(t, g, http.MethodGet, "/integrations", 26, true)
	if code != http.StatusOK || !reached {
		t.Fatalf("con orders.read debe pasar, llego %d", code)
	}

	g2 := newGuard(adminAccess("customers.read"), ModeEnforce)
	g2.table = table
	code, _ = runGuard(t, g2, http.MethodGet, "/integrations", 26, true)
	if code != http.StatusForbidden {
		t.Fatalf("sin ninguno de los permisos debe dar 403, llego %d", code)
	}
}

func TestGuard_Enforce_SuscripcionVencidaBloqueaSalvoBilletera(t *testing.T) {
	uc := &fakeUC{access: &entities.Access{BusinessID: 26, SubscriptionStatus: "expired", Permissions: []string{"orders.read"}}}
	code, _ := runGuard(t, newGuard(uc, ModeEnforce), http.MethodGet, "/orders", 26, true)
	if code != http.StatusPaymentRequired {
		t.Fatalf("suscripcion vencida debe dar 402, llego %d", code)
	}
}

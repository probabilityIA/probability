package authz

import "testing"

func TestPoliticas_ResolucionPorPrefijoYMetodo(t *testing.T) {
	cases := []struct {
		method, path string
		kind         PolicyKind
		permission   string
	}{
		{"GET", "/orders/:id", PolicyPermission, "orders.read"},
		{"PUT", "/orders/:id/status", PolicyPermission, "orders.update"},
		{"DELETE", "/orders/:id", PolicyPermission, "orders.delete"},
		{"POST", "/orders/:id/request-confirmation", PolicyPermission, "orders.update"},
		{"GET", "/orders-compare", PolicyPermission, "orders.read"},
		{"POST", "/pay/wallet/admin/adjust-balance", PolicySuperAdmin, ""},
		{"GET", "/pay/wallet/balance", PolicyPermission, "wallet.read"},
		{"POST", "/auth/login", PolicyPublic, ""},
		{"POST", "/integrations/shopify/webhook/:integration_id", PolicyPublic, ""},
		{"GET", "/integration-types", PolicyAuthenticated, ""},
		{"POST", "/integration-types", PolicySuperAdmin, ""},
		{"GET", "/integration-types/:id/platform-credentials", PolicySuperAdmin, ""},
		{"GET", "/subscriptions/me", PolicyAuthenticated, ""},
		{"POST", "/subscriptions/extend-days", PolicySuperAdmin, ""},
		{"GET", "/tickets/:id", PolicySuperAdmin, ""},
	}
	for _, c := range cases {
		p, ok := RoutePolicies.Resolve(c.method, c.path)
		if !ok {
			t.Fatalf("%s %s sin politica", c.method, c.path)
		}
		if p.Kind != c.kind {
			t.Fatalf("%s %s: se esperaba %s, llego %s", c.method, c.path, c.kind, p.Kind)
		}
		if c.permission != "" && PermissionCode(p.Resource, p.Action) != c.permission {
			t.Fatalf("%s %s: se esperaba %s, llego %s", c.method, c.path, c.permission, PermissionCode(p.Resource, p.Action))
		}
	}
}

func TestPoliticas_PrefijoNoConfundeRutasParecidas(t *testing.T) {
	p, _ := RoutePolicies.Resolve("GET", "/orders-compare")
	if p.Resource != "orders" {
		t.Fatalf("orders-compare debe resolver a orders, llego %+v", p)
	}
	if _, ok := RoutePolicies.Resolve("GET", "/ordersx"); ok {
		t.Fatalf("/ordersx no debe coincidir con el prefijo /orders")
	}
}

func TestPoliticas_TodoRecursoReferenciadoExiste(t *testing.T) {
	for route, p := range buildExact() {
		if p.Kind == PolicyPermission {
			if _, ok := ResourceByCode(p.Resource); !ok {
				t.Fatalf("%s referencia un recurso inexistente: %s", route, p.Resource)
			}
		}
	}
	for _, rule := range prefixRules {
		if rule.Policy.Kind == PolicyPermission {
			if _, ok := ResourceByCode(rule.Policy.Resource); !ok {
				t.Fatalf("prefijo %s referencia un recurso inexistente: %s", rule.Prefix, rule.Policy.Resource)
			}
		}
	}
}

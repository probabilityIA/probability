package authz

import (
	"testing"

	"github.com/secamc93/probability/back/central/shared/moduleregistry"
)

func TestCatalogo_CodigosDeRecursoUnicos(t *testing.T) {
	seen := map[string]bool{}
	for _, r := range Resources {
		if seen[r.Code] {
			t.Fatalf("codigo de recurso repetido: %s", r.Code)
		}
		seen[r.Code] = true
	}
}

func TestCatalogo_NombresLegadosNoSeRepitenEntreRecursos(t *testing.T) {
	owner := map[string]string{}
	for _, r := range Resources {
		for _, name := range r.LegacyNames {
			key := normalize(name)
			if prev, ok := owner[key]; ok && prev != r.Code {
				t.Fatalf("el nombre %q apunta a %s y a %s", name, prev, r.Code)
			}
			owner[key] = r.Code
		}
	}
}

func TestCatalogo_ModulosValidos(t *testing.T) {
	for _, r := range Resources {
		if r.Module != "" && !moduleregistry.IsValid(r.Module) {
			t.Fatalf("recurso %s con modulo invalido %s", r.Code, r.Module)
		}
	}
}

func TestCatalogo_NavegacionConsistente(t *testing.T) {
	keys := map[string]bool{}
	routes := map[string]bool{}
	for _, n := range Navigation {
		if keys[n.Key] {
			t.Fatalf("clave de navegacion repetida: %s", n.Key)
		}
		keys[n.Key] = true
		if routes[n.Route] {
			t.Fatalf("ruta de navegacion repetida: %s", n.Route)
		}
		routes[n.Route] = true
		if n.Rule == NavRulePermission {
			if _, ok := ResourceByCode(n.Resource); !ok {
				t.Fatalf("la entrada %s exige un recurso inexistente: %s", n.Key, n.Resource)
			}
		}
		if n.Rule == NavRuleRole && len(n.RoleCodes) == 0 {
			t.Fatalf("la entrada %s es por rol y no declara roles", n.Key)
		}
	}
}

func TestCatalogo_LosTreintaYCuatroRecursosDeBDTienenCodigo(t *testing.T) {
	enBD := []string{"Usuarios", "Permisos", "Roles", "Recursos", "Empresas", "Ordenes", "Productos",
		"Integraciones", "Envios", "Facturacion", "Integraciones-Platform", "Integraciones-E-commerce",
		"Integraciones-Facturacion-Electronica", "Integraciones-Mensajeria", "Integraciones-Pagos",
		"Integraciones-Logistica", "Integraciones-Tipos-de-integracion", "Notificaciones", "Clientes",
		"Ultima Milla", "Billetera", "Inventario", "Bodegas", "Storefront", "Inventario-Stock",
		"Inventario-Movimientos", "Inventario-Trazabilidad", "Inventario-Kardex", "Inventario-Operaciones",
		"Inventario-Slotting", "Inventario-Auditoria", "Inventario-LPN", "Inventario-Scan", "Inventario-Sync-Logs"}
	for _, name := range enBD {
		if _, ok := ResourceByLegacyName(name); !ok {
			t.Fatalf("el recurso de BD %q no tiene codigo en el catalogo", name)
		}
	}
}

func TestRoleCode(t *testing.T) {
	cases := map[string]string{"Administrador": "administrador", "Super Admin": "super_admin", "cliente_final": "cliente_final", "demo": "demo"}
	for in, want := range cases {
		if got := RoleCode(in); got != want {
			t.Fatalf("RoleCode(%q) = %q, se esperaba %q", in, got, want)
		}
	}
}

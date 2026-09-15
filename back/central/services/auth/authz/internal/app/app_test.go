package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/log"
)

type fakeRepo struct {
	staff       map[uint]*entities.StaffRole
	superStaff  *entities.StaffRole
	permissions []entities.RolePermission
	configs     []entities.ResourceConfig
	names       map[uint]string
	staffCalls  int
}

func (f *fakeRepo) GetStaffRole(_ context.Context, _ uint, businessID *uint) (*entities.StaffRole, error) {
	f.staffCalls++
	if businessID == nil {
		return f.superStaff, nil
	}
	return f.staff[*businessID], nil
}
func (f *fakeRepo) GetRolePermissions(context.Context, uint) ([]entities.RolePermission, error) {
	return f.permissions, nil
}
func (f *fakeRepo) GetBusinessResourceConfig(context.Context, uint) ([]entities.ResourceConfig, error) {
	return f.configs, nil
}
func (f *fakeRepo) GetBusinessName(_ context.Context, id uint) (string, error) {
	return f.names[id], nil
}

type fakeModules struct{ modules []string }

func (f fakeModules) GetAccessibleModules(context.Context, uint) ([]string, error) {
	return f.modules, nil
}

type memoryCache struct{ data map[[2]uint]*entities.Access }

func (m *memoryCache) Get(_ context.Context, userID, businessID uint) (*entities.Access, bool) {
	a, ok := m.data[[2]uint{userID, businessID}]
	return a, ok
}
func (m *memoryCache) Set(_ context.Context, a *entities.Access) {
	m.data[[2]uint{a.UserID, a.BusinessID}] = a
}

func roleID(id uint) *uint { return &id }

func adminStaff() *entities.StaffRole {
	return &entities.StaffRole{BusinessName: "Demo", SubscriptionStatus: "active", RoleID: roleID(4), RoleName: "Administrador", RoleScopeCode: "business"}
}

func adminPermissions() []entities.RolePermission {
	return []entities.RolePermission{
		{ResourceName: "Ordenes", ActionName: "Read"},
		{ResourceName: "Ordenes", ActionName: "Update"},
		{ResourceName: "Inventario", ActionName: "Read"},
		{ResourceName: "Roles", ActionName: "Read"},
		{ResourceName: "Recursos", ActionName: "Read"},
		{ResourceName: "Recurso-Inventado", ActionName: "Read"},
	}
}

func newUC(repo *fakeRepo, modules []string) *UseCase {
	return &UseCase{repo: repo, modules: fakeModules{modules: modules}, log: log.New()}
}

func has(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func TestAcceso_ElPlanQuitaLosModulosNoContratados(t *testing.T) {
	repo := &fakeRepo{staff: map[uint]*entities.StaffRole{26: adminStaff()}, permissions: adminPermissions()}
	uc := newUC(repo, []string{"orders", "iam"})

	access, err := uc.GetEffectiveAccess(context.Background(), 8, 26, 0)
	if err != nil {
		t.Fatal(err)
	}
	if has(access.Permissions, "inventory.read") {
		t.Fatalf("inventory no esta en el plan y aparecio: %v", access.Permissions)
	}
	if !has(access.Permissions, "orders.read") || !has(access.Permissions, "orders.update") {
		t.Fatalf("faltan permisos de ordenes: %v", access.Permissions)
	}
}

func TestAcceso_RecursosSoloSuperAdminNuncaSeConcedenAUnNegocio(t *testing.T) {
	repo := &fakeRepo{staff: map[uint]*entities.StaffRole{26: adminStaff()}, permissions: adminPermissions()}
	uc := newUC(repo, []string{"orders", "iam", "inventory"})

	access, _ := uc.GetEffectiveAccess(context.Background(), 8, 26, 0)
	if has(access.Permissions, "resources.read") {
		t.Fatalf("resources es solo super admin: %v", access.Permissions)
	}
}

func TestAcceso_ConfiguracionDelNegocioRestringe(t *testing.T) {
	repo := &fakeRepo{
		staff:       map[uint]*entities.StaffRole{26: adminStaff()},
		permissions: adminPermissions(),
		configs: []entities.ResourceConfig{
			{ResourceName: "Ordenes", Active: true},
			{ResourceName: "Roles", Active: false},
		},
	}
	uc := newUC(repo, []string{"orders", "iam", "inventory"})

	access, _ := uc.GetEffectiveAccess(context.Background(), 8, 26, 0)
	if has(access.Permissions, "roles.read") {
		t.Fatalf("roles esta desactivado para el negocio: %v", access.Permissions)
	}
	if has(access.Permissions, "inventory.read") {
		t.Fatalf("inventario no esta configurado como activo: %v", access.Permissions)
	}
	if !has(access.Permissions, "orders.read") {
		t.Fatalf("ordenes esta activo: %v", access.Permissions)
	}
}

func TestAcceso_SinConfiguracionNoAbreTodo_ManDaElPlan(t *testing.T) {
	repo := &fakeRepo{staff: map[uint]*entities.StaffRole{26: adminStaff()}, permissions: adminPermissions()}
	uc := newUC(repo, []string{"orders"})

	access, _ := uc.GetEffectiveAccess(context.Background(), 8, 26, 0)
	for _, p := range access.Permissions {
		if p != "orders.read" && p != "orders.update" {
			t.Fatalf("sin configuracion solo debe quedar lo del plan, aparecio %s", p)
		}
	}
}

func TestAcceso_UsuarioSinRelacionConElNegocio(t *testing.T) {
	repo := &fakeRepo{staff: map[uint]*entities.StaffRole{}}
	uc := newUC(repo, []string{"orders"})

	_, err := uc.GetEffectiveAccess(context.Background(), 8, 46, 0)
	if !errors.Is(err, domainerrors.ErrNoBusinessRelation) {
		t.Fatalf("se esperaba ErrNoBusinessRelation, llego %v", err)
	}
}

func TestAcceso_SuperAdminLoPuedeTodo(t *testing.T) {
	repo := &fakeRepo{superStaff: &entities.StaffRole{RoleID: roleID(1), RoleName: "Super Admin", RoleScopeCode: "platform"}, names: map[uint]string{26: "Demo"}}
	uc := newUC(repo, nil)

	access, err := uc.GetEffectiveAccess(context.Background(), 1, 0, 26)
	if err != nil {
		t.Fatal(err)
	}
	if !access.IsSuper || access.BusinessName != "Demo" {
		t.Fatalf("acceso de super admin inesperado: %+v", access)
	}
	if !uc.Can(access, "invoicing", "delete") {
		t.Fatalf("el super admin debe poder todo")
	}
	nav := uc.Navigation(access)
	if !navHas(nav, "tickets") || !navHas(nav, "resources") {
		t.Fatalf("el super admin debe ver tickets y recursos")
	}
}

func TestNavegacion_DemoNoVeSuscripcionNiSeccionesDePlataforma(t *testing.T) {
	staff := adminStaff()
	staff.RoleName = "demo"
	repo := &fakeRepo{staff: map[uint]*entities.StaffRole{30: staff}, permissions: []entities.RolePermission{{ResourceName: "Ordenes", ActionName: "Read"}}}
	uc := newUC(repo, []string{"orders"})

	access, _ := uc.GetEffectiveAccess(context.Background(), 71, 30, 0)
	nav := uc.Navigation(access)
	if navHas(nav, "subscription") || navHas(nav, "tickets") || navHas(nav, "invoicing") {
		t.Fatalf("navegacion de demo con entradas prohibidas: %+v", nav)
	}
	if !navHas(nav, "orders") || !navHas(nav, "home") {
		t.Fatalf("demo debe ver inicio y ordenes: %+v", nav)
	}
}

func TestNavegacion_SitioWebSoloParaAdministrador(t *testing.T) {
	repo := &fakeRepo{staff: map[uint]*entities.StaffRole{26: adminStaff()}}
	uc := newUC(repo, []string{"orders"})

	access, _ := uc.GetEffectiveAccess(context.Background(), 8, 26, 0)
	if !navHas(uc.Navigation(access), "website_config") {
		t.Fatalf("el administrador debe ver sitio web")
	}
}

func TestAcceso_SeCachea(t *testing.T) {
	repo := &fakeRepo{staff: map[uint]*entities.StaffRole{26: adminStaff()}, permissions: adminPermissions()}
	uc := newUC(repo, []string{"orders"})
	uc.cache = &memoryCache{data: map[[2]uint]*entities.Access{}}

	_, _ = uc.GetEffectiveAccess(context.Background(), 8, 26, 0)
	_, _ = uc.GetEffectiveAccess(context.Background(), 8, 26, 0)
	if repo.staffCalls != 1 {
		t.Fatalf("la segunda consulta debia salir de cache, llamadas=%d", repo.staffCalls)
	}
}

func TestAcceso_SubrecursosDeInventarioHeredanDelPadre(t *testing.T) {
	repo := &fakeRepo{
		staff:       map[uint]*entities.StaffRole{30: adminStaff()},
		permissions: []entities.RolePermission{{ResourceName: "Inventario", ActionName: "Read"}},
		configs: []entities.ResourceConfig{
			{ResourceName: "Inventario", Active: true},
			{ResourceName: "Inventario-Stock", Active: true},
			{ResourceName: "Inventario-Kardex", Active: false},
		},
	}
	uc := newUC(repo, []string{"inventory"})

	access, _ := uc.GetEffectiveAccess(context.Background(), 8, 30, 0)
	if !has(access.Permissions, "inventory.stock.read") {
		t.Fatalf("stock activo debe heredar lectura del padre: %v", access.Permissions)
	}
	if has(access.Permissions, "inventory.kardex.read") {
		t.Fatalf("kardex esta desactivado para el negocio: %v", access.Permissions)
	}
}

func TestAcceso_DemoNoHeredaSubrecursosAvanzados(t *testing.T) {
	staff := adminStaff()
	staff.RoleName = "demo"
	repo := &fakeRepo{staff: map[uint]*entities.StaffRole{30: staff}, permissions: []entities.RolePermission{{ResourceName: "Inventario", ActionName: "Read"}}}
	uc := newUC(repo, []string{"inventory"})

	access, _ := uc.GetEffectiveAccess(context.Background(), 71, 30, 0)
	if !has(access.Permissions, "inventory.movements.read") {
		t.Fatalf("demo ve movimientos: %v", access.Permissions)
	}
	if has(access.Permissions, "inventory.kardex.read") || has(access.Permissions, "inventory.lpn.read") {
		t.Fatalf("demo no debe ver kardex ni LPN: %v", access.Permissions)
	}
}

func navHas(nav []entities.NavItem, key string) bool {
	for _, n := range nav {
		if n.Key == key {
			return true
		}
	}
	return false
}

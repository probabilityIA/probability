package middleware

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/app"
	authmw "github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/authz"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	ModeAudit   = "audit"
	ModeEnforce = "enforce"
)

var suspendedStatuses = map[string]bool{"expired": true, "cancelled": true}

var allowedWhileSuspended = map[string]bool{"wallet": true}

type Guard struct {
	uc             app.IUseCase
	table          *authz.PolicyTable
	log            log.ILogger
	mode           string
	enforceModules map[string]bool
	apiPrefix      string
	authenticate   func(c *gin.Context) bool
}

func New(uc app.IUseCase, table *authz.PolicyTable, logger log.ILogger, apiPrefix string) *Guard {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("AUTHZ_MODE")))
	if mode != ModeEnforce {
		mode = ModeAudit
	}
	enforceModules := map[string]bool{}
	for _, m := range strings.Split(os.Getenv("AUTHZ_ENFORCE_MODULES"), ",") {
		if m = strings.TrimSpace(m); m != "" {
			enforceModules[m] = true
		}
	}
	return &Guard{uc: uc, table: table, log: logger, mode: mode, enforceModules: enforceModules, apiPrefix: apiPrefix, authenticate: authmw.Authenticate}
}

type decision struct {
	allowed    bool
	status     int
	reason     string
	permission string
	module     string
}

func (g *Guard) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		fullPath := c.FullPath()
		if fullPath == "" {
			c.Next()
			return
		}

		d := g.evaluate(c, fullPath)
		if d.allowed {
			c.Next()
			return
		}

		userID, _ := authmw.GetUserID(c)
		businessID, _ := authmw.GetBusinessIDFromContext(c)
		enforced := g.shouldEnforce(d)

		event := g.log.Warn(c.Request.Context()).
			Str("authz_mode", g.modeFor(enforced)).
			Str("method", c.Request.Method).
			Str("route", fullPath).
			Str("reason", d.reason).
			Str("permission", d.permission).
			Uint("user_id", userID).
			Uint("business_id", businessID)
		if enforced {
			event.Msg("[authz] acceso denegado")
			message := denialMessage(d)
			c.AbortWithStatusJSON(d.status, gin.H{
				"success":    false,
				"code":       denialCode(d),
				"error":      message,
				"message":    message,
				"reason":     d.reason,
				"permission": d.permission,
			})
			return
		}
		event.Msg("[authz] denegacion en auditoria")
		c.Next()
	}
}

var actionVerbs = map[string]string{
	authz.ActionRead:   "consulta",
	authz.ActionCreate: "creaci\u00f3n",
	authz.ActionUpdate: "edici\u00f3n",
	authz.ActionDelete: "eliminaci\u00f3n",
}

func denialCode(d decision) string {
	switch d.reason {
	case "unauthenticated":
		return "unauthenticated"
	case "subscription_suspended":
		return "subscription_suspended"
	default:
		return "forbidden"
	}
}

func denialMessage(d decision) string {
	switch d.reason {
	case "unauthenticated":
		return "Tu sesi\u00f3n expir\u00f3. Vuelve a iniciar sesi\u00f3n."
	case "subscription_suspended":
		return "Tu suscripci\u00f3n est\u00e1 vencida. Ponte al d\u00eda en el m\u00f3dulo de Suscripci\u00f3n para seguir usando esta funci\u00f3n."
	case "super_admin_required":
		return "No tienes permiso para esta acci\u00f3n: es exclusiva del equipo de Probability."
	}
	if resourceCode, action, ok := splitPermission(d.permission); ok && !strings.Contains(d.permission, "|") {
		verb, hasVerb := actionVerbs[action]
		resource, hasResource := authz.ResourceByCode(resourceCode)
		if hasVerb && hasResource {
			return "No tienes permiso de " + verb + " en " + resource.Label + "."
		}
	}
	return "No tienes permiso para esta acci\u00f3n."
}

func (g *Guard) modeFor(enforced bool) string {
	if enforced {
		return ModeEnforce
	}
	return ModeAudit
}

func (g *Guard) shouldEnforce(d decision) bool {
	if g.mode == ModeEnforce {
		return true
	}
	return d.module != "" && g.enforceModules[d.module]
}

func (g *Guard) evaluate(c *gin.Context, fullPath string) decision {
	policy, ok := g.table.Resolve(c.Request.Method, strings.TrimPrefix(fullPath, g.apiPrefix))
	if !ok {
		return decision{status: http.StatusForbidden, reason: "undeclared_route", module: "undeclared"}
	}

	if policy.Kind == authz.PolicyPublic {
		return decision{allowed: true}
	}

	if !g.authenticate(c) {
		return decision{status: http.StatusUnauthorized, reason: "unauthenticated", module: "auth"}
	}

	userID, _ := authmw.GetUserID(c)
	tokenBusinessID, _ := authmw.GetBusinessIDFromContext(c)

	if policy.Kind == authz.PolicySuperAdmin {
		if tokenBusinessID == 0 {
			return decision{allowed: true}
		}
		return decision{status: http.StatusForbidden, reason: "super_admin_required", module: "platform"}
	}

	access, err := g.uc.GetEffectiveAccess(c.Request.Context(), userID, tokenBusinessID, 0)
	if err != nil {
		return decision{status: http.StatusForbidden, reason: "access_unresolved", module: "unresolved"}
	}
	if access.IsSuper || policy.Kind == authz.PolicyAuthenticated {
		return decision{allowed: true}
	}

	if len(policy.AnyOf) > 0 {
		if suspendedStatuses[access.SubscriptionStatus] {
			return decision{status: http.StatusPaymentRequired, reason: "subscription_suspended", permission: strings.Join(policy.AnyOf, "|"), module: "any"}
		}
		for _, candidate := range policy.AnyOf {
			if resourceCode, action, ok := splitPermission(candidate); ok && g.uc.Can(access, resourceCode, action) {
				return decision{allowed: true}
			}
		}
		return decision{status: http.StatusForbidden, reason: "missing_permission", permission: strings.Join(policy.AnyOf, "|"), module: "any"}
	}

	resource, _ := authz.ResourceByCode(policy.Resource)
	permission := authz.PermissionCode(policy.Resource, policy.Action)

	if suspendedStatuses[access.SubscriptionStatus] && !allowedWhileSuspended[policy.Resource] {
		return decision{status: http.StatusPaymentRequired, reason: "subscription_suspended", permission: permission, module: resource.Module}
	}

	if !g.uc.Can(access, policy.Resource, policy.Action) {
		return decision{status: http.StatusForbidden, reason: "missing_permission", permission: permission, module: resource.Module}
	}
	return decision{allowed: true}
}

func splitPermission(code string) (string, string, bool) {
	idx := strings.LastIndex(code, ".")
	if idx <= 0 || idx == len(code)-1 {
		return "", "", false
	}
	return code[:idx], code[idx+1:], true
}

func Coverage(table *authz.PolicyTable, routes gin.RoutesInfo, apiPrefix string) []string {
	missing := make([]string, 0)
	for _, r := range routes {
		if !strings.HasPrefix(r.Path, apiPrefix) {
			continue
		}
		if _, ok := table.Resolve(r.Method, strings.TrimPrefix(r.Path, apiPrefix)); !ok {
			missing = append(missing, r.Method+" "+r.Path)
		}
	}
	sort.Strings(missing)
	return missing
}

func DumpRoutes(table *authz.PolicyTable, routes gin.RoutesInfo, apiPrefix, path string) error {
	lines := make([]string, 0, len(routes))
	for _, r := range routes {
		if !strings.HasPrefix(r.Path, apiPrefix) {
			continue
		}
		p, ok := table.Resolve(r.Method, strings.TrimPrefix(r.Path, apiPrefix))
		policy := "UNDECLARED"
		if ok {
			policy = string(p.Kind)
			if p.Kind == authz.PolicyPermission {
				policy += ":" + p.Resource + "." + p.Action
			}
		}
		lines = append(lines, fmt.Sprintf("%s %s %s", r.Method, strings.TrimPrefix(r.Path, apiPrefix), policy))
	}
	sort.Strings(lines)
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

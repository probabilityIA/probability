package authz

import (
	"net/http"
	"sort"
	"strings"
)

type PolicyKind string

const (
	PolicyPublic        PolicyKind = "public"
	PolicyAuthenticated PolicyKind = "authenticated"
	PolicySuperAdmin    PolicyKind = "super_admin"
	PolicyPermission    PolicyKind = "permission"
)

type RoutePolicy struct {
	Kind     PolicyKind
	Resource string
	Action   string
	AnyOf    []string
}

func AnyPermission(permissions ...string) RoutePolicy {
	return RoutePolicy{Kind: PolicyPermission, AnyOf: permissions}
}

type PrefixRule struct {
	Prefix string
	Policy RoutePolicy
}

func Public() RoutePolicy        { return RoutePolicy{Kind: PolicyPublic} }
func Authenticated() RoutePolicy { return RoutePolicy{Kind: PolicyAuthenticated} }
func SuperAdmin() RoutePolicy    { return RoutePolicy{Kind: PolicySuperAdmin} }

func Perm(resource string) RoutePolicy {
	return RoutePolicy{Kind: PolicyPermission, Resource: resource}
}

func PermAction(resource, action string) RoutePolicy {
	return RoutePolicy{Kind: PolicyPermission, Resource: resource, Action: action}
}

type PolicyTable struct {
	exact    map[string]RoutePolicy
	prefixes []PrefixRule
}

func NewPolicyTable(exact map[string]RoutePolicy, prefixes []PrefixRule) *PolicyTable {
	sorted := make([]PrefixRule, len(prefixes))
	copy(sorted, prefixes)
	sort.SliceStable(sorted, func(i, j int) bool {
		return len(sorted[i].Prefix) > len(sorted[j].Prefix)
	})
	return &PolicyTable{exact: exact, prefixes: sorted}
}

func (t *PolicyTable) Resolve(method, fullPath string) (RoutePolicy, bool) {
	if p, ok := t.exact[method+" "+fullPath]; ok {
		return withDefaultAction(p, method), true
	}
	for _, rule := range t.prefixes {
		if matchesPrefix(fullPath, rule.Prefix) {
			return withDefaultAction(rule.Policy, method), true
		}
	}
	return RoutePolicy{}, false
}

func matchesPrefix(path, prefix string) bool {
	if path == prefix {
		return true
	}
	return strings.HasPrefix(path, strings.TrimSuffix(prefix, "/")+"/")
}

func withDefaultAction(p RoutePolicy, method string) RoutePolicy {
	if p.Kind != PolicyPermission || p.Action != "" || len(p.AnyOf) > 0 {
		return p
	}
	p.Action = ActionForMethod(method)
	return p
}

func ActionForMethod(method string) string {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return ActionRead
	case http.MethodPost:
		return ActionCreate
	case http.MethodPut, http.MethodPatch:
		return ActionUpdate
	case http.MethodDelete:
		return ActionDelete
	default:
		return ActionRead
	}
}

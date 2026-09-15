package response

import "github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/entities"

type Business struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Role struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Subscription struct {
	Status string `json:"status"`
}

type NavItem struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Route       string `json:"route"`
	Section     string `json:"section"`
	Description string `json:"description"`
}

type Access struct {
	IsSuper      bool         `json:"is_super"`
	Business     *Business    `json:"business"`
	Role         Role         `json:"role"`
	Subscription Subscription `json:"subscription"`
	Modules      []string     `json:"modules"`
	Permissions  []string     `json:"permissions"`
	Navigation   []NavItem    `json:"navigation"`
}

func FromAccess(access *entities.Access, nav []entities.NavItem) Access {
	out := Access{
		IsSuper:      access.IsSuper,
		Role:         Role{ID: access.RoleID, Code: access.RoleCode, Name: access.RoleName},
		Subscription: Subscription{Status: access.SubscriptionStatus},
		Modules:      access.Modules,
		Permissions:  access.Permissions,
		Navigation:   make([]NavItem, len(nav)),
	}
	if out.Modules == nil {
		out.Modules = []string{}
	}
	if out.Permissions == nil {
		out.Permissions = []string{}
	}
	if access.BusinessID > 0 {
		out.Business = &Business{ID: access.BusinessID, Name: access.BusinessName}
	}
	for i, n := range nav {
		out.Navigation[i] = NavItem{Key: n.Key, Label: n.Label, Route: n.Route, Section: n.Section, Description: n.Description}
	}
	return out
}

package repository

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (r *Repository) ListDaneStates(ctx context.Context) ([]domain.DaneItem, error) {
	var rows []struct {
		Code string
		Name string
	}
	err := r.db.Conn(ctx).
		Table("geozones").
		Select("code, name").
		Where("type = ? AND deleted_at IS NULL AND code IS NOT NULL AND is_active = ?", "state", true).
		Order("name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.DaneItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.DaneItem{Code: row.Code, Name: row.Name})
	}
	return out, nil
}

var daneStateAliases = map[string]string{
	"capital district": "11",
	"distrito capital": "11",
	"bogota":           "11",
	"bogota d.c.":      "11",
	"bogota, d.c.":     "11",
	"bogota dc":        "11",
	"santafe de bogota": "11",
	"san andres y providencia": "88",
	"valle":            "76",
	"norte de santander (n. de santander)": "54",
}

func normalizeStateAlias(state string) string {
	key := strings.ToLower(strings.TrimSpace(state))
	repl := strings.NewReplacer("\u00e1", "a", "\u00e9", "e", "\u00ed", "i", "\u00f3", "o", "\u00fa", "u", "\u00f1", "n")
	key = repl.Replace(key)
	if code, ok := daneStateAliases[key]; ok {
		return code
	}
	return state
}

func (r *Repository) ListDaneCitiesByState(ctx context.Context, state string) ([]domain.DaneItem, error) {
	state = normalizeStateAlias(state)
	var rows []struct {
		Code string
		Name string
	}
	err := r.db.Conn(ctx).
		Table("geozones g").
		Select("g.code, g.name").
		Joins("JOIN geozones p ON p.id = g.parent_id AND p.type = 'state'").
		Where("g.type = ? AND g.deleted_at IS NULL AND g.code IS NOT NULL AND g.is_active = ?", "city", true).
		Where(`p.code = @state
		  OR unaccent(lower(p.name)) = unaccent(lower(@state))
		  OR unaccent(lower(p.name)) LIKE unaccent(lower(@state)) || '%'
		  OR unaccent(lower(@state)) LIKE unaccent(lower(p.name)) || '%'`,
			map[string]interface{}{"state": state}).
		Order("g.name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.DaneItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.DaneItem{Code: row.Code, Name: row.Name})
	}
	return out, nil
}

func (r *Repository) SearchDaneCities(ctx context.Context, state, term string, limit int) ([]domain.DaneItem, error) {
	term = strings.TrimSpace(term)
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	db := r.db.Conn(ctx).
		Table("geozones g").
		Select("g.code, g.name, p.name AS state_name").
		Joins("JOIN geozones p ON p.id = g.parent_id AND p.type = 'state'").
		Where("g.type = ? AND g.deleted_at IS NULL AND g.code IS NOT NULL AND g.is_active = ?", "city", true)

	if state = normalizeStateAlias(strings.TrimSpace(state)); state != "" {
		db = db.Where(`p.code = @state
		  OR unaccent(lower(p.name)) = unaccent(lower(@state))
		  OR unaccent(lower(p.name)) LIKE unaccent(lower(@state)) || '%'
		  OR unaccent(lower(@state)) LIKE unaccent(lower(p.name)) || '%'`,
			map[string]interface{}{"state": state})
	}

	if term != "" {
		db = db.Where(`unaccent(lower(g.name)) LIKE '%' || unaccent(lower(@term)) || '%'`,
			map[string]interface{}{"term": term}).
			Order(gorm.Expr(`(unaccent(lower(g.name)) = unaccent(lower(?))) DESC,
			  (unaccent(lower(g.name)) LIKE unaccent(lower(?)) || '%') DESC,
			  length(g.name) ASC, g.name ASC`, term, term))
	} else {
		db = db.Order("g.name ASC")
	}

	var rows []struct {
		Code      string
		Name      string
		StateName string
	}
	if err := db.Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]domain.DaneItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.DaneItem{Code: row.Code, Name: row.Name, StateName: row.StateName})
	}
	return out, nil
}

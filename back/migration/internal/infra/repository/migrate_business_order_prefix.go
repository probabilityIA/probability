package repository

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateBusinessOrderPrefix(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Business{}); err != nil {
		return fmt.Errorf("failed to auto-migrate business for order_prefix: %w", err)
	}

	var businesses []models.Business
	if err := r.db.Conn(ctx).
		Where("order_prefix IS NULL OR order_prefix = ''").
		Find(&businesses).Error; err != nil {
		return fmt.Errorf("failed to load businesses without order_prefix: %w", err)
	}

	if len(businesses) == 0 {
		return nil
	}

	taken := make(map[string]bool)
	if err := r.db.Conn(ctx).
		Model(&models.Business{}).
		Where("order_prefix IS NOT NULL AND order_prefix <> ''").
		Pluck("order_prefix", &[]string{}).Error; err == nil {
		var prefixes []string
		r.db.Conn(ctx).Model(&models.Business{}).
			Where("order_prefix IS NOT NULL AND order_prefix <> ''").
			Pluck("order_prefix", &prefixes)
		for _, p := range prefixes {
			taken[strings.ToUpper(p)] = true
		}
	}

	for i := range businesses {
		b := &businesses[i]
		base := derivePrefix(b.Name)
		prefix := base
		suffix := 2
		for taken[prefix] {
			prefix = fmt.Sprintf("%s%d", base, suffix)
			suffix++
		}
		taken[prefix] = true
		if err := r.db.Conn(ctx).
			Model(&models.Business{}).
			Where("id = ?", b.ID).
			Update("order_prefix", prefix).Error; err != nil {
			return fmt.Errorf("failed to backfill order_prefix for business %d: %w", b.ID, err)
		}
	}

	return nil
}

func derivePrefix(name string) string {
	letters := make([]rune, 0, 3)
	for _, r := range name {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToUpper(r))
			if len(letters) == 3 {
				break
			}
		}
	}
	if len(letters) == 0 {
		return "BIZ"
	}
	for len(letters) < 3 {
		letters = append(letters, 'X')
	}
	return string(letters)
}

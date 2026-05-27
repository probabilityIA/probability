package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateGuideFormats(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.GuideFormat{}); err != nil {
		return fmt.Errorf("automigrate guide_formats: %w", err)
	}
	return r.seedGuideFormats(ctx)
}

func (r *Repository) seedGuideFormats(ctx context.Context) error {
	obsolete := []string{
		"interrapidisimo-compact",
		"envia-original", "envia-10x15",
		"coordinadora-original", "coordinadora-10x15",
		"envioclick-original", "envioclick-10x15",
		"interrapidisimo-10x15",
		"tcc-original", "tcc-10x15",
		"servientrega-original", "servientrega-10x15",
		"99minutos-10x15",
		"deprisa-10x15",
	}
	if err := r.db.Conn(ctx).Where("code IN ?", obsolete).Delete(&models.GuideFormat{}).Error; err != nil {
		return fmt.Errorf("cleanup obsolete guide formats: %w", err)
	}

	seeds := []models.GuideFormat{
		{Carrier: "ENVIA", Code: "envia-compact", Label: "Guia recortada", WidthCm: 21.6, HeightCm: 9.3, Strategy: "crop", CropLLyFrac: 0.667, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "COORDINADORA", Code: "coordinadora-compact", Label: "Guia recortada", WidthCm: 21.6, HeightCm: 9.3, Strategy: "crop", CropLLyFrac: 0.667, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "ENVIOCLICK", Code: "envioclick-compact", Label: "Guia recortada", WidthCm: 21.6, HeightCm: 9.3, Strategy: "crop", CropLLyFrac: 0.667, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "INTERRAPIDISIMO", Code: "interrapidisimo-original", Label: "Guia recortada", WidthCm: 21, HeightCm: 29.7, Strategy: "passthrough", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "TCC", Code: "tcc-compact", Label: "Guia recortada", WidthCm: 21.6, HeightCm: 8.4, Strategy: "crop", CropLLyFrac: 0.7, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "SERVIENTREGA", Code: "servientrega-compact", Label: "Guia recortada", WidthCm: 21.6, HeightCm: 14, Strategy: "crop", CropLLyFrac: 0.5, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "99MINUTOS", Code: "99minutos-original", Label: "Guia recortada", WidthCm: 21.6, HeightCm: 27.9, Strategy: "passthrough", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "DEPRISA", Code: "deprisa-original", Label: "Guia recortada", WidthCm: 10.4, HeightCm: 11.2, Strategy: "passthrough", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},

		{Carrier: "*", Code: "probability-10x15", Label: "Adhesiva 10x15 cm", WidthCm: 10, HeightCm: 15, Adhesive: true, Strategy: "rebuild", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 100},
		{Carrier: "*", Code: "probability-4x6", Label: "Termica 4x6 in", WidthCm: 10.16, HeightCm: 15.24, Adhesive: true, Strategy: "rebuild", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 110},
		{Carrier: "*", Code: "probability-6x10", Label: "Mini 6x10 cm", WidthCm: 6, HeightCm: 10, Adhesive: true, Strategy: "rebuild", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 120},
		{Carrier: "*", Code: "probability-letter", Label: "Carta completa", WidthCm: 21.6, HeightCm: 27.9, Strategy: "rebuild", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 130},
	}

	for _, sd := range seeds {
		var existing models.GuideFormat
		err := r.db.Conn(ctx).Where("code = ?", sd.Code).First(&existing).Error
		if err == nil {
			existing.Carrier = sd.Carrier
			existing.Label = sd.Label
			existing.WidthCm = sd.WidthCm
			existing.HeightCm = sd.HeightCm
			existing.Adhesive = sd.Adhesive
			existing.Strategy = sd.Strategy
			existing.CropLLxFrac = sd.CropLLxFrac
			existing.CropLLyFrac = sd.CropLLyFrac
			existing.CropURxFrac = sd.CropURxFrac
			existing.CropURyFrac = sd.CropURyFrac
			existing.SourcePage = sd.SourcePage
			existing.IsDefault = sd.IsDefault
			existing.SortOrder = sd.SortOrder
			existing.IsActive = true
			if err := r.db.Conn(ctx).Save(&existing).Error; err != nil {
				return fmt.Errorf("update guide format %s: %w", sd.Code, err)
			}
			continue
		}
		seed := sd
		seed.IsActive = true
		if err := r.db.Conn(ctx).Create(&seed).Error; err != nil {
			return fmt.Errorf("create guide format %s: %w", sd.Code, err)
		}
	}
	return nil
}

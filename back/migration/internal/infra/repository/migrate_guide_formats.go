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
	if err := r.db.Conn(ctx).Where("code IN ?", []string{"interrapidisimo-compact"}).Delete(&models.GuideFormat{}).Error; err != nil {
		return fmt.Errorf("cleanup obsolete guide formats: %w", err)
	}

	seeds := []models.GuideFormat{
		// ENVIA — letter, 3 copias verticales
		{Carrier: "ENVIA", Code: "envia-original", Label: "Original carta", WidthCm: 21.6, HeightCm: 27.9, Strategy: "passthrough", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 30},
		{Carrier: "ENVIA", Code: "envia-compact", Label: "Compacta (1 guia)", WidthCm: 21.6, HeightCm: 9.3, Strategy: "crop", CropLLyFrac: 0.667, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "ENVIA", Code: "envia-10x15", Label: "Adhesiva 10x15 cm", WidthCm: 10, HeightCm: 15, Adhesive: true, Strategy: "resize", CropLLyFrac: 0.667, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 20},

		// COORDINADORA — letter, 1 copia top 1/3
		{Carrier: "COORDINADORA", Code: "coordinadora-original", Label: "Original carta", WidthCm: 21.6, HeightCm: 27.9, Strategy: "passthrough", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 30},
		{Carrier: "COORDINADORA", Code: "coordinadora-compact", Label: "Compacta (1 guia)", WidthCm: 21.6, HeightCm: 9.3, Strategy: "crop", CropLLyFrac: 0.667, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "COORDINADORA", Code: "coordinadora-10x15", Label: "Adhesiva 10x15 cm", WidthCm: 10, HeightCm: 15, Adhesive: true, Strategy: "resize", CropLLyFrac: 0.667, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 20},

		// ENVIOCLICK — alias de Coordinadora
		{Carrier: "ENVIOCLICK", Code: "envioclick-original", Label: "Original carta", WidthCm: 21.6, HeightCm: 27.9, Strategy: "passthrough", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 30},
		{Carrier: "ENVIOCLICK", Code: "envioclick-compact", Label: "Compacta (1 guia)", WidthCm: 21.6, HeightCm: 9.3, Strategy: "crop", CropLLyFrac: 0.667, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "ENVIOCLICK", Code: "envioclick-10x15", Label: "Adhesiva 10x15 cm", WidthCm: 10, HeightCm: 15, Adhesive: true, Strategy: "resize", CropLLyFrac: 0.667, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 20},

		// INTERRAPIDISIMO — A4 horizontal con rotacion 90, ocupa toda la pagina (1 sola guia)
		{Carrier: "INTERRAPIDISIMO", Code: "interrapidisimo-original", Label: "Original A4", WidthCm: 21, HeightCm: 29.7, Strategy: "passthrough", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "INTERRAPIDISIMO", Code: "interrapidisimo-10x15", Label: "Adhesiva 10x15 cm", WidthCm: 10, HeightCm: 15, Adhesive: true, Strategy: "resize", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 20},

		// TCC — letter, contenido top ~30%
		{Carrier: "TCC", Code: "tcc-original", Label: "Original carta", WidthCm: 21.6, HeightCm: 27.9, Strategy: "passthrough", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 30},
		{Carrier: "TCC", Code: "tcc-compact", Label: "Compacta (1 guia)", WidthCm: 21.6, HeightCm: 8.4, Strategy: "crop", CropLLyFrac: 0.7, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "TCC", Code: "tcc-10x15", Label: "Adhesiva 10x15 cm", WidthCm: 10, HeightCm: 15, Adhesive: true, Strategy: "resize", CropLLyFrac: 0.7, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 20},

		// SERVIENTREGA — letter, 4 copias 2x2 (tomamos la mitad superior)
		{Carrier: "SERVIENTREGA", Code: "servientrega-original", Label: "Original carta (4 copias)", WidthCm: 21.6, HeightCm: 27.9, Strategy: "passthrough", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 30},
		{Carrier: "SERVIENTREGA", Code: "servientrega-compact", Label: "Compacta (1 guia)", WidthCm: 21.6, HeightCm: 14, Strategy: "crop", CropLLyFrac: 0.5, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "SERVIENTREGA", Code: "servientrega-10x15", Label: "Adhesiva 10x15 cm", WidthCm: 10, HeightCm: 15, Adhesive: true, Strategy: "resize", CropLLyFrac: 0.5, CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 20},

		// 99MINUTOS — letter completa, ya viene 1 guia
		{Carrier: "99MINUTOS", Code: "99minutos-original", Label: "Original carta", WidthCm: 21.6, HeightCm: 27.9, Strategy: "passthrough", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "99MINUTOS", Code: "99minutos-10x15", Label: "Adhesiva 10x15 cm", WidthCm: 10, HeightCm: 15, Adhesive: true, Strategy: "resize", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 20},

		// DEPRISA — ya viene chico (~10x11), passthrough es ideal
		{Carrier: "DEPRISA", Code: "deprisa-original", Label: "Original (10x11 cm)", WidthCm: 10.4, HeightCm: 11.2, Strategy: "passthrough", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, IsDefault: true, SortOrder: 10},
		{Carrier: "DEPRISA", Code: "deprisa-10x15", Label: "Adhesiva 10x15 cm", WidthCm: 10, HeightCm: 15, Adhesive: true, Strategy: "resize", CropURxFrac: 1, CropURyFrac: 1, SourcePage: 1, SortOrder: 20},
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

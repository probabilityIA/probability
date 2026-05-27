package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type guideFormatRow struct {
	ID          uint
	Carrier     string
	Code        string
	Label       string
	WidthCm     float64
	HeightCm    float64
	Adhesive    bool
	Strategy    string
	CropLLxFrac float64
	CropLLyFrac float64
	CropURxFrac float64
	CropURyFrac float64
	SourcePage  int
	IsDefault   bool
	SortOrder   int
}

func (g *guideFormatRow) toDomain() *domain.GuideFormat {
	return &domain.GuideFormat{
		ID:          g.ID,
		Carrier:     g.Carrier,
		Code:        g.Code,
		Label:       g.Label,
		WidthCm:     g.WidthCm,
		HeightCm:    g.HeightCm,
		Adhesive:    g.Adhesive,
		Strategy:    g.Strategy,
		CropLLxFrac: g.CropLLxFrac,
		CropLLyFrac: g.CropLLyFrac,
		CropURxFrac: g.CropURxFrac,
		CropURyFrac: g.CropURyFrac,
		SourcePage:  g.SourcePage,
		IsDefault:   g.IsDefault,
		SortOrder:   g.SortOrder,
	}
}

const guideFormatSelect = `SELECT id, carrier, code, label, width_cm, height_cm, adhesive, strategy,
		crop_l_lx_frac AS crop_l_lx_frac,
		crop_l_ly_frac AS crop_l_ly_frac,
		crop_u_rx_frac AS crop_u_rx_frac,
		crop_u_ry_frac AS crop_u_ry_frac,
		source_page, is_default, sort_order
	FROM guide_formats
	WHERE deleted_at IS NULL AND is_active = TRUE`

func (r *Repository) scanGuideFormats(ctx context.Context, sql string, args ...any) ([]domain.GuideFormat, error) {
	var rows []guideFormatRow
	if err := r.db.Conn(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.GuideFormat, 0, len(rows))
	for i := range rows {
		out = append(out, *rows[i].toDomain())
	}
	return out, nil
}

func (r *Repository) GetGuideFormatByCode(ctx context.Context, code string) (*domain.GuideFormat, error) {
	rows, err := r.scanGuideFormats(ctx, guideFormatSelect+" AND code = ? LIMIT 1", code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

func (r *Repository) ListGuideFormatsByCarrier(ctx context.Context, carrier string) ([]domain.GuideFormat, error) {
	return r.scanGuideFormats(ctx, guideFormatSelect+" AND UPPER(carrier) = UPPER(?) ORDER BY sort_order, id", strings.TrimSpace(carrier))
}

func (r *Repository) ListGuideFormats(ctx context.Context) ([]domain.GuideFormat, error) {
	return r.scanGuideFormats(ctx, guideFormatSelect+" ORDER BY carrier, sort_order, id")
}

func (r *Repository) GetDefaultGuideFormat(ctx context.Context, carrier string) (*domain.GuideFormat, error) {
	rows, err := r.scanGuideFormats(ctx, guideFormatSelect+" AND UPPER(carrier) = UPPER(?) AND is_default = TRUE ORDER BY sort_order LIMIT 1", strings.TrimSpace(carrier))
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		fallback, err := r.scanGuideFormats(ctx, guideFormatSelect+" AND UPPER(carrier) = UPPER(?) ORDER BY sort_order LIMIT 1", strings.TrimSpace(carrier))
		if err != nil || len(fallback) == 0 {
			return nil, err
		}
		return &fallback[0], nil
	}
	return &rows[0], nil
}

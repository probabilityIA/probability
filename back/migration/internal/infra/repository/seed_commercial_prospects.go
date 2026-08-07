package repository

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/migration/shared/models"
)

//go:embed data/commercial_prospects_seed.json
var commercialProspectsSeedJSON []byte

type commercialProspectSeedRow struct {
	Domain                 string  `json:"domain"`
	MeliSellerID           *string `json:"meli_seller_id"`
	BusinessName           *string `json:"business_name"`
	Platform               string  `json:"platform"`
	City                   *string `json:"city"`
	Phone                  *string `json:"phone"`
	WhatsAppNumber         *string `json:"whatsapp_number"`
	InstagramHandle        *string `json:"instagram_handle"`
	ProductCount           *int    `json:"product_count"`
	HasCOD                 int     `json:"has_cod"`
	NumCouriersMentioned   int     `json:"num_couriers_mentioned"`
	UsesProbabilityChannel int     `json:"uses_probability_channel"`
	RobotsOK               int     `json:"robots_ok"`
	DiscoveryQuery         *string `json:"discovery_query"`
	Score                  float64 `json:"score"`
	ScoreBreakdown         string  `json:"score_breakdown"`
	Status                 string  `json:"status"`
	DiscoveredAt           *string `json:"discovered_at"`
	LastScannedAt          *string `json:"last_scanned_at"`
	InstagramFollowers     *int    `json:"instagram_followers"`
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func intOrZero(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func parseSeedTime(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339Nano, *s)
	if err != nil {
		return nil
	}
	return &t
}

func (r *Repository) seedCommercialProspects(ctx context.Context) error {
	db := r.db.Conn(ctx)

	var count int64
	if err := db.Model(&models.CommercialProspect{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check commercial_prospects: %w", err)
	}
	if count > 0 {
		return nil
	}

	var rows []commercialProspectSeedRow
	if err := json.Unmarshal(commercialProspectsSeedJSON, &rows); err != nil {
		return fmt.Errorf("failed to parse commercial prospects seed json: %w", err)
	}

	const batchSize = 200
	for i := 0; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}

		batch := make([]models.CommercialProspect, 0, end-i)
		for _, row := range rows[i:end] {
			if row.Domain == "" {
				continue
			}
			batch = append(batch, models.CommercialProspect{
				Domain:                 row.Domain,
				MeliSellerID:           strOrEmpty(row.MeliSellerID),
				BusinessName:           strOrEmpty(row.BusinessName),
				Platform:               row.Platform,
				City:                   strOrEmpty(row.City),
				Phone:                  strOrEmpty(row.Phone),
				WhatsAppNumber:         strOrEmpty(row.WhatsAppNumber),
				InstagramHandle:        strOrEmpty(row.InstagramHandle),
				InstagramFollowers:     intOrZero(row.InstagramFollowers),
				ProductCount:           intOrZero(row.ProductCount),
				HasCOD:                 row.HasCOD != 0,
				NumCouriersMentioned:   row.NumCouriersMentioned,
				UsesProbabilityChannel: row.UsesProbabilityChannel != 0,
				RobotsOK:               row.RobotsOK != 0,
				DiscoveryQuery:         strOrEmpty(row.DiscoveryQuery),
				Score:                  row.Score,
				ScoreBreakdown:         []byte(row.ScoreBreakdown),
				Status:                 row.Status,
				SourceDiscoveredAt:     parseSeedTime(row.DiscoveredAt),
				SourceLastScannedAt:    parseSeedTime(row.LastScannedAt),
			})
		}

		if len(batch) == 0 {
			continue
		}
		if err := db.Create(&batch).Error; err != nil {
			return fmt.Errorf("failed to insert commercial_prospects batch %d: %w", i/batchSize, err)
		}
	}

	fmt.Printf("seed commercial prospects: %d filas importadas\n", len(rows))
	return nil
}

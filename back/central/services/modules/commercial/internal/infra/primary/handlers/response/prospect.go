package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/entities"
)

type ProspectResponse struct {
	ID                     uint                   `json:"id"`
	Domain                 string                 `json:"domain"`
	ProfileURL             string                 `json:"profile_url"`
	MeliSellerID           string                 `json:"meli_seller_id,omitempty"`
	BusinessName           string                 `json:"business_name,omitempty"`
	Platform               string                 `json:"platform"`
	City                   string                 `json:"city,omitempty"`
	Phone                  string                 `json:"phone,omitempty"`
	WhatsAppNumber         string                 `json:"whatsapp_number,omitempty"`
	WhatsAppURL            string                 `json:"whatsapp_url,omitempty"`
	InstagramHandle        string                 `json:"instagram_handle,omitempty"`
	InstagramURL           string                 `json:"instagram_url,omitempty"`
	InstagramFollowers     int                    `json:"instagram_followers"`
	ProductCount           int                    `json:"product_count"`
	HasCOD                 bool                   `json:"has_cod"`
	NumCouriersMentioned   int                    `json:"num_couriers_mentioned"`
	UsesProbabilityChannel bool                   `json:"uses_probability_channel"`
	DiscoveryQuery         string                 `json:"discovery_query,omitempty"`
	Score                  float64                `json:"score"`
	ScoreBreakdown         map[string]interface{} `json:"score_breakdown,omitempty"`
	Status                 string                 `json:"status"`
	Seen                   bool                   `json:"seen"`
	SourceDiscoveredAt     *time.Time             `json:"source_discovered_at,omitempty"`
	SourceLastScannedAt    *time.Time             `json:"source_last_scanned_at,omitempty"`
}

func buildProfileURL(p *entities.Prospect) string {
	switch p.Platform {
	case "instagram_only":
		if p.InstagramHandle != "" {
			return "https://instagram.com/" + p.InstagramHandle
		}
	default:
		if p.Domain != "" {
			return "https://" + p.Domain
		}
	}
	return ""
}

func buildWhatsAppURL(number string) string {
	if number == "" {
		return ""
	}
	return "https://wa.me/" + number
}

func buildInstagramURL(handle string) string {
	if handle == "" {
		return ""
	}
	return "https://instagram.com/" + handle
}

func FromEntity(p *entities.Prospect) ProspectResponse {
	return ProspectResponse{
		ID:                     p.ID,
		Domain:                 p.Domain,
		ProfileURL:             buildProfileURL(p),
		MeliSellerID:           p.MeliSellerID,
		BusinessName:           p.BusinessName,
		Platform:               p.Platform,
		City:                   p.City,
		Phone:                  p.Phone,
		WhatsAppNumber:         p.WhatsAppNumber,
		WhatsAppURL:            buildWhatsAppURL(p.WhatsAppNumber),
		InstagramHandle:        p.InstagramHandle,
		InstagramURL:           buildInstagramURL(p.InstagramHandle),
		InstagramFollowers:     p.InstagramFollowers,
		ProductCount:           p.ProductCount,
		HasCOD:                 p.HasCOD,
		NumCouriersMentioned:   p.NumCouriersMentioned,
		UsesProbabilityChannel: p.UsesProbabilityChannel,
		DiscoveryQuery:         p.DiscoveryQuery,
		Score:                  p.Score,
		ScoreBreakdown:         p.ScoreBreakdown,
		Status:                 p.Status,
		Seen:                   p.Seen,
		SourceDiscoveredAt:     p.SourceDiscoveredAt,
		SourceLastScannedAt:    p.SourceLastScannedAt,
	}
}

func FromEntities(prospects []entities.Prospect) []ProspectResponse {
	out := make([]ProspectResponse, 0, len(prospects))
	for i := range prospects {
		out = append(out, FromEntity(&prospects[i]))
	}
	return out
}

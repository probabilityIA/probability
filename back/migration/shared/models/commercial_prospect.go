package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CommercialProspect struct {
	gorm.Model
	Domain                 string `gorm:"size:255;not null;uniqueIndex"`
	MeliSellerID           string `gorm:"size:50;index"`
	BusinessName           string `gorm:"size:255"`
	Platform               string `gorm:"size:50;index"`
	City                   string `gorm:"size:100;index"`
	Phone                  string `gorm:"size:50"`
	WhatsAppNumber         string `gorm:"size:50"`
	InstagramHandle        string `gorm:"size:150"`
	InstagramFollowers     int
	ProductCount           int
	HasCOD                 bool `gorm:"index"`
	NumCouriersMentioned   int
	UsesProbabilityChannel bool
	RobotsOK               bool
	DiscoveryQuery         string         `gorm:"size:255"`
	Score                  float64        `gorm:"index"`
	ScoreBreakdown         datatypes.JSON `gorm:"type:jsonb"`
	Status                 string         `gorm:"size:30;index"`
	SourceDiscoveredAt     *time.Time
	SourceLastScannedAt    *time.Time
}

func (CommercialProspect) TableName() string {
	return "commercial_prospects"
}

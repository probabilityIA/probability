package models

import (
	"gorm.io/gorm"
)

type GuideFormat struct {
	gorm.Model

	Carrier   string `gorm:"size:64;not null;index:idx_guide_formats_carrier"`
	Code      string `gorm:"size:64;not null;uniqueIndex:idx_guide_formats_code"`
	Label     string `gorm:"size:128;not null"`
	WidthCm   float64 `gorm:"type:decimal(6,2)"`
	HeightCm  float64 `gorm:"type:decimal(6,2)"`
	Adhesive  bool    `gorm:"default:false"`
	Strategy  string  `gorm:"size:32;not null"`

	CropLLxFrac float64 `gorm:"type:decimal(5,3);default:0"`
	CropLLyFrac float64 `gorm:"type:decimal(5,3);default:0"`
	CropURxFrac float64 `gorm:"type:decimal(5,3);default:1"`
	CropURyFrac float64 `gorm:"type:decimal(5,3);default:1"`

	SourcePage int  `gorm:"default:1"`
	IsDefault  bool `gorm:"default:false;index"`
	SortOrder  int  `gorm:"default:0"`
	IsActive   bool `gorm:"default:true;index"`
}

func (GuideFormat) TableName() string {
	return "guide_formats"
}

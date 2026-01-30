package mappers

import (
	"time"

	"gorm.io/datatypes"
)

// IntegrationNotificationConfigModel representa la tabla integration_notification_configs
type IntegrationNotificationConfigModel struct {
	ID               uint           `gorm:"primaryKey;autoIncrement"`
	IntegrationID    uint           `gorm:"not null;index:idx_inc_integration_trigger"`
	NotificationType string         `gorm:"type:varchar(50);not null"`
	IsActive         bool           `gorm:"default:true;index:idx_inc_integration_trigger"`
	Conditions       datatypes.JSON `gorm:"type:jsonb;not null"`
	Config           datatypes.JSON `gorm:"type:jsonb;not null"`
	Description      string         `gorm:"type:text"`
	Priority         int            `gorm:"default:0"`
	CreatedAt        time.Time      `gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime"`
}

// TableName especifica el nombre de la tabla
func (IntegrationNotificationConfigModel) TableName() string {
	return "integration_notification_configs"
}

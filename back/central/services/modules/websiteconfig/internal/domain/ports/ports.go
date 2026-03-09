package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/entities"
)

type IRepository interface {
	GetConfig(ctx context.Context, businessID uint) (*entities.WebsiteConfig, error)
	UpsertConfig(ctx context.Context, businessID uint, dto *dtos.UpdateConfigDTO) (*entities.WebsiteConfig, error)
}

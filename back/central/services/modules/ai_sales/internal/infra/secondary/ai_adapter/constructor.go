package ai_adapter

import (
	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/shared/bedrock"
	"github.com/secamc93/probability/back/central/shared/log"
)

type adapter struct {
	bedrock bedrock.IBedrock
	log     log.ILogger
}

// New crea un nuevo adaptador de AI que implementa domain.IAIProvider
func New(bedrockClient bedrock.IBedrock, logger log.ILogger) domain.IAIProvider {
	return &adapter{
		bedrock: bedrockClient,
		log:     logger,
	}
}

package bedrock

import (
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/ports"
	sharedbedrock "github.com/secamc93/probability/back/central/shared/bedrock"
	"github.com/secamc93/probability/back/central/shared/env"
)

const defaultModelID = "qwen.qwen3-next-80b-a3b"

type AssistantModel struct {
	client  sharedbedrock.IBedrock
	modelID string
}

func New(client sharedbedrock.IBedrock, cfg env.IConfig) ports.IAssistantModel {
	modelID := defaultModelID
	if cfg != nil {
		if configured := cfg.Get("AI_ASSISTANT_MODEL_ID"); configured != "" {
			modelID = configured
		}
	}
	return &AssistantModel{client: client, modelID: modelID}
}

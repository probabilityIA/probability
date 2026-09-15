package bedrock

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

const (
	replyToolName  = "responder"
	noDestination  = "none"
	requestTimeout = 30 * time.Second
	maxTokens      = 600
	temperature    = 0.3
)

type toolPayload struct {
	Message     string `json:"message"`
	Destination string `json:"destination"`
}

func (m *AssistantModel) Reply(ctx context.Context, req dtos.ModelRequest) (*dtos.ModelReply, error) {
	if m.client == nil || m.client.GetClient() == nil {
		return nil, fmt.Errorf("cliente de Bedrock no inicializado")
	}

	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	output, err := m.client.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId:  aws.String(m.modelID),
		System:   []types.SystemContentBlock{&types.SystemContentBlockMemberText{Value: req.SystemPrompt}},
		Messages: toBedrockMessages(req.Messages),
		InferenceConfig: &types.InferenceConfiguration{
			MaxTokens:   aws.Int32(maxTokens),
			Temperature: aws.Float32(temperature),
		},
		ToolConfig: replyToolConfig(req.DestinationKeys),
	})
	if err != nil {
		return nil, err
	}
	reply, err := parseOutput(output)
	if err != nil {
		return nil, err
	}
	reply.Model = m.modelID
	if output.Usage != nil {
		reply.InputTokens = int(aws.ToInt32(output.Usage.InputTokens))
		reply.OutputTokens = int(aws.ToInt32(output.Usage.OutputTokens))
	}
	return reply, nil
}

func toBedrockMessages(messages []entities.ChatMessage) []types.Message {
	result := make([]types.Message, 0, len(messages))
	for _, msg := range messages {
		role := types.ConversationRoleUser
		if msg.Role == entities.RoleAssistant {
			role = types.ConversationRoleAssistant
		}
		result = append(result, types.Message{
			Role:    role,
			Content: []types.ContentBlock{&types.ContentBlockMemberText{Value: msg.Text}},
		})
	}
	return result
}

func replyToolConfig(destinationKeys []string) *types.ToolConfiguration {
	enum := append([]string{noDestination}, destinationKeys...)
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"message": map[string]interface{}{
				"type":        "string",
				"description": "Respuesta breve para el usuario, sin URLs ni rutas.",
			},
			"destination": map[string]interface{}{
				"type":        "string",
				"enum":        enum,
				"description": "Clave del destino permitido que resuelve la pregunta, o none.",
			},
		},
		"required": []string{"message", "destination"},
	}

	return &types.ToolConfiguration{
		Tools: []types.Tool{
			&types.ToolMemberToolSpec{Value: types.ToolSpecification{
				Name:        aws.String(replyToolName),
				Description: aws.String("Entrega la respuesta al usuario y, si aplica, el destino de la plataforma."),
				InputSchema: &types.ToolInputSchemaMemberJson{Value: document.NewLazyDocument(schema)},
			}},
		},
		ToolChoice: &types.ToolChoiceMemberTool{Value: types.SpecificToolChoice{Name: aws.String(replyToolName)}},
	}
}

func parseOutput(output *bedrockruntime.ConverseOutput) (*dtos.ModelReply, error) {
	if output == nil || output.Output == nil {
		return nil, fmt.Errorf("respuesta vacia de Bedrock")
	}
	message, ok := output.Output.(*types.ConverseOutputMemberMessage)
	if !ok {
		return nil, fmt.Errorf("respuesta de Bedrock sin mensaje")
	}

	reply := &dtos.ModelReply{}
	var texts []string
	for _, block := range message.Value.Content {
		switch v := block.(type) {
		case *types.ContentBlockMemberText:
			texts = append(texts, v.Value)
		case *types.ContentBlockMemberToolUse:
			if aws.ToString(v.Value.Name) != replyToolName || v.Value.Input == nil {
				continue
			}
			raw, err := v.Value.Input.MarshalSmithyDocument()
			if err != nil {
				continue
			}
			var payload toolPayload
			if err := json.Unmarshal(raw, &payload); err != nil {
				continue
			}
			reply.Message = payload.Message
			reply.DestinationKey = payload.Destination
		}
	}

	if strings.TrimSpace(reply.Message) == "" {
		reply.Message = strings.Join(texts, "\n")
	}
	return reply, nil
}

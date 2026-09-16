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
	maxTokens      = 700
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

	messages, err := toBedrockMessages(req.Messages)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	output, err := m.client.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId:  aws.String(m.modelID),
		System:   []types.SystemContentBlock{&types.SystemContentBlockMemberText{Value: req.SystemPrompt}},
		Messages: messages,
		InferenceConfig: &types.InferenceConfiguration{
			MaxTokens:   aws.Int32(maxTokens),
			Temperature: aws.Float32(temperature),
		},
		ToolConfig: toolConfig(req),
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

func toBedrockMessages(messages []dtos.ModelMessage) ([]types.Message, error) {
	result := make([]types.Message, 0, len(messages))
	for _, msg := range messages {
		role := types.ConversationRoleUser
		if msg.Role == entities.RoleAssistant {
			role = types.ConversationRoleAssistant
		}

		var content []types.ContentBlock
		if strings.TrimSpace(msg.Text) != "" {
			content = append(content, &types.ContentBlockMemberText{Value: msg.Text})
		}
		for _, call := range msg.ToolCalls {
			input := call.Input
			if input == nil {
				input = map[string]any{}
			}
			content = append(content, &types.ContentBlockMemberToolUse{Value: types.ToolUseBlock{
				ToolUseId: aws.String(call.ID),
				Name:      aws.String(call.Name),
				Input:     document.NewLazyDocument(input),
			}})
		}
		for _, res := range msg.ToolResults {
			payload := res.Content
			if payload == nil {
				payload = map[string]any{}
			}
			content = append(content, &types.ContentBlockMemberToolResult{Value: types.ToolResultBlock{
				ToolUseId: aws.String(res.ToolCallID),
				Content:   []types.ToolResultContentBlock{&types.ToolResultContentBlockMemberJson{Value: document.NewLazyDocument(payload)}},
			}})
		}
		if len(content) == 0 {
			return nil, fmt.Errorf("mensaje vacio para Bedrock")
		}
		result = append(result, types.Message{Role: role, Content: content})
	}
	return result, nil
}

func toolConfig(req dtos.ModelRequest) *types.ToolConfiguration {
	tools := make([]types.Tool, 0, len(req.Tools)+1)
	for _, def := range req.Tools {
		tools = append(tools, &types.ToolMemberToolSpec{Value: types.ToolSpecification{
			Name:        aws.String(def.Name),
			Description: aws.String(def.Description),
			InputSchema: &types.ToolInputSchemaMemberJson{Value: document.NewLazyDocument(def.Schema)},
		}})
	}
	tools = append(tools, replyTool(req.DestinationKeys))

	var choice types.ToolChoice = &types.ToolChoiceMemberTool{Value: types.SpecificToolChoice{Name: aws.String(replyToolName)}}
	if len(req.Tools) > 0 && !req.ForceReply {
		choice = &types.ToolChoiceMemberAny{Value: types.AnyToolChoice{}}
	}
	return &types.ToolConfiguration{Tools: tools, ToolChoice: choice}
}

func replyTool(destinationKeys []string) types.Tool {
	enum := append([]string{noDestination}, destinationKeys...)
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"message": map[string]any{
				"type":        "string",
				"description": "Respuesta breve para el usuario, sin URLs ni rutas.",
			},
			"destination": map[string]any{
				"type":        "string",
				"enum":        enum,
				"description": "Clave del destino permitido que resuelve la pregunta, o none.",
			},
		},
		"required": []string{"message", "destination"},
	}
	return &types.ToolMemberToolSpec{Value: types.ToolSpecification{
		Name:        aws.String(replyToolName),
		Description: aws.String("Entrega la respuesta final al usuario y, si aplica, el destino de la plataforma."),
		InputSchema: &types.ToolInputSchemaMemberJson{Value: document.NewLazyDocument(schema)},
	}}
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
	var calls []dtos.ToolCall
	final := false
	for _, block := range message.Value.Content {
		switch v := block.(type) {
		case *types.ContentBlockMemberText:
			texts = append(texts, v.Value)
		case *types.ContentBlockMemberToolUse:
			input := decodeInput(v.Value.Input)
			name := aws.ToString(v.Value.Name)
			if name == replyToolName {
				raw, _ := json.Marshal(input)
				var payload toolPayload
				if err := json.Unmarshal(raw, &payload); err == nil {
					reply.Message = payload.Message
					reply.DestinationKey = payload.Destination
					final = true
				}
				continue
			}
			calls = append(calls, dtos.ToolCall{ID: aws.ToString(v.Value.ToolUseId), Name: name, Input: input})
		}
	}

	if !final {
		reply.ToolCalls = calls
	}
	if strings.TrimSpace(reply.Message) == "" {
		reply.Message = strings.Join(texts, "\n")
	}
	return reply, nil
}

func decodeInput(doc document.Interface) map[string]any {
	out := map[string]any{}
	if doc == nil {
		return out
	}
	raw, err := doc.MarshalSmithyDocument()
	if err != nil {
		return out
	}
	_ = json.Unmarshal(raw, &out)
	return out
}

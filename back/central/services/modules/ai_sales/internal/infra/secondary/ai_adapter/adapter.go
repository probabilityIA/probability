package ai_adapter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
)

func (a *adapter) Converse(ctx context.Context, messages []domain.AIMessage, systemPrompt string, tools []domain.ToolDefinition) (*domain.AIResponse, error) {
	if a.bedrock.GetClient() == nil {
		return nil, &domain.ErrBedrockUnavailable{Cause: fmt.Errorf("client not initialized")}
	}

	// Mapear mensajes de dominio a tipos de Bedrock
	bedrockMessages, err := a.mapMessages(messages)
	if err != nil {
		return nil, fmt.Errorf("error mapping messages: %w", err)
	}

	// Construir input
	input := &bedrockruntime.ConverseInput{
		ModelId: aws.String("amazon.nova-micro-v1:0"),
		System: []types.SystemContentBlock{
			&types.SystemContentBlockMemberText{Value: systemPrompt},
		},
		Messages: bedrockMessages,
		InferenceConfig: &types.InferenceConfiguration{
			MaxTokens:   aws.Int32(1024),
			Temperature: aws.Float32(0.7),
		},
	}

	// Agregar tools si hay
	if len(tools) > 0 {
		toolSpecs, err := a.mapTools(tools)
		if err != nil {
			return nil, fmt.Errorf("error mapping tools: %w", err)
		}
		input.ToolConfig = &types.ToolConfiguration{
			Tools: toolSpecs,
		}
	}

	// Llamar Bedrock
	output, err := a.bedrock.Converse(ctx, input)
	if err != nil {
		return nil, &domain.ErrBedrockUnavailable{Cause: err}
	}

	// Mapear respuesta a dominio
	return a.mapResponse(output), nil
}

func (a *adapter) mapMessages(messages []domain.AIMessage) ([]types.Message, error) {
	var result []types.Message

	for _, msg := range messages {
		var contentBlocks []types.ContentBlock

		for _, block := range msg.Content {
			switch block.Type {
			case domain.ContentTypeText:
				contentBlocks = append(contentBlocks, &types.ContentBlockMemberText{
					Value: block.Text,
				})

			case domain.ContentTypeToolUse:
				var inputDoc map[string]interface{}
				if err := json.Unmarshal([]byte(block.Input), &inputDoc); err != nil {
					return nil, fmt.Errorf("error parsing tool input: %w", err)
				}
				contentBlocks = append(contentBlocks, &types.ContentBlockMemberToolUse{
					Value: types.ToolUseBlock{
						ToolUseId: aws.String(block.ToolUseID),
						Name:      aws.String(block.ToolName),
						Input:     document.NewLazyDocument(inputDoc),
					},
				})

			case domain.ContentTypeToolResult:
				contentBlocks = append(contentBlocks, &types.ContentBlockMemberToolResult{
					Value: types.ToolResultBlock{
						ToolUseId: aws.String(block.ToolUseID),
						Content: []types.ToolResultContentBlock{
							&types.ToolResultContentBlockMemberText{
								Value: block.Content,
							},
						},
					},
				})
			}
		}

		role := types.ConversationRoleUser
		if msg.Role == domain.RoleAssistant {
			role = types.ConversationRoleAssistant
		}

		result = append(result, types.Message{
			Role:    role,
			Content: contentBlocks,
		})
	}

	return result, nil
}

func (a *adapter) mapTools(tools []domain.ToolDefinition) ([]types.Tool, error) {
	var result []types.Tool

	for _, tool := range tools {
		var schema map[string]interface{}
		if err := json.Unmarshal([]byte(tool.InputSchema), &schema); err != nil {
			return nil, fmt.Errorf("error parsing tool schema for %s: %w", tool.Name, err)
		}

		result = append(result, &types.ToolMemberToolSpec{
			Value: types.ToolSpecification{
				Name:        aws.String(tool.Name),
				Description: aws.String(tool.Description),
				InputSchema: &types.ToolInputSchemaMemberJson{
					Value: document.NewLazyDocument(schema),
				},
			},
		})
	}

	return result, nil
}

func (a *adapter) mapResponse(output *bedrockruntime.ConverseOutput) *domain.AIResponse {
	resp := &domain.AIResponse{}

	// Mapear stop reason
	switch output.StopReason {
	case types.StopReasonToolUse:
		resp.StopReason = domain.StopReasonToolUse
	case types.StopReasonMaxTokens:
		resp.StopReason = domain.StopReasonMaxToken
	default:
		resp.StopReason = domain.StopReasonEndTurn
	}

	// Mapear content blocks
	if output.Output != nil {
		if msgOutput, ok := output.Output.(*types.ConverseOutputMemberMessage); ok {
			for _, block := range msgOutput.Value.Content {
				switch v := block.(type) {
				case *types.ContentBlockMemberText:
					resp.Content = append(resp.Content, domain.ContentBlock{
						Type: domain.ContentTypeText,
						Text: v.Value,
					})

				case *types.ContentBlockMemberToolUse:
					inputBytes, _ := v.Value.Input.MarshalSmithyDocument()
					resp.Content = append(resp.Content, domain.ContentBlock{
						Type:      domain.ContentTypeToolUse,
						ToolUseID: aws.ToString(v.Value.ToolUseId),
						ToolName:  aws.ToString(v.Value.Name),
						Input:     string(inputBytes),
					})
				}
			}
		}
	}

	return resp
}

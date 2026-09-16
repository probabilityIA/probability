package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers/response"
)

func (h *handler) GetSystemFlows(c *gin.Context) {
	flows := h.useCase.SystemFlows()
	out := make([]response.SystemFlow, 0, len(flows))
	for _, flow := range flows {
		out = append(out, response.SystemFlow{
			Key:         flow.Key,
			Name:        flow.Name,
			Description: flow.Description,
			Trigger:     flow.Trigger,
			Root:        toSystemFlowNode(flow.Root),
		})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": out})
}

func toSystemFlowNode(node *entities.SystemFlowNode) *response.SystemFlowNode {
	if node == nil {
		return nil
	}
	out := &response.SystemFlowNode{
		Kind:    string(node.Kind),
		Title:   node.Title,
		Detail:  node.Detail,
		Effects: node.Effects,
	}
	if node.Template != "" {
		out.Template = node.Template
		if tpl, ok := entities.Templates[node.Template]; ok {
			out.Description = tpl.Description
			out.Body = tpl.Body
			out.Buttons = tpl.ButtonLabels
		}
	}
	for _, branch := range node.Branches {
		out.Branches = append(out.Branches, response.SystemFlowBranch{
			Label:  branch.Label,
			Node:   toSystemFlowNode(branch.Node),
			BackTo: branch.BackTo,
		})
	}
	return out
}

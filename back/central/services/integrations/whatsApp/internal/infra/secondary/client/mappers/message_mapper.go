package mappers

import (
		"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/client/request"
)

// MapDomainToRequest transforma el mensaje de dominio al DTO de request de WhatsApp
func MapDomainToRequest(d entities.TemplateMessage) any {
	if d.Type == "text" {
		return request.TextMessageRequest{
			MessagingProduct: d.MessagingProduct,
			To:               d.To,
			Type:             d.Type,
			Text:             request.TextBody{Body: d.TextBody},
		}
	}

	// Mapear componentes si existen
	var components []request.Component
	for _, comp := range d.Template.Components {
		var parameters []request.Parameter
		for _, param := range comp.Parameters {
			parameters = append(parameters, request.Parameter{
				Type:          param.Type,
				ParameterName: param.ParameterName,
				Text:          param.Text,
			})
		}
		components = append(components, request.Component{
			Type:       comp.Type,
			SubType:    comp.SubType,
			Index:      comp.Index,
			Parameters: parameters,
		})
	}

	return request.TemplateMessageRequest{
		MessagingProduct: d.MessagingProduct,
		RecipientType:    d.RecipientType,
		To:               d.To,
		Type:             d.Type,
		Template: request.Template{
			Name:       d.Template.Name,
			Language:   request.Language{Code: d.Template.Language.Code},
			Components: components,
		},
	}
}

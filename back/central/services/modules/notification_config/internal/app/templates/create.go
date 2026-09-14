package templates

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

const (
	maxTemplateNameLength = 512
	maxBodyLength         = 1024
	maxHeaderLength       = 60
	maxFooterLength       = 60
)

var (
	templateNamePattern = regexp.MustCompile(`^[a-z0-9_]+$`)
	placeholderPattern  = regexp.MustCompile(`\{\{(\d+)\}\}`)
)

func (uc *useCase) Create(ctx context.Context, dto dtos.CreateTemplateDTO) (*entities.WhatsappTemplate, error) {
	template, err := buildTemplate(dto)
	if err != nil {
		return nil, err
	}

	existing, err := uc.repository.GetTemplateByName(ctx, dto.BusinessID, template.Name, template.Language)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("ya existe una plantilla llamada %s en %s para este negocio", template.Name, template.Language)
	}

	if err := uc.repository.CreateTemplate(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

func (uc *useCase) sendToMeta(ctx context.Context, template *entities.WhatsappTemplate) error {
	template.Components = BuildMetaComponents(template)
	if template.MetaTemplateID != "" {
		return uc.submitEdit(ctx, template)
	}
	return uc.submit(ctx, template)
}

func (uc *useCase) submitEdit(ctx context.Context, template *entities.WhatsappTemplate) error {
	if uc.publisher == nil {
		uc.logger.Warn().Uint("template_id", template.ID).
			Msg("Sin publicador de plantillas: el cambio queda en borrador y no se envia a Meta")
		return nil
	}

	businessID := uint(0)
	if template.BusinessID != nil {
		businessID = *template.BusinessID
	}

	message := dtos.TemplateSubmissionMessage{
		Action:         "update",
		TemplateID:     template.ID,
		BusinessID:     businessID,
		Name:           template.Name,
		Language:       template.Language,
		Category:       template.Category,
		MetaTemplateID: template.MetaTemplateID,
		HeaderMediaURL: template.HeaderMediaURL,
		Components:     template.Components,
	}

	if err := uc.publisher.PublishTemplateSubmission(ctx, message); err != nil {
		uc.logger.Error().Err(err).Uint("template_id", template.ID).
			Msg("Error publicando la edicion de la plantilla a la cola de Meta")
		return err
	}

	now := time.Now()
	template.Status = entities.TemplateStatusPending
	template.SubmittedAt = &now

	if err := uc.repository.UpdateTemplate(ctx, template); err != nil {
		return err
	}

	uc.logger.Info().Uint("template_id", template.ID).Str("name", template.Name).
		Str("meta_template_id", template.MetaTemplateID).
		Msg("Edicion de plantilla encolada para Meta")

	return nil
}

func (uc *useCase) submit(ctx context.Context, template *entities.WhatsappTemplate) error {
	if uc.publisher == nil {
		uc.logger.Warn().Uint("template_id", template.ID).
			Msg("Sin publicador de plantillas: la plantilla queda en borrador y no se envia a Meta")
		return nil
	}

	businessID := uint(0)
	if template.BusinessID != nil {
		businessID = *template.BusinessID
	}

	message := dtos.TemplateSubmissionMessage{
		Action:         "create",
		TemplateID:     template.ID,
		BusinessID:     businessID,
		Name:           template.Name,
		Language:       template.Language,
		Category:       template.Category,
		HeaderMediaURL: template.HeaderMediaURL,
		Components:     template.Components,
	}

	if err := uc.publisher.PublishTemplateSubmission(ctx, message); err != nil {
		uc.logger.Error().Err(err).Uint("template_id", template.ID).
			Msg("Error publicando la plantilla a la cola de envio a Meta")
		return err
	}

	now := time.Now()
	template.Status = entities.TemplateStatusPending
	template.SubmittedAt = &now

	if err := uc.repository.UpdateTemplate(ctx, template); err != nil {
		return err
	}

	uc.logger.Info().Uint("template_id", template.ID).Str("name", template.Name).
		Msg("Plantilla encolada para envio a Meta")

	return nil
}

func buildTemplate(dto dtos.CreateTemplateDTO) (*entities.WhatsappTemplate, error) {
	name := strings.ToLower(strings.TrimSpace(dto.Name))
	if name == "" {
		return nil, fmt.Errorf("el nombre de la plantilla es obligatorio")
	}
	if len(name) > maxTemplateNameLength {
		return nil, fmt.Errorf("el nombre de la plantilla supera %d caracteres", maxTemplateNameLength)
	}
	if !templateNamePattern.MatchString(name) {
		return nil, fmt.Errorf("el nombre solo admite minusculas, numeros y guion bajo: %s", name)
	}

	if dto.BusinessID > 0 {
		prefix := fmt.Sprintf("%d_", dto.BusinessID)
		if !strings.HasPrefix(name, prefix) {
			name = prefix + name
		}
		if len(name) > maxTemplateNameLength {
			return nil, fmt.Errorf("el nombre de la plantilla supera %d caracteres", maxTemplateNameLength)
		}
	}

	body := strings.TrimSpace(dto.BodyText)
	if body == "" {
		return nil, fmt.Errorf("el cuerpo de la plantilla es obligatorio")
	}
	if len(body) > maxBodyLength {
		return nil, fmt.Errorf("el cuerpo supera %d caracteres", maxBodyLength)
	}

	header := strings.TrimSpace(dto.HeaderText)
	if len(header) > maxHeaderLength {
		return nil, fmt.Errorf("el encabezado supera %d caracteres", maxHeaderLength)
	}

	headerType := strings.ToUpper(strings.TrimSpace(dto.HeaderType))
	if headerType == "" {
		headerType = entities.TemplateHeaderTypeText
	}
	if headerType != entities.TemplateHeaderTypeText && headerType != entities.TemplateHeaderTypeImage {
		return nil, fmt.Errorf("tipo de encabezado no soportado: %s", headerType)
	}

	headerMediaURL := strings.TrimSpace(dto.HeaderMediaURL)
	if headerType == entities.TemplateHeaderTypeImage {
		if headerMediaURL == "" {
			return nil, fmt.Errorf("el encabezado de imagen necesita una imagen cargada")
		}
		header = ""
	} else {
		headerMediaURL = ""
	}

	footer := strings.TrimSpace(dto.FooterText)
	if len(footer) > maxFooterLength {
		return nil, fmt.Errorf("el pie supera %d caracteres", maxFooterLength)
	}

	category := strings.ToUpper(strings.TrimSpace(dto.Category))
	if category == "" {
		category = entities.TemplateCategoryMarketing
	}
	if category != entities.TemplateCategoryMarketing && category != entities.TemplateCategoryUtility {
		return nil, fmt.Errorf("categoria no soportada: %s", category)
	}

	language := strings.TrimSpace(dto.Language)
	if language == "" {
		language = "es"
	}

	variables, err := buildVariables(dto.Variables, body)
	if err != nil {
		return nil, err
	}

	buttons, err := buildButtonList(dto.Buttons, category)
	if err != nil {
		return nil, err
	}

	businessID := dto.BusinessID

	scope := strings.TrimSpace(dto.Scope)
	if !entities.IsAllowedScope(scope) {
		scope = entities.TemplateScopeScheduled
	}

	template := &entities.WhatsappTemplate{
		BusinessID:     &businessID,
		Origin:         entities.TemplateOriginBusiness,
		Scope:          scope,
		Name:           name,
		Language:       language,
		Category:       category,
		BodyText:       body,
		HeaderText:     header,
		HeaderType:     headerType,
		HeaderMediaURL: headerMediaURL,
		FooterText:     footer,
		Variables:      variables,
		Buttons:        buttons,
		Status:         entities.TemplateStatusDraft,
		CreatedByID:    dto.CreatedBy,
	}

	template.Components = BuildMetaComponents(template)

	return template, nil
}

func buildButtonList(input []dtos.TemplateButtonDTO, category string) ([]entities.TemplateButton, error) {
	limit := entities.MaxUserButtons(category)
	if len(input) > limit {
		if category == entities.TemplateCategoryMarketing {
			return nil, fmt.Errorf("en marketing solo caben %d botones: el de baja ocupa el tercero", limit)
		}
		return nil, fmt.Errorf("una plantilla admite maximo %d botones", limit)
	}

	buttons := make([]entities.TemplateButton, 0, len(input))
	seen := make(map[string]bool, len(input))

	for _, button := range input {
		text := strings.TrimSpace(button.Text)
		if text == "" {
			return nil, fmt.Errorf("un boton no puede ir sin texto")
		}
		if len(text) > entities.MaxTemplateButtonTextLen {
			return nil, fmt.Errorf("el boton %q supera %d caracteres", text, entities.MaxTemplateButtonTextLen)
		}

		key := strings.ToLower(text)
		if seen[key] {
			return nil, fmt.Errorf("hay dos botones con el texto %q", text)
		}
		if strings.EqualFold(text, OptOutButtonText) {
			return nil, fmt.Errorf("%q lo agrega Meta solo en las plantillas de marketing", OptOutButtonText)
		}
		seen[key] = true

		buttonType := strings.ToUpper(strings.TrimSpace(button.Type))
		if buttonType == "" {
			buttonType = entities.TemplateButtonTypeQuickReply
		}

		url := strings.TrimSpace(button.URL)
		switch buttonType {
		case entities.TemplateButtonTypeQuickReply:
			url = ""
		case entities.TemplateButtonTypeURL:
			if url == "" {
				return nil, fmt.Errorf("el boton %q es de enlace y no tiene URL", text)
			}
		default:
			return nil, fmt.Errorf("tipo de boton no soportado: %s", buttonType)
		}

		buttons = append(buttons, entities.TemplateButton{
			Type: buttonType,
			Text: text,
			URL:  url,
		})
	}

	return buttons, nil
}

func buildVariables(input []dtos.TemplateVariableDTO, body string) ([]entities.TemplateVariable, error) {
	placeholders := placeholderPattern.FindAllStringSubmatch(body, -1)
	unique := make(map[string]bool, len(placeholders))
	for _, match := range placeholders {
		unique[match[1]] = true
	}

	if len(unique) != len(input) {
		return nil, fmt.Errorf("el cuerpo usa %d variables y se declararon %d", len(unique), len(input))
	}

	variables := make([]entities.TemplateVariable, 0, len(input))
	seen := make(map[int]bool, len(input))

	for _, item := range input {
		if item.Position <= 0 {
			return nil, fmt.Errorf("la posicion de una variable debe ser mayor a cero")
		}
		if seen[item.Position] {
			return nil, fmt.Errorf("la posicion %d esta declarada dos veces", item.Position)
		}
		seen[item.Position] = true

		source := strings.TrimSpace(item.Source)
		if !entities.IsAllowedVariableSource(source) {
			return nil, fmt.Errorf("variable no permitida: %s", source)
		}

		label := strings.TrimSpace(item.Label)
		if label == "" {
			label = entities.VariableSourceLabel(source)
		}

		variables = append(variables, entities.TemplateVariable{
			Position: item.Position,
			Source:   source,
			Label:    label,
			Fallback: strings.TrimSpace(item.Fallback),
		})
	}

	for position := 1; position <= len(variables); position++ {
		if !seen[position] {
			return nil, fmt.Errorf("falta declarar la variable {{%d}}", position)
		}
	}

	sorted := make([]entities.TemplateVariable, len(variables))
	for _, variable := range variables {
		sorted[variable.Position-1] = variable
	}

	return sorted, nil
}

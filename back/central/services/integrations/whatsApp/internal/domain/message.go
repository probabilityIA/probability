package domain

// TemplateLanguage representa el idioma del template (dominio)
type TemplateLanguage struct {
	Code string
}

// TemplateParameter representa un parámetro del template
type TemplateParameter struct {
	Type          string
	ParameterName string
	Text          string
}

// TemplateComponent representa un componente del template
type TemplateComponent struct {
	Type       string
	Parameters []TemplateParameter
}

// TemplateData representa los datos del template (dominio)
type TemplateData struct {
	Name       string
	Language   TemplateLanguage
	Components []TemplateComponent
}

// TemplateMessage representa el mensaje a enviar (dominio)
type TemplateMessage struct {
	MessagingProduct string
	RecipientType    string
	To               string
	Type             string
	Template         TemplateData
	TextBody         string
}

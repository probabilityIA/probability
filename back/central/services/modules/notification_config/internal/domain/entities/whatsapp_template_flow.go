package entities

const MaxFlowDepth = 4

type TemplateFlow struct {
	ID               uint
	BusinessID       uint
	FlowID           *uint
	SourceTemplateID uint
	ButtonText       string
	TargetTemplateID uint
	Enabled          bool

	SourceName           string
	TargetName           string
	TargetLanguage       string
	TargetStatus         string
	TargetHeaderMediaURL string
}

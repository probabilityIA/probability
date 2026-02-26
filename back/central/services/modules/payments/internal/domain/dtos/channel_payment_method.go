package dtos

type ChannelPaymentMethodInfo struct {
	ID              uint
	IntegrationType string
	Code            string
	Name            string
	Description     string
	IsActive        bool
	DisplayOrder    int
}

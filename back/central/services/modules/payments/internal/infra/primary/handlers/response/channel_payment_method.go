package response

type ChannelPaymentMethodResponse struct {
	ID              uint   `json:"id"`
	IntegrationType string `json:"integration_type"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	IsActive        bool   `json:"is_active"`
	DisplayOrder    int    `json:"display_order"`
}

type ChannelPaymentMethodListResponse struct {
	Success bool                            `json:"success"`
	Message string                          `json:"message"`
	Data    []ChannelPaymentMethodResponse  `json:"data"`
}

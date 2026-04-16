package dtos

// SendMessageRequest representa la solicitud simplificada para enviar un mensaje de WhatsApp
type SendMessageRequest struct {
	OrderNumber string // Número de orden
	PhoneNumber string // Número de celular al que se va a enviar (formato internacional: +573001234567)
}

package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/queue/messages"
)

// BulkJobDTOToMessage convierte un DTO de dominio a un mensaje de queue
func BulkJobDTOToMessage(dto *dtos.BulkInvoiceJobMessage) *messages.BulkInvoiceJobMessage {
	return &messages.BulkInvoiceJobMessage{
		JobID:         dto.JobID,
		OrderID:       dto.OrderID,
		BusinessID:    dto.BusinessID,
		IsManual:      dto.IsManual,
		CreatedBy:     dto.CreatedBy,
		AttemptNumber: dto.AttemptNumber,
	}
}

// BulkJobMessageToDTO convierte un mensaje de queue a un DTO de dominio
func BulkJobMessageToDTO(msg *messages.BulkInvoiceJobMessage) *dtos.BulkInvoiceJobMessage {
	return &dtos.BulkInvoiceJobMessage{
		JobID:         msg.JobID,
		OrderID:       msg.OrderID,
		BusinessID:    msg.BusinessID,
		IsManual:      msg.IsManual,
		CreatedBy:     msg.CreatedBy,
		AttemptNumber: msg.AttemptNumber,
	}
}

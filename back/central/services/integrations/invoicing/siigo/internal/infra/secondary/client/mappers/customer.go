package mappers

import (
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/response"
)

// CustomerToDTO convierte un Customer de Siigo a CustomerResult de dominio
func CustomerToDTO(customer *response.Customer) *dtos.CustomerResult {
	if customer == nil {
		return nil
	}

	result := &dtos.CustomerResult{
		ID:             customer.ID,
		Identification: customer.Identification,
	}

	// Combinar nombres
	if len(customer.Name) > 0 {
		result.Name = strings.Join(customer.Name, " ")
	}

	// Extraer email del primer contacto
	if len(customer.Emails) > 0 {
		result.Email = customer.Emails[0].Email
	}

	// Extraer teléfono del primer teléfono
	if len(customer.Phones) > 0 {
		result.Phone = customer.Phones[0].Number
	}

	// Extraer dirección
	if customer.Address != nil {
		result.Address = customer.Address.Address
	}

	return result
}

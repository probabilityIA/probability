package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/request"
)

// RequestToCreateOrderDTO maps a create order request to the domain DTO
func RequestToCreateOrderDTO(req *request.CreateOrderRequest) *dtos.StorefrontCreateOrderDTO {
	dto := &dtos.StorefrontCreateOrderDTO{
		Notes: req.Notes,
	}

	for _, item := range req.Items {
		dto.Items = append(dto.Items, dtos.StorefrontOrderItemDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	if req.Address != nil {
		dto.Address = &dtos.StorefrontAddressDTO{
			FirstName:    req.Address.FirstName,
			LastName:     req.Address.LastName,
			Phone:        req.Address.Phone,
			Street:       req.Address.Street,
			Street2:      req.Address.Street2,
			City:         req.Address.City,
			State:        req.Address.State,
			Country:      req.Address.Country,
			PostalCode:   req.Address.PostalCode,
			Instructions: req.Address.Instructions,
		}
	}

	return dto
}

// RequestToRegisterDTO maps a register request to the domain DTO
func RequestToRegisterDTO(req *request.RegisterRequest) *dtos.RegisterDTO {
	return &dtos.RegisterDTO{
		Name:         req.Name,
		Email:        req.Email,
		Password:     req.Password,
		Phone:        req.Phone,
		Dni:          req.Dni,
		BusinessCode: req.BusinessCode,
	}
}

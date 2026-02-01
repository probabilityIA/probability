package usecases

import (
	"context"
	"math"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/app/usecases/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/errors"
)

// ═══════════════════════════════════════════
// PAYMENT METHODS USE CASES
// ═══════════════════════════════════════════

// ListPaymentMethods obtiene una lista paginada de métodos de pago
func (uc *UseCase) ListPaymentMethods(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*dtos.PaymentMethodsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	methods, total, err := uc.repo.ListPaymentMethods(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	// Mapear a response
	data := mappers.EntitiesToResponses(methods)

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dtos.PaymentMethodsListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetPaymentMethodByID obtiene un método de pago por ID
func (uc *UseCase) GetPaymentMethodByID(ctx context.Context, id uint) (*dtos.PaymentMethodResponse, error) {
	method, err := uc.repo.GetPaymentMethodByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mappers.EntityToResponse(method)
	return &response, nil
}

// GetPaymentMethodByCode obtiene un método de pago por código
func (uc *UseCase) GetPaymentMethodByCode(ctx context.Context, code string) (*dtos.PaymentMethodResponse, error) {
	method, err := uc.repo.GetPaymentMethodByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	response := mappers.EntityToResponse(method)
	return &response, nil
}

// CreatePaymentMethod crea un nuevo método de pago
func (uc *UseCase) CreatePaymentMethod(ctx context.Context, req *dtos.CreatePaymentMethod) (*dtos.PaymentMethodResponse, error) {
	// Validar que el código no exista
	exists, err := uc.repo.PaymentMethodExists(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrPaymentMethodCodeAlreadyExists
	}

	// Convertir DTO a entidad
	method := mappers.CreateDTOToEntity(req)

	if err := uc.repo.CreatePaymentMethod(ctx, method); err != nil {
		return nil, err
	}

	response := mappers.EntityToResponse(method)
	return &response, nil
}

// UpdatePaymentMethod actualiza un método de pago existente
func (uc *UseCase) UpdatePaymentMethod(ctx context.Context, id uint, req *dtos.UpdatePaymentMethod) (*dtos.PaymentMethodResponse, error) {
	// Obtener método existente
	method, err := uc.repo.GetPaymentMethodByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	method.Name = req.Name
	method.Description = req.Description
	method.Category = req.Category
	method.Provider = req.Provider
	method.Icon = req.Icon
	method.Color = req.Color

	if err := uc.repo.UpdatePaymentMethod(ctx, method); err != nil {
		return nil, err
	}

	response := mappers.EntityToResponse(method)
	return &response, nil
}

// DeletePaymentMethod elimina un método de pago
func (uc *UseCase) DeletePaymentMethod(ctx context.Context, id uint) error {
	// Verificar que no tenga mapeos activos
	hasActive, err := uc.repo.PaymentMethodHasActiveMappings(ctx, id)
	if err != nil {
		return err
	}
	if hasActive {
		return errors.ErrPaymentMethodHasActiveMappings
	}

	return uc.repo.DeletePaymentMethod(ctx, id)
}

// TogglePaymentMethodActive activa/desactiva un método de pago
func (uc *UseCase) TogglePaymentMethodActive(ctx context.Context, id uint) (*dtos.PaymentMethodResponse, error) {
	method, err := uc.repo.TogglePaymentMethodActive(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mappers.EntityToResponse(method)
	return &response, nil
}

// ═══════════════════════════════════════════
// PAYMENT MAPPINGS USE CASES
// ═══════════════════════════════════════════

// ListPaymentMappings obtiene una lista de mapeos
func (uc *UseCase) ListPaymentMappings(ctx context.Context, filters map[string]interface{}) (*dtos.PaymentMappingsListResponse, error) {
	mappings, total, err := uc.repo.ListPaymentMappingsWithMethods(ctx, filters)
	if err != nil {
		return nil, err
	}

	data := mappers.MappingEntitiesToResponses(mappings)

	return &dtos.PaymentMappingsListResponse{
		Data:  data,
		Total: total,
	}, nil
}

// GetPaymentMappingByID obtiene un mapeo por ID
func (uc *UseCase) GetPaymentMappingByID(ctx context.Context, id uint) (*dtos.PaymentMappingResponse, error) {
	mapping, err := uc.repo.GetPaymentMappingByIDWithMethod(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mappers.MappingEntityToResponse(mapping)
	return &response, nil
}

// GetPaymentMappingsByIntegrationType obtiene mapeos por tipo de integración
func (uc *UseCase) GetPaymentMappingsByIntegrationType(ctx context.Context, integrationType string) ([]dtos.PaymentMappingResponse, error) {
	mappings, err := uc.repo.GetPaymentMappingsByIntegrationTypeWithMethods(ctx, integrationType)
	if err != nil {
		return nil, err
	}

	responses := mappers.MappingEntitiesToResponses(mappings)
	return responses, nil
}

// GetAllPaymentMappingsGroupedByIntegration obtiene todos los mapeos agrupados por tipo de integración
func (uc *UseCase) GetAllPaymentMappingsGroupedByIntegration(ctx context.Context) ([]dtos.PaymentMappingsByIntegrationResponse, error) {
	mappings, _, err := uc.repo.ListPaymentMappingsWithMethods(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Agrupar por tipo de integración
	grouped := make(map[string][]dtos.PaymentMappingResponse)
	for _, mapping := range mappings {
		response := mappers.MappingEntityToResponse(&mapping)
		grouped[mapping.IntegrationType] = append(grouped[mapping.IntegrationType], response)
	}

	// Convertir a slice
	result := make([]dtos.PaymentMappingsByIntegrationResponse, 0, len(grouped))
	for integrationType, mappings := range grouped {
		result = append(result, dtos.PaymentMappingsByIntegrationResponse{
			IntegrationType: integrationType,
			Mappings:        mappings,
		})
	}

	return result, nil
}

// CreatePaymentMapping crea un nuevo mapeo
func (uc *UseCase) CreatePaymentMapping(ctx context.Context, req *dtos.CreatePaymentMapping) (*dtos.PaymentMappingResponse, error) {
	// Validar que no exista
	exists, err := uc.repo.PaymentMappingExists(ctx, req.IntegrationType, req.OriginalMethod)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrPaymentMappingAlreadyExists
	}

	// Validar que el método de pago exista
	_, err = uc.repo.GetPaymentMethodByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, errors.ErrPaymentMethodNotFound
	}

	// Convertir DTO a entidad
	mapping := mappers.CreateMappingDTOToEntity(req)

	if err := uc.repo.CreatePaymentMapping(ctx, mapping); err != nil {
		return nil, err
	}

	// Obtener con método de pago
	created, err := uc.repo.GetPaymentMappingByIDWithMethod(ctx, mapping.ID)
	if err != nil {
		return nil, err
	}

	response := mappers.MappingEntityToResponse(created)
	return &response, nil
}

// UpdatePaymentMapping actualiza un mapeo existente
func (uc *UseCase) UpdatePaymentMapping(ctx context.Context, id uint, req *dtos.UpdatePaymentMapping) (*dtos.PaymentMappingResponse, error) {
	// Obtener mapeo existente
	mapping, err := uc.repo.GetPaymentMappingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validar que el método de pago exista
	_, err = uc.repo.GetPaymentMethodByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, errors.ErrPaymentMethodNotFound
	}

	// Actualizar campos
	mapping.OriginalMethod = req.OriginalMethod
	mapping.PaymentMethodID = req.PaymentMethodID
	mapping.Priority = req.Priority

	if err := uc.repo.UpdatePaymentMapping(ctx, mapping); err != nil {
		return nil, err
	}

	// Obtener con método de pago
	updated, err := uc.repo.GetPaymentMappingByIDWithMethod(ctx, mapping.ID)
	if err != nil {
		return nil, err
	}

	response := mappers.MappingEntityToResponse(updated)
	return &response, nil
}

// DeletePaymentMapping elimina un mapeo
func (uc *UseCase) DeletePaymentMapping(ctx context.Context, id uint) error {
	return uc.repo.DeletePaymentMapping(ctx, id)
}

// TogglePaymentMappingActive activa/desactiva un mapeo
func (uc *UseCase) TogglePaymentMappingActive(ctx context.Context, id uint) (*dtos.PaymentMappingResponse, error) {
	mapping, err := uc.repo.TogglePaymentMappingActive(ctx, id)
	if err != nil {
		return nil, err
	}

	// Obtener con método de pago
	updated, err := uc.repo.GetPaymentMappingByIDWithMethod(ctx, mapping.ID)
	if err != nil {
		return nil, err
	}

	response := mappers.MappingEntityToResponse(updated)
	return &response, nil
}

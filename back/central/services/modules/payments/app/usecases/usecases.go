package usecases

import (
	"context"
	"errors"
	"math"

	"github.com/secamc93/probability/back/central/services/modules/payments/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ═══════════════════════════════════════════
// PAYMENT METHODS USE CASES
// ═══════════════════════════════════════════

// ListPaymentMethods obtiene una lista paginada de métodos de pago
func (uc *UseCase) ListPaymentMethods(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*domain.PaymentMethodsListResponse, error) {
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
	data := make([]domain.PaymentMethodResponse, len(methods))
	for i, method := range methods {
		data[i] = mapToPaymentMethodResponse(&method)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &domain.PaymentMethodsListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetPaymentMethodByID obtiene un método de pago por ID
func (uc *UseCase) GetPaymentMethodByID(ctx context.Context, id uint) (*domain.PaymentMethodResponse, error) {
	method, err := uc.repo.GetPaymentMethodByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mapToPaymentMethodResponse(method)
	return &response, nil
}

// GetPaymentMethodByCode obtiene un método de pago por código
func (uc *UseCase) GetPaymentMethodByCode(ctx context.Context, code string) (*domain.PaymentMethodResponse, error) {
	method, err := uc.repo.GetPaymentMethodByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	response := mapToPaymentMethodResponse(method)
	return &response, nil
}

// CreatePaymentMethod crea un nuevo método de pago
func (uc *UseCase) CreatePaymentMethod(ctx context.Context, req *domain.CreatePaymentMethodRequest) (*domain.PaymentMethodResponse, error) {
	// Validar que el código no exista
	exists, err := uc.repo.PaymentMethodExists(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("payment method with this code already exists")
	}

	// Crear modelo
	method := &models.PaymentMethod{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Provider:    req.Provider,
		Icon:        req.Icon,
		Color:       req.Color,
		IsActive:    true,
	}

	if err := uc.repo.CreatePaymentMethod(ctx, method); err != nil {
		return nil, err
	}

	response := mapToPaymentMethodResponse(method)
	return &response, nil
}

// UpdatePaymentMethod actualiza un método de pago existente
func (uc *UseCase) UpdatePaymentMethod(ctx context.Context, id uint, req *domain.UpdatePaymentMethodRequest) (*domain.PaymentMethodResponse, error) {
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

	response := mapToPaymentMethodResponse(method)
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
		return errors.New("cannot delete payment method with active mappings")
	}

	return uc.repo.DeletePaymentMethod(ctx, id)
}

// TogglePaymentMethodActive activa/desactiva un método de pago
func (uc *UseCase) TogglePaymentMethodActive(ctx context.Context, id uint) (*domain.PaymentMethodResponse, error) {
	method, err := uc.repo.TogglePaymentMethodActive(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mapToPaymentMethodResponse(method)
	return &response, nil
}

// ═══════════════════════════════════════════
// PAYMENT MAPPINGS USE CASES
// ═══════════════════════════════════════════

// ListPaymentMappings obtiene una lista de mapeos
func (uc *UseCase) ListPaymentMappings(ctx context.Context, filters map[string]interface{}) (*domain.PaymentMappingsListResponse, error) {
	mappings, total, err := uc.repo.ListPaymentMappingsWithMethods(ctx, filters)
	if err != nil {
		return nil, err
	}

	data := make([]domain.PaymentMappingResponse, len(mappings))
	for i, mapping := range mappings {
		data[i] = mapToPaymentMappingResponse(&mapping)
	}

	return &domain.PaymentMappingsListResponse{
		Data:  data,
		Total: total,
	}, nil
}

// GetPaymentMappingByID obtiene un mapeo por ID
func (uc *UseCase) GetPaymentMappingByID(ctx context.Context, id uint) (*domain.PaymentMappingResponse, error) {
	mapping, err := uc.repo.GetPaymentMappingByIDWithMethod(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mapToPaymentMappingResponse(mapping)
	return &response, nil
}

// GetPaymentMappingsByIntegrationType obtiene mapeos por tipo de integración
func (uc *UseCase) GetPaymentMappingsByIntegrationType(ctx context.Context, integrationType string) ([]domain.PaymentMappingResponse, error) {
	mappings, err := uc.repo.GetPaymentMappingsByIntegrationTypeWithMethods(ctx, integrationType)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.PaymentMappingResponse, len(mappings))
	for i, mapping := range mappings {
		responses[i] = mapToPaymentMappingResponse(&mapping)
	}

	return responses, nil
}

// GetAllPaymentMappingsGroupedByIntegration obtiene todos los mapeos agrupados por tipo de integración
func (uc *UseCase) GetAllPaymentMappingsGroupedByIntegration(ctx context.Context) ([]domain.PaymentMappingsByIntegrationResponse, error) {
	mappings, _, err := uc.repo.ListPaymentMappingsWithMethods(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Agrupar por tipo de integración
	grouped := make(map[string][]domain.PaymentMappingResponse)
	for _, mapping := range mappings {
		response := mapToPaymentMappingResponse(&mapping)
		grouped[mapping.IntegrationType] = append(grouped[mapping.IntegrationType], response)
	}

	// Convertir a slice
	result := make([]domain.PaymentMappingsByIntegrationResponse, 0, len(grouped))
	for integrationType, mappings := range grouped {
		result = append(result, domain.PaymentMappingsByIntegrationResponse{
			IntegrationType: integrationType,
			Mappings:        mappings,
		})
	}

	return result, nil
}

// CreatePaymentMapping crea un nuevo mapeo
func (uc *UseCase) CreatePaymentMapping(ctx context.Context, req *domain.CreatePaymentMappingRequest) (*domain.PaymentMappingResponse, error) {
	// Validar que no exista
	exists, err := uc.repo.PaymentMappingExists(ctx, req.IntegrationType, req.OriginalMethod)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("mapping already exists for this integration type and original method")
	}

	// Validar que el método de pago exista
	_, err = uc.repo.GetPaymentMethodByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, errors.New("payment method not found")
	}

	// Crear mapeo
	mapping := &models.PaymentMethodMapping{
		IntegrationType: req.IntegrationType,
		OriginalMethod:  req.OriginalMethod,
		PaymentMethodID: req.PaymentMethodID,
		Priority:        req.Priority,
		IsActive:        true,
	}

	if err := uc.repo.CreatePaymentMapping(ctx, mapping); err != nil {
		return nil, err
	}

	// Obtener con método de pago
	created, err := uc.repo.GetPaymentMappingByIDWithMethod(ctx, mapping.ID)
	if err != nil {
		return nil, err
	}

	response := mapToPaymentMappingResponse(created)
	return &response, nil
}

// UpdatePaymentMapping actualiza un mapeo existente
func (uc *UseCase) UpdatePaymentMapping(ctx context.Context, id uint, req *domain.UpdatePaymentMappingRequest) (*domain.PaymentMappingResponse, error) {
	// Obtener mapeo existente
	mapping, err := uc.repo.GetPaymentMappingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validar que el método de pago exista
	_, err = uc.repo.GetPaymentMethodByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, errors.New("payment method not found")
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

	response := mapToPaymentMappingResponse(updated)
	return &response, nil
}

// DeletePaymentMapping elimina un mapeo
func (uc *UseCase) DeletePaymentMapping(ctx context.Context, id uint) error {
	return uc.repo.DeletePaymentMapping(ctx, id)
}

// TogglePaymentMappingActive activa/desactiva un mapeo
func (uc *UseCase) TogglePaymentMappingActive(ctx context.Context, id uint) (*domain.PaymentMappingResponse, error) {
	mapping, err := uc.repo.TogglePaymentMappingActive(ctx, id)
	if err != nil {
		return nil, err
	}

	// Obtener con método de pago
	updated, err := uc.repo.GetPaymentMappingByIDWithMethod(ctx, mapping.ID)
	if err != nil {
		return nil, err
	}

	response := mapToPaymentMappingResponse(updated)
	return &response, nil
}

// ═══════════════════════════════════════════
// MAPPERS
// ═══════════════════════════════════════════

func mapToPaymentMethodResponse(method *models.PaymentMethod) domain.PaymentMethodResponse {
	return domain.PaymentMethodResponse{
		ID:          method.ID,
		Code:        method.Code,
		Name:        method.Name,
		Description: method.Description,
		Category:    method.Category,
		Provider:    method.Provider,
		IsActive:    method.IsActive,
		Icon:        method.Icon,
		Color:       method.Color,
		CreatedAt:   method.CreatedAt,
		UpdatedAt:   method.UpdatedAt,
	}
}

func mapToPaymentMappingResponse(mapping *models.PaymentMethodMapping) domain.PaymentMappingResponse {
	return domain.PaymentMappingResponse{
		ID:              mapping.ID,
		IntegrationType: mapping.IntegrationType,
		OriginalMethod:  mapping.OriginalMethod,
		PaymentMethodID: mapping.PaymentMethodID,
		PaymentMethod: domain.PaymentMethodResponse{
			ID:          mapping.PaymentMethod.ID,
			Code:        mapping.PaymentMethod.Code,
			Name:        mapping.PaymentMethod.Name,
			Description: mapping.PaymentMethod.Description,
			Category:    mapping.PaymentMethod.Category,
			Provider:    mapping.PaymentMethod.Provider,
			IsActive:    mapping.PaymentMethod.IsActive,
			Icon:        mapping.PaymentMethod.Icon,
			Color:       mapping.PaymentMethod.Color,
			CreatedAt:   mapping.PaymentMethod.CreatedAt,
			UpdatedAt:   mapping.PaymentMethod.UpdatedAt,
		},
		IsActive:  mapping.IsActive,
		Priority:  mapping.Priority,
		CreatedAt: mapping.CreatedAt,
		UpdatedAt: mapping.UpdatedAt,
	}
}

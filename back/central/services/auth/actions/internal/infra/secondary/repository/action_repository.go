package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"

	"gorm.io/gorm"
)

// GetActions obtiene todos los actions con filtros y paginación
func (r *Repository) GetActions(ctx context.Context, page, pageSize int, name string) ([]domain.Action, int64, error) {
	r.logger.Info().Int("page", page).Int("page_size", pageSize).Str("name", name).Msg("Iniciando búsqueda de actions")

	// Configurar paginación por defecto
	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pageSize

	// Construir query base
	query := r.database.Conn(ctx).Model(&models.Action{})

	// Aplicar filtros
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	// Contar total antes de aplicar paginación
	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al contar actions")
		return nil, 0, err
	}

	// Aplicar paginación y ordenamiento
	var actions []models.Action
	if err := query.
		Order("name ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&actions).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener actions")
		return nil, 0, err
	}

	// Convertir a entidades de dominio
	domainActions := make([]domain.Action, len(actions))
	for i, action := range actions {
		var deletedAt *time.Time
		if action.DeletedAt.Valid {
			deletedAt = &action.DeletedAt.Time
		}

		domainActions[i] = domain.Action{
			ID:          action.ID,
			Name:        action.Name,
			Description: action.Description,
			CreatedAt:   action.CreatedAt,
			UpdatedAt:   action.UpdatedAt,
			DeletedAt:   deletedAt,
		}
	}

	r.logger.Info().Int("total", len(domainActions)).Int64("total_count", total).Msg("Actions obtenidos exitosamente")
	return domainActions, total, nil
}

// GetActionByID obtiene un action por su ID
func (r *Repository) GetActionByID(ctx context.Context, id uint) (*domain.Action, error) {
	r.logger.Info().Uint("id", id).Msg("Obteniendo action por ID")

	var action models.Action
	if err := r.database.Conn(ctx).Where("id = ?", id).First(&action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warn().Uint("id", id).Msg("Action no encontrado")
			return nil, fmt.Errorf("action con ID %d no encontrado", id)
		}
		r.logger.Error().Err(err).Uint("id", id).Msg("Error al obtener action por ID")
		return nil, err
	}

	var deletedAt *time.Time
	if action.DeletedAt.Valid {
		deletedAt = &action.DeletedAt.Time
	}

	domainAction := &domain.Action{
		ID:          action.ID,
		Name:        action.Name,
		Description: action.Description,
		CreatedAt:   action.CreatedAt,
		UpdatedAt:   action.UpdatedAt,
		DeletedAt:   deletedAt,
	}

	r.logger.Info().Uint("action_id", id).Str("name", action.Name).Msg("Action obtenido exitosamente")
	return domainAction, nil
}

// GetActionByName obtiene un action por su nombre
func (r *Repository) GetActionByName(ctx context.Context, name string) (*domain.Action, error) {
	r.logger.Info().Str("name", name).Msg("Obteniendo action por nombre")

	var action models.Action
	if err := r.database.Conn(ctx).Where("name = ?", name).First(&action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warn().Str("name", name).Msg("Action no encontrado por nombre")
			return nil, nil
		}
		r.logger.Error().Err(err).Str("name", name).Msg("Error al obtener action por nombre")
		return nil, err
	}

	var deletedAt *time.Time
	if action.DeletedAt.Valid {
		deletedAt = &action.DeletedAt.Time
	}

	domainAction := &domain.Action{
		ID:          action.ID,
		Name:        action.Name,
		Description: action.Description,
		CreatedAt:   action.CreatedAt,
		UpdatedAt:   action.UpdatedAt,
		DeletedAt:   deletedAt,
	}

	r.logger.Info().Uint("action_id", action.ID).Str("name", name).Msg("Action obtenido exitosamente por nombre")
	return domainAction, nil
}

// CreateAction crea un nuevo action
func (r *Repository) CreateAction(ctx context.Context, action domain.Action) (uint, error) {
	r.logger.Info().Str("name", action.Name).Msg("Creando nuevo action")

	modelAction := models.Action{
		Name:        strings.TrimSpace(action.Name),
		Description: strings.TrimSpace(action.Description),
	}

	if err := r.database.Conn(ctx).Create(&modelAction).Error; err != nil {
		r.logger.Error().Err(err).Str("name", action.Name).Msg("Error al crear action")
		return 0, err
	}

	r.logger.Info().
		Uint("action_id", modelAction.ID).
		Str("name", modelAction.Name).
		Msg("Action creado exitosamente")

	return modelAction.ID, nil
}

// UpdateAction actualiza un action existente
func (r *Repository) UpdateAction(ctx context.Context, id uint, action domain.Action) (string, error) {
	r.logger.Info().Uint("id", id).Str("name", action.Name).Msg("Actualizando action")

	// Verificar que el action existe
	var existingAction models.Action
	if err := r.database.Conn(ctx).Where("id = ?", id).First(&existingAction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warn().Uint("id", id).Msg("Action no encontrado para actualizar")
			return "", fmt.Errorf("action con ID %d no encontrado", id)
		}
		r.logger.Error().Err(err).Uint("id", id).Msg("Error al verificar action existente")
		return "", err
	}

	// Actualizar campos
	updates := map[string]interface{}{
		"name":        strings.TrimSpace(action.Name),
		"description": strings.TrimSpace(action.Description),
	}

	if err := r.database.Conn(ctx).
		Model(&models.Action{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Error al actualizar action")
		return "", err
	}

	r.logger.Info().Uint("action_id", id).Str("name", action.Name).Msg("Action actualizado exitosamente")
	return fmt.Sprintf("Action actualizado con ID: %d", id), nil
}

// DeleteAction elimina un action (soft delete)
func (r *Repository) DeleteAction(ctx context.Context, id uint) (string, error) {
	r.logger.Info().Uint("id", id).Msg("Eliminando action")

	// Verificar que el action existe
	var action models.Action
	if err := r.database.Conn(ctx).Where("id = ?", id).First(&action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warn().Uint("id", id).Msg("Action no encontrado para eliminar")
			return "", fmt.Errorf("action con ID %d no encontrado", id)
		}
		r.logger.Error().Err(err).Uint("id", id).Msg("Error al verificar action existente")
		return "", err
	}

	// Verificar si tiene permisos asociados
	var permissionsCount int64
	if err := r.database.Conn(ctx).
		Model(&models.Permission{}).
		Where("action_id = ?", id).
		Count(&permissionsCount).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Error al verificar permisos asociados")
		return "", err
	}

	if permissionsCount > 0 {
		r.logger.Warn().Uint("id", id).Int64("permissions_count", permissionsCount).Msg("Action tiene permisos asociados")
		return "", fmt.Errorf("no se puede eliminar el action porque tiene %d permiso(s) asociado(s)", permissionsCount)
	}

	// Realizar soft delete
	if err := r.database.Conn(ctx).
		Model(&models.Action{}).
		Where("id = ?", id).
		Delete(&models.Action{}).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Error al eliminar action")
		return "", err
	}

	r.logger.Info().Uint("action_id", id).Str("name", action.Name).Msg("Action eliminado exitosamente")
	return fmt.Sprintf("Action eliminado con ID: %d", id), nil
}

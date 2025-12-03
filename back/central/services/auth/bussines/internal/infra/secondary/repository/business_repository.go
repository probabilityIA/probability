package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type Repository struct {
	database db.IDatabase
	logger   log.ILogger
}

func New(database db.IDatabase, logger log.ILogger) domain.IBusinessRepository {
	return &Repository{
		database: database,
		logger:   logger,
	}
}

// GetBusinesses obtiene todos los negocios con paginación y filtros
func (r *Repository) GetBusinesses(ctx context.Context, page, perPage int, name string, businessTypeID *uint, isActive *bool) ([]domain.Business, int64, error) {
	var businesses []models.Business
	var total int64

	// Calcular offset
	offset := (page - 1) * perPage

	// Construir query base
	query := r.database.Conn(ctx).Model(&models.Business{})

	// Aplicar filtros
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if businessTypeID != nil {
		query = query.Where("business_type_id = ?", *businessTypeID)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Contar total con filtros aplicados
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al contar negocios")
		return nil, 0, err
	}

	// Obtener negocios con paginación y filtros
	if err := query.
		Preload("BusinessType").
		Limit(perPage).
		Offset(offset).
		Order("created_at DESC").
		Find(&businesses).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener negocios")
		return nil, 0, err
	}

	return mappers.ToBusinessEntitySlice(businesses), total, nil
}

// GetBusinessByID obtiene un negocio por su ID
func (r *Repository) GetBusinessByID(ctx context.Context, id uint) (*domain.Business, error) {
	var business models.Business
	if err := r.database.Conn(ctx).
		Model(&models.Business{}).
		Preload("BusinessType").
		Where("id = ?", id).
		First(&business).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al obtener negocio por ID")
		return nil, err
	}

	entity := mappers.ToBusinessEntity(business)
	return &entity, nil
}

// GetBusinessByCode obtiene un negocio por su código
func (r *Repository) GetBusinessByCode(ctx context.Context, code string) (*domain.Business, error) {
	var business models.Business
	if err := r.database.Conn(ctx).
		Model(&models.Business{}).
		Preload("BusinessType").
		Where("code = ?", code).
		First(&business).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBusinessNotFound
		}
		r.logger.Error().Str("code", code).Err(err).Msg("Error al obtener negocio por código")
		return nil, fmt.Errorf("error al consultar negocio por código: %w", err)
	}

	entity := mappers.ToBusinessEntity(business)
	return &entity, nil
}

// GetBusinessByCustomDomain obtiene un negocio por su dominio personalizado
func (r *Repository) GetBusinessByCustomDomain(ctx context.Context, customDomain string) (*domain.Business, error) {
	var business models.Business
	if err := r.database.Conn(ctx).
		Model(&models.Business{}).
		Preload("BusinessType").
		Where("custom_domain = ?", customDomain).
		First(&business).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBusinessNotFound
		}
		r.logger.Error().Str("domain", customDomain).Err(err).Msg("Error al obtener negocio por dominio personalizado")
		return nil, fmt.Errorf("error al consultar negocio por dominio personalizado: %w", err)
	}

	entity := mappers.ToBusinessEntity(business)
	return &entity, nil
}

// CreateBusiness crea un nuevo negocio
func (r *Repository) CreateBusiness(ctx context.Context, business domain.Business) (uint, error) {
	businessModel := mappers.ToBusinessModel(business)

	if err := r.database.Conn(ctx).Create(&businessModel).Error; err != nil {
		// Detectar errores de foreign key constraint y proporcionar mensaje más claro
		errMsg := err.Error()
		if strings.Contains(errMsg, "foreign key constraint") || strings.Contains(errMsg, "SQLSTATE 23503") {
			if strings.Contains(errMsg, "fk_business_type_businesses") {
				r.logger.Error().
					Err(err).
					Uint("business_type_id", business.BusinessTypeID).
					Msg("Error de foreign key: el tipo de negocio especificado no existe")
				return 0, fmt.Errorf("el tipo de negocio con ID %d no existe o no es válido", business.BusinessTypeID)
			}
		}
		r.logger.Error().Err(err).
			Str("name", business.Name).
			Str("code", business.Code).
			Uint("business_type_id", business.BusinessTypeID).
			Msg("Error al crear negocio en la base de datos")
		return 0, fmt.Errorf("error al guardar el negocio en la base de datos: %w", err)
	}

	// Obtener todos los recursos permitidos para el tipo de negocio
	permittedResources, err := r.GetBusinessTypeResourcesPermitted(ctx, business.BusinessTypeID)
	if err != nil {
		r.logger.Error().Err(err).Uint("business_type_id", business.BusinessTypeID).Msg("Error al obtener recursos permitidos")
		return 0, err
	}

	// Crear relaciones con todos los recursos permitidos (inactivas por defecto)
	for _, resource := range permittedResources {
		businessResource := models.BusinessResourceConfigured{
			BusinessID: businessModel.Model.ID,
			ResourceID: resource.ResourceID,
			Active:     false, // Nuevo negocio con recursos inactivos por defecto
		}

		if err := r.database.Conn(ctx).Create(&businessResource).Error; err != nil {
			r.logger.Error().Err(err).
				Uint("business_id", businessModel.Model.ID).
				Uint("resource_id", resource.ResourceID).
				Msg("Error al crear relación business-resource")
			return 0, err
		}
	}

	r.logger.Info().
		Uint("business_id", businessModel.Model.ID).
		Uint("business_type_id", business.BusinessTypeID).
		Int("resources_count", len(permittedResources)).
		Msg("Negocio creado con relaciones a recursos exitosamente")

	return businessModel.Model.ID, nil
}

// UpdateBusiness actualiza un negocio existente
func (r *Repository) UpdateBusiness(ctx context.Context, id uint, business domain.Business) (string, error) {
	businessModel := mappers.ToBusinessModel(business)

	if err := r.database.Conn(ctx).
		Model(&models.Business{}).
		Where("id = ?", id).
		Updates(&businessModel).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al actualizar negocio")
		return "", err
	}
	return fmt.Sprintf("Negocio actualizado con ID: %d", id), nil
}

// DeleteBusiness elimina un negocio
func (r *Repository) DeleteBusiness(ctx context.Context, id uint) (string, error) {
	if err := r.database.Conn(ctx).Delete(&models.Business{}, id).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al eliminar negocio")
		return "", err
	}
	return fmt.Sprintf("Negocio eliminado con ID: %d", id), nil
}

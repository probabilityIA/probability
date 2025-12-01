package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
)

// BusinessTypeRepository implementa ports.IBusinessTypeRepository

// GetBusinessTypes obtiene todos los tipos de negocio
func (r *Repository) GetBusinessTypes(ctx context.Context) ([]domain.BusinessType, error) {
	var businessTypes []domain.BusinessType
	if err := r.database.Conn(ctx).Table("business_type").Find(&businessTypes).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener tipos de negocio")
		return nil, err
	}
	return businessTypes, nil
}

// GetBusinessTypeByID obtiene un tipo de negocio por su ID
func (r *Repository) GetBusinessTypeByID(ctx context.Context, id uint) (*domain.BusinessType, error) {
	var businessType domain.BusinessType
	if err := r.database.Conn(ctx).Table("business_type").Where("id = ?", id).First(&businessType).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al obtener tipo de negocio por ID")
		return nil, err
	}
	return &businessType, nil
}

// GetBusinessTypeByCode obtiene un tipo de negocio por su código
func (r *Repository) GetBusinessTypeByCode(ctx context.Context, code string) (*domain.BusinessType, error) {
	var businessType domain.BusinessType
	if err := r.database.Conn(ctx).Table("business_type").Where("code = ?", code).First(&businessType).Error; err != nil {
		r.logger.Error().Str("code", code).Err(err).Msg("Error al obtener tipo de negocio por código")
		return nil, err
	}
	return &businessType, nil
}

// CreateBusinessType crea un nuevo tipo de negocio
func (r *Repository) CreateBusinessType(ctx context.Context, businessType domain.BusinessType) (string, error) {
	if err := r.database.Conn(ctx).Table("business_type").Create(&businessType).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al crear tipo de negocio")
		return "", err
	}
	return fmt.Sprintf("Tipo de negocio creado con ID: %d", businessType.ID), nil
}

// UpdateBusinessType actualiza un tipo de negocio existente
func (r *Repository) UpdateBusinessType(ctx context.Context, id uint, businessType domain.BusinessType) (string, error) {
	if err := r.database.Conn(ctx).Table("business_type").Where("id = ?", id).Updates(&businessType).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al actualizar tipo de negocio")
		return "", err
	}
	return fmt.Sprintf("Tipo de negocio actualizado con ID: %d", id), nil
}

// DeleteBusinessType elimina un tipo de negocio
func (r *Repository) DeleteBusinessType(ctx context.Context, id uint) (string, error) {
	if err := r.database.Conn(ctx).Table("business_type").Where("id = ?", id).Delete(&domain.BusinessType{}).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al eliminar tipo de negocio")
		return "", err
	}
	return fmt.Sprintf("Tipo de negocio eliminado con ID: %d", id), nil
}

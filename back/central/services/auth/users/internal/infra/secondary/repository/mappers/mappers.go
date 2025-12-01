package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ToUserModel convierte entities.User a models.User
func ToUserModel(user domain.UsersEntity) models.User {
	var deletedAt gorm.DeletedAt
	if user.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *user.DeletedAt, Valid: true}
	}

	return models.User{
		Model: gorm.Model{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: deletedAt,
		},
		Name:        user.Name,
		Email:       user.Email,
		Password:    user.Password,
		Phone:       user.Phone,
		AvatarURL:   user.AvatarURL,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
	}
}

// ToUserEntity convierte models.User a entities.User
func ToUserEntity(user models.User) domain.UsersEntity {
	var deletedAt *time.Time
	if user.Model.DeletedAt.Valid {
		deletedAt = &user.Model.DeletedAt.Time
	}

	return domain.UsersEntity{
		ID:          user.Model.ID,
		Name:        user.Name,
		Email:       user.Email,
		Password:    user.Password,
		Phone:       user.Phone,
		AvatarURL:   user.AvatarURL,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.Model.CreatedAt,
		UpdatedAt:   user.Model.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

// CreateUserModel crea un models.User para inserción (sin ID)
func CreateUserModel(user domain.UsersEntity) models.User {
	return models.User{
		Name:        user.Name,
		Email:       user.Email,
		Password:    user.Password,
		Phone:       user.Phone,
		AvatarURL:   user.AvatarURL,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
	}
}

// ToAPIKeyModel convierte una entidad APIKey del dominio al modelo de base de datos
func ToAPIKeyModel(apiKey domain.APIKey) models.APIKey {
	return models.APIKey{
		Model: gorm.Model{
			ID:        apiKey.ID,
			CreatedAt: apiKey.CreatedAt,
			UpdatedAt: apiKey.UpdatedAt,
		},
		UserID:      apiKey.UserID,
		BusinessID:  apiKey.BusinessID,
		CreatedByID: apiKey.CreatedByID,
		Name:        apiKey.Name,
		Description: apiKey.Description,
		KeyHash:     apiKey.KeyHash,
		LastUsedAt:  apiKey.LastUsedAt,
		Revoked:     apiKey.Revoked,
		RevokedAt:   apiKey.RevokedAt,
		RateLimit:   apiKey.RateLimit,
		IPWhitelist: apiKey.IPWhitelist,
	}
}

// ToAPIKeyEntity convierte un modelo APIKey de base de datos a entidad del dominio
func ToAPIKeyEntity(model models.APIKey) domain.APIKey {
	return domain.APIKey{
		ID:          model.Model.ID,
		UserID:      model.UserID,
		BusinessID:  model.BusinessID,
		CreatedByID: model.CreatedByID,
		Name:        model.Name,
		Description: model.Description,
		KeyHash:     model.KeyHash,
		LastUsedAt:  model.LastUsedAt,
		Revoked:     model.Revoked,
		RevokedAt:   model.RevokedAt,
		RateLimit:   model.RateLimit,
		IPWhitelist: model.IPWhitelist,
		CreatedAt:   model.Model.CreatedAt,
		UpdatedAt:   model.Model.UpdatedAt,
	}
}

// ToAPIKeyInfoEntity convierte un modelo APIKey a APIKeyInfo del dominio
func ToAPIKeyInfoEntity(model models.APIKey) domain.APIKeyInfo {
	return domain.APIKeyInfo{
		ID:          model.Model.ID,
		UserID:      model.UserID,
		BusinessID:  model.BusinessID,
		Name:        model.Name,
		Description: model.Description,
		LastUsedAt:  model.LastUsedAt,
		Revoked:     model.Revoked,
		RateLimit:   model.RateLimit,
		CreatedAt:   model.Model.CreatedAt,
	}
}

// ToAPIKeyInfoEntitySlice convierte un slice de modelos APIKey a slice de APIKeyInfo
func ToAPIKeyInfoEntitySlice(models []models.APIKey) []domain.APIKeyInfo {
	result := make([]domain.APIKeyInfo, len(models))
	for i, model := range models {
		result[i] = ToAPIKeyInfoEntity(model)
	}
	return result
}

// ToAPIKeyEntitySlice convierte un slice de modelos APIKey a slice de entidades APIKey
func ToAPIKeyEntitySlice(models []models.APIKey) []domain.APIKey {
	result := make([]domain.APIKey, len(models))
	for i, model := range models {
		result[i] = ToAPIKeyEntity(model)
	}
	return result
}

// CreateAPIKeyModel crea un modelo APIKey para inserción (sin ID)
func CreateAPIKeyModel(apiKey domain.APIKey, keyHash string) models.APIKey {
	return models.APIKey{
		UserID:      apiKey.UserID,
		BusinessID:  apiKey.BusinessID,
		CreatedByID: apiKey.CreatedByID,
		Name:        apiKey.Name,
		Description: apiKey.Description,
		KeyHash:     keyHash,
		RateLimit:   apiKey.RateLimit,
		IPWhitelist: apiKey.IPWhitelist,
	}
}

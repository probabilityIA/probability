package businessdata

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Reader) DescribeIdentity(ctx context.Context, userID uint, businessID *uint) (*entities.ChatIdentity, error) {
	identity := &entities.ChatIdentity{}

	var user models.User
	err := r.db.Conn(ctx).Model(&models.User{}).Select("id", "name").Where("id = ?", userID).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	identity.UserName = user.Name

	if businessID != nil && *businessID > 0 {
		var business models.Business
		err := r.db.Conn(ctx).Model(&models.Business{}).Select("id", "name").Where("id = ?", *businessID).First(&business).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		identity.BusinessName = business.Name
	}
	return identity, nil
}

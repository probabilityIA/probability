package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/siigoreferrals/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/siigoreferrals/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/siigoreferrals/internal/domain/errors"
)

var validOrderRanges = map[string]bool{
	"0-50":     true,
	"51-200":   true,
	"201-500":  true,
	"501-1000": true,
	"1000+":    true,
}

func (uc *UseCase) CreateReferral(ctx context.Context, dto dtos.CreateReferralDTO) error {
	name := strings.TrimSpace(dto.Name)
	email := strings.TrimSpace(dto.Email)
	phone := strings.TrimSpace(dto.Phone)
	orderRange := strings.TrimSpace(dto.OrderRange)
	if name == "" || email == "" || phone == "" || !validOrderRanges[orderRange] {
		return domainerrors.ErrInvalidReferral
	}

	referral := &entities.SiigoReferral{
		Name:       name,
		Email:      email,
		Phone:      phone,
		OrderRange: orderRange,
	}

	return uc.repo.CreateReferral(ctx, referral)
}

func (uc *UseCase) ListReferrals(ctx context.Context, page, pageSize int) ([]entities.SiigoReferral, dtos.ListReferralsResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	referrals, total, err := uc.repo.ListReferrals(ctx, page, pageSize)
	if err != nil {
		return nil, dtos.ListReferralsResult{}, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	return referrals, dtos.ListReferralsResult{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

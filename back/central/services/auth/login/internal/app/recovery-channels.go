package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
)

func (uc *AuthUseCase) RecoveryChannels(ctx context.Context, request domain.RecoveryChannelsRequest) (*domain.RecoveryChannelsResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))

	response := &domain.RecoveryChannelsResponse{
		Email:    true,
		WhatsApp: domain.WhatsAppChannelInfo{Available: false},
	}

	if email == "" {
		return response, nil
	}

	user, err := uc.repository.GetUserByEmail(ctx, email)
	if err != nil {
		uc.log.Error().Err(err).Str("email", email).Msg("Error buscando usuario para canales de recuperacion")
		return response, nil
	}
	if user == nil || !user.IsActive {
		return response, nil
	}

	if phone := strings.TrimSpace(user.Phone); phone != "" {
		response.WhatsApp = domain.WhatsAppChannelInfo{
			Available:   true,
			MaskedPhone: maskPhone(phone),
		}
	}

	return response, nil
}

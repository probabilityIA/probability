package usecaseintegrationtype

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

func (uc *integrationTypeUseCase) GetPlatformCredentialsByCode(ctx context.Context, code string) (map[string]interface{}, *domain.IntegrationType, error) {
	ctx = log.WithFunctionCtx(ctx, "GetPlatformCredentialsByCode")

	intType, err := uc.repo.GetIntegrationTypeByCode(ctx, code)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("code", code).Msg("integration_type not found")
		return nil, nil, fmt.Errorf("%w: %w", domain.ErrIntegrationTypeNotFound, err)
	}

	if len(intType.PlatformCredentialsEncrypted) == 0 {
		return map[string]interface{}{}, intType, nil
	}

	creds, err := uc.encryption.DecryptCredentials(ctx, intType.PlatformCredentialsEncrypted)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("code", code).Msg("decrypt platform credentials failed")
		return nil, intType, fmt.Errorf("decrypt credentials: %w", err)
	}

	return creds, intType, nil
}

func (uc *integrationTypeUseCase) RecordRevealAudit(ctx context.Context, audit *domain.CredentialRevealAudit) error {
	if audit == nil {
		return nil
	}
	if err := uc.repo.RecordCredentialReveal(ctx, audit); err != nil {
		uc.log.Warn(ctx).Err(err).Str("code", audit.IntegrationCode).Msg("record reveal audit failed")
		return err
	}
	return nil
}

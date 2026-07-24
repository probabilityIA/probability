package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
)

const dianIntegrationName = "Factus Probability (administrativa)"

func (uc *UseCase) GetDianConfig(ctx context.Context) (*dtos.DianConfigStatus, error) {
	status := &dtos.DianConfigStatus{}
	integrationID, err := uc.dianIntegration.GetGlobalIntegrationID(ctx)
	if err != nil {
		return status, nil
	}
	status.Configured = true
	status.IntegrationID = integrationID
	config, err := uc.dianIntegration.GetConfig(ctx, integrationID)
	if err != nil {
		return status, nil
	}
	status.NumberingRangeID = configInt(config, "numbering_range_id")
	status.MunicipalityID = configString(config, "municipality_id")
	status.DocumentIDType = configString(config, "identification_document_id")
	status.LegalOrgID = configString(config, "legal_organization_id")
	status.ApiURL = configString(config, "api_url_display")
	status.Username = configString(config, "username_display")
	return status, nil
}

func (uc *UseCase) SaveDianConfig(ctx context.Context, dto dtos.DianConfigDTO) (*dtos.DianConfigStatus, error) {
	dto.ClientID = strings.TrimSpace(dto.ClientID)
	dto.ClientSecret = strings.TrimSpace(dto.ClientSecret)
	dto.Username = strings.TrimSpace(dto.Username)
	dto.Password = strings.TrimSpace(dto.Password)
	dto.ApiURL = strings.TrimSpace(dto.ApiURL)
	if dto.ClientID == "" || dto.ClientSecret == "" || dto.Username == "" || dto.Password == "" {
		return nil, domainerrors.ErrDianCredentials
	}

	credentials := map[string]interface{}{
		"client_id":     dto.ClientID,
		"client_secret": dto.ClientSecret,
		"username":      dto.Username,
		"password":      dto.Password,
	}
	if dto.ApiURL != "" {
		credentials["api_url"] = dto.ApiURL
	}

	config := map[string]interface{}{
		"api_url_display":  dto.ApiURL,
		"username_display": dto.Username,
	}
	if dto.NumberingRangeID > 0 {
		config["numbering_range_id"] = dto.NumberingRangeID
	}
	if dto.MunicipalityID != "" {
		config["municipality_id"] = dto.MunicipalityID
	}
	if dto.DocumentIDType != "" {
		config["identification_document_id"] = dto.DocumentIDType
	}
	if dto.LegalOrgID != "" {
		config["legal_organization_id"] = dto.LegalOrgID
	}

	if err := uc.dianIntegration.TestConnection(ctx, config, credentials); err != nil {
		return nil, err
	}

	if _, err := uc.dianIntegration.Ensure(ctx, dianIntegrationName, config, credentials); err != nil {
		return nil, err
	}
	uc.log.Info(ctx).Msg("Accounting DIAN config saved (Factus platform integration)")
	return uc.GetDianConfig(ctx)
}

func (uc *UseCase) EmitInvoiceDian(ctx context.Context, dto dtos.EmitInvoiceDianDTO) (*entities.Invoice, error) {
	inv, err := uc.repo.GetInvoiceByID(ctx, dto.ID)
	if err != nil {
		return nil, err
	}
	if inv.Status == entities.InvoiceStatusCancelled {
		return nil, domainerrors.ErrInvoiceCancelled
	}
	if inv.Cufe != "" || inv.DianStatus == entities.DianStatusValidated {
		return nil, domainerrors.ErrDianAlreadyEmitted
	}

	document := strings.TrimSpace(dto.CustomerDocument)
	if document == "" {
		document = inv.CustomerDocument
	}
	if document == "" {
		return nil, domainerrors.ErrDianDocumentRequired
	}
	phone := strings.TrimSpace(dto.CustomerPhone)
	if phone == "" {
		phone = inv.CustomerPhone
	}
	address := strings.TrimSpace(dto.CustomerAddress)
	if address == "" {
		address = inv.CustomerAddress
	}
	if document != inv.CustomerDocument || phone != inv.CustomerPhone || address != inv.CustomerAddress {
		if err := uc.repo.SetInvoiceCustomer(ctx, inv.ID, document, phone, address); err != nil {
			return nil, err
		}
		inv.CustomerDocument = document
		inv.CustomerPhone = phone
		inv.CustomerAddress = address
	}

	integrationID, err := uc.dianIntegration.GetGlobalIntegrationID(ctx)
	if err != nil {
		return nil, domainerrors.ErrDianNotConfigured
	}

	config := map[string]interface{}{
		"reference_code": inv.Number,
	}

	result, err := uc.dianEmitter.Emit(ctx, integrationID, inv, config)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("invoice_id", inv.ID).Msg("Accounting DIAN emit failed")
		return nil, err
	}
	if err := uc.repo.SetInvoiceDian(ctx, inv.ID, result); err != nil {
		return nil, err
	}
	uc.log.Info(ctx).
		Uint("invoice_id", inv.ID).
		Str("dian_number", result.Number).
		Str("cufe", result.CUFE).
		Msg("Accounting invoice emitted to DIAN via Factus")
	return uc.repo.GetInvoiceByID(ctx, inv.ID)
}

func configString(config map[string]interface{}, key string) string {
	if config == nil {
		return ""
	}
	if v, ok := config[key].(string); ok {
		return v
	}
	return ""
}

func configInt(config map[string]interface{}, key string) int {
	if config == nil {
		return 0
	}
	switch v := config[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	}
	return 0
}

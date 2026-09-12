package usecasenumbers

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	whaerrors "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

func (u *usecase) GetState(ctx context.Context, businessID uint) (*NumberState, error) {
	integrationID, config, credentials, err := u.load(ctx, businessID)
	if err != nil {
		return nil, err
	}

	state := stateFromConfig(integrationID, businessID, config)

	if state.PhoneNumberID == "" {
		return state, nil
	}

	platform, err := u.credentialsCache.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return state, nil
	}

	api := u.apiFactory(platform.WhatsAppURL)
	number, err := api.GetPhoneNumber(ctx, state.PhoneNumberID, tokenFor(credentials, platform.AccessToken))
	if err != nil {
		u.log.Warn(ctx).Err(err).
			Str("phone_number_id", state.PhoneNumberID).
			Msg("no se pudo consultar el estado del número en Meta")
		return state, nil
	}

	applyRemote(state, number)

	return state, nil
}

func (u *usecase) AddNumber(ctx context.Context, businessID uint, input AddNumberInput) (*NumberState, error) {
	integrationID, config, _, err := u.load(ctx, businessID)
	if err != nil {
		return nil, err
	}

	countryCode := onlyDigits(input.CountryCode)
	phoneNumber := onlyDigits(input.PhoneNumber)
	verifiedName := strings.TrimSpace(input.VerifiedName)

	if verifiedName == "" {
		return nil, fmt.Errorf("se necesita el nombre que verán tus clientes")
	}
	if err := validarNumero(countryCode, phoneNumber); err != nil {
		return nil, err
	}

	if existing := stringValue(config["phone_number_id"]); existing != "" {
		return nil, fmt.Errorf("este negocio ya tiene un número en proceso: quítalo antes de agregar otro")
	}

	platform, err := u.credentialsCache.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("credenciales de plataforma no disponibles: %w", err)
	}
	if platform.WABAID == "" {
		return nil, fmt.Errorf("las credenciales de plataforma no tienen waba_id: no se puede agregar el número")
	}

	api := u.apiFactory(platform.WhatsAppURL)

	phoneNumberID, err := api.AddPhoneNumber(ctx, platform.WABAID, platform.AccessToken, countryCode, phoneNumber, verifiedName)
	if err != nil {
		return nil, fmt.Errorf("Meta no aceptó el número: %w", err)
	}

	if err := u.saveConfig(ctx, integrationID, map[string]any{
		"phone_number_id":    phoneNumberID,
		"waba_id":            platform.WABAID,
		"hosted_by_platform": true,
		"number_status":      StatusEsperandoCodigo,
		"verified_name":      verifiedName,
		"use_platform_token": true,
	}); err != nil {
		return nil, err
	}

	u.log.Info(ctx).
		Uint("business_id", businessID).
		Str("phone_number_id", phoneNumberID).
		Msg("número agregado al WABA de Probability, falta verificar el código")

	return &NumberState{
		IntegrationID:    integrationID,
		BusinessID:       businessID,
		PhoneNumberID:    phoneNumberID,
		Status:           StatusEsperandoCodigo,
		VerifiedName:     verifiedName,
		HostedByPlatform: true,
	}, nil
}

func (u *usecase) RequestCode(ctx context.Context, businessID uint, method string) (*NumberState, error) {
	integrationID, config, credentials, err := u.load(ctx, businessID)
	if err != nil {
		return nil, err
	}

	state := stateFromConfig(integrationID, businessID, config)
	if state.PhoneNumberID == "" {
		return nil, fmt.Errorf("todavía no hay un número agregado")
	}

	platform, err := u.credentialsCache.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("credenciales de plataforma no disponibles: %w", err)
	}

	api := u.apiFactory(platform.WhatsAppURL)
	if err := api.RequestCode(ctx, state.PhoneNumberID, tokenFor(credentials, platform.AccessToken), method, "es"); err != nil {
		return nil, fmt.Errorf("Meta no envió el código: %w", err)
	}

	if err := u.saveConfig(ctx, integrationID, map[string]any{"number_status": StatusEsperandoCodigo}); err != nil {
		return nil, err
	}

	state.Status = StatusEsperandoCodigo

	u.log.Info(ctx).
		Uint("business_id", businessID).
		Str("phone_number_id", state.PhoneNumberID).
		Str("metodo", method).
		Msg("código de verificación solicitado a Meta")

	return state, nil
}

func (u *usecase) VerifyCode(ctx context.Context, businessID uint, code string) (*NumberState, error) {
	integrationID, config, credentials, err := u.load(ctx, businessID)
	if err != nil {
		return nil, err
	}

	code = onlyDigits(code)
	if code == "" {
		return nil, fmt.Errorf("falta el código que llegó al teléfono")
	}

	state := stateFromConfig(integrationID, businessID, config)
	if state.PhoneNumberID == "" {
		return nil, fmt.Errorf("todavía no hay un número agregado")
	}

	platform, err := u.credentialsCache.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("credenciales de plataforma no disponibles: %w", err)
	}

	api := u.apiFactory(platform.WhatsAppURL)
	if err := api.VerifyCode(ctx, state.PhoneNumberID, tokenFor(credentials, platform.AccessToken), code); err != nil {
		return nil, fmt.Errorf("el código no fue aceptado: %w", err)
	}

	if err := u.saveConfig(ctx, integrationID, map[string]any{"number_status": StatusVerificado}); err != nil {
		return nil, err
	}

	state.Status = StatusVerificado

	u.log.Info(ctx).
		Uint("business_id", businessID).
		Str("phone_number_id", state.PhoneNumberID).
		Msg("número verificado, falta registrarlo en la Cloud API")

	return state, nil
}

func (u *usecase) Register(ctx context.Context, businessID uint) (*NumberState, error) {
	integrationID, config, credentials, err := u.load(ctx, businessID)
	if err != nil {
		return nil, err
	}

	state := stateFromConfig(integrationID, businessID, config)
	if state.PhoneNumberID == "" {
		return nil, fmt.Errorf("todavía no hay un número agregado")
	}

	platform, err := u.credentialsCache.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("credenciales de plataforma no disponibles: %w", err)
	}

	pin := stringValue(credentials["two_factor_pin"])
	nuevoPin := pin == ""
	if nuevoPin {
		pin, err = generarPin()
		if err != nil {
			return nil, fmt.Errorf("no se pudo generar el PIN de dos pasos: %w", err)
		}
	}

	api := u.apiFactory(platform.WhatsAppURL)
	if err := api.Register(ctx, state.PhoneNumberID, tokenFor(credentials, platform.AccessToken), pin); err != nil {
		return nil, fmt.Errorf("Meta no registró el número: %w", err)
	}

	if nuevoPin {
		if err := u.resolver.UpdateIntegrationCredentials(ctx, strconv.FormatUint(uint64(integrationID), 10), map[string]any{
			"two_factor_pin": pin,
		}); err != nil {
			return nil, fmt.Errorf("el número quedó registrado pero no se pudo guardar el PIN: %w", err)
		}
	}

	if err := u.saveConfig(ctx, integrationID, map[string]any{
		"number_status":      StatusRegistrado,
		"use_platform_token": false,
	}); err != nil {
		return nil, err
	}

	state.Status = StatusRegistrado
	state.Active = true
	if nuevoPin {
		state.Pin = pin
	}

	u.log.Info(ctx).
		Uint("business_id", businessID).
		Str("phone_number_id", state.PhoneNumberID).
		Msg("número registrado: el negocio ya envía y recibe por su propio número")

	return state, nil
}

func (u *usecase) RemoveNumber(ctx context.Context, businessID uint) (*NumberState, error) {
	integrationID, config, credentials, err := u.load(ctx, businessID)
	if err != nil {
		return nil, err
	}

	state := stateFromConfig(integrationID, businessID, config)
	if state.PhoneNumberID == "" {
		return nil, fmt.Errorf("este negocio no tiene ningún número en proceso")
	}

	platform, err := u.credentialsCache.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("credenciales de plataforma no disponibles: %w", err)
	}

	if state.HostedByPlatform {
		api := u.apiFactory(platform.WhatsAppURL)
		token := tokenFor(credentials, platform.AccessToken)

		if state.Status == StatusRegistrado || state.Active {
			if err := api.DeregisterPhoneNumber(ctx, state.PhoneNumberID, token); err != nil {
				u.log.Warn(ctx).Err(err).
					Str("phone_number_id", state.PhoneNumberID).
					Msg("no se pudo dar de baja el número de la Cloud API: se intenta borrarlo igual")
			}
		}

		if err := api.DeletePhoneNumber(ctx, state.PhoneNumberID, token); err != nil {
			if !numeroYaNoExisteEnMeta(err) {
				return nil, fmt.Errorf("Meta no dejó quitar el número: %w", err)
			}
			u.log.Warn(ctx).Err(err).
				Str("phone_number_id", state.PhoneNumberID).
				Msg("el número ya no existía en Meta: se limpia la configuración local")
		}
	}

	if err := u.saveConfig(ctx, integrationID, map[string]any{
		"phone_number_id":    "",
		"waba_id":            "",
		"number_status":      "",
		"verified_name":      "",
		"hosted_by_platform": false,
		"use_platform_token": true,
	}); err != nil {
		return nil, err
	}

	u.log.Info(ctx).
		Uint("business_id", businessID).
		Str("phone_number_id", state.PhoneNumberID).
		Msg("número quitado: el negocio vuelve a enviar por el número de Probability")

	return &NumberState{
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Status:        StatusSinNumero,
	}, nil
}

func (u *usecase) load(ctx context.Context, businessID uint) (uint, map[string]any, map[string]any, error) {
	if u.resolver == nil {
		return 0, nil, nil, fmt.Errorf("el módulo de integraciones no está disponible")
	}
	if businessID == 0 {
		return 0, nil, nil, fmt.Errorf("business_id es requerido")
	}

	integrationID, err := u.resolver.GetIntegrationIDByBusinessAndType(ctx, businessID, whatsAppTypeID)
	if err != nil || integrationID == 0 {
		return 0, nil, nil, fmt.Errorf("el negocio %d no tiene integración de WhatsApp", businessID)
	}

	config, credentials, err := u.resolver.GetIntegrationConfigAndCredentials(ctx, integrationID)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("no se pudo leer la integración de WhatsApp: %w", err)
	}

	return integrationID, config, credentials, nil
}

func (u *usecase) saveConfig(ctx context.Context, integrationID uint, config map[string]any) error {
	if err := u.resolver.UpdateIntegrationConfig(ctx, strconv.FormatUint(uint64(integrationID), 10), config); err != nil {
		return fmt.Errorf("error guardando el estado del número: %w", err)
	}
	return nil
}

func stateFromConfig(integrationID, businessID uint, config map[string]any) *NumberState {
	status := stringValue(config["number_status"])
	if status == "" {
		status = StatusSinNumero
	}

	return &NumberState{
		IntegrationID:    integrationID,
		BusinessID:       businessID,
		PhoneNumberID:    stringValue(config["phone_number_id"]),
		Status:           status,
		VerifiedName:     stringValue(config["verified_name"]),
		HostedByPlatform: config["hosted_by_platform"] == true,
		Active:           config["use_platform_token"] == false,
	}
}

func applyRemote(state *NumberState, number *ports.WABAPhoneNumber) {
	state.DisplayPhoneNumber = number.DisplayPhoneNumber
	state.QualityRating = number.QualityRating
	state.CodeVerification = number.CodeVerificationStatus
	state.NameStatus = number.NameStatus
	if number.VerifiedName != "" {
		state.VerifiedName = number.VerifiedName
	}
	if strings.EqualFold(number.NameStatus, "PENDING_REVIEW") && state.Status == StatusRegistrado {
		state.Status = StatusNombreEnRevision
	}
}

func tokenFor(credentials map[string]any, platformToken string) string {
	if token := stringValue(credentials["access_token"]); token != "" {
		return token
	}
	return platformToken
}

func stringValue(raw any) string {
	value, _ := raw.(string)
	return strings.TrimSpace(value)
}

func onlyDigits(raw string) string {
	var b strings.Builder
	for _, r := range raw {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func generarPin() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func numeroYaNoExisteEnMeta(err error) bool {
	var metaErr *whaerrors.MetaGraphError
	if !errors.As(err, &metaErr) {
		return false
	}
	if metaErr.StatusCode == 404 {
		return true
	}
	return metaErr.Code == 803 || (metaErr.Code == 100 && metaErr.Subcode == 33)
}

func validarNumero(countryCode, phoneNumber string) error {
	if countryCode == "" || phoneNumber == "" {
		return fmt.Errorf("se necesitan el indicativo del país y el número")
	}
	if len(countryCode) > 3 {
		return fmt.Errorf("el indicativo del país no puede pasar de 3 dígitos")
	}
	if len(countryCode)+len(phoneNumber) > 15 {
		return fmt.Errorf("indicativo y número suman %d dígitos: un número internacional no pasa de 15", len(countryCode)+len(phoneNumber))
	}
	if countryCode == "57" && (len(phoneNumber) != 10 || phoneNumber[0] != '3') {
		return fmt.Errorf("un celular colombiano tiene 10 dígitos y empieza por 3; escribiste %d (%s)", len(phoneNumber), phoneNumber)
	}
	if len(phoneNumber) < 6 {
		return fmt.Errorf("el número es demasiado corto")
	}
	return nil
}

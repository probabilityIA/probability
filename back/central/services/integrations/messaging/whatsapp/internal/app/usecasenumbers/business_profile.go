package usecasenumbers

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

const (
	maxAbout       = 139
	maxDescription = 512
	maxWebsites    = 2
)

var verticalesValidas = map[string]bool{
	"OTHER": true, "AUTO": true, "BEAUTY": true, "APPAREL": true, "EDU": true,
	"ENTERTAIN": true, "EVENT_PLAN": true, "FINANCE": true, "GROCERY": true,
	"GOVT": true, "HOTEL": true, "HEALTH": true, "NONPROFIT": true,
	"PROF_SERVICES": true, "RETAIL": true, "TRAVEL": true, "RESTAURANT": true,
	"ALCOHOL": true, "ONLINE_GAMBLING": true, "PHYSICAL_GAMBLING": true,
	"OTC_DRUGS": true, "MATRIMONY_SERVICE": true,
}

type BusinessProfile struct {
	About             string   `json:"about"`
	Address           string   `json:"address"`
	Description       string   `json:"description"`
	Email             string   `json:"email"`
	ProfilePictureURL string   `json:"profile_picture_url"`
	Websites          []string `json:"websites"`
	Vertical          string   `json:"vertical"`
}

type BusinessProfileInput struct {
	About       *string
	Address     *string
	Description *string
	Email       *string
	Vertical    *string
	Websites    *[]string
}

func (u *usecase) GetProfile(ctx context.Context, businessID uint) (*BusinessProfile, error) {
	phoneNumberID, token, baseURL, err := u.perfilTarget(ctx, businessID)
	if err != nil {
		return nil, err
	}

	remoto, err := u.profileFactory(baseURL).GetProfile(ctx, phoneNumberID, token)
	if err != nil {
		return nil, err
	}

	return &BusinessProfile{
		About:             remoto.About,
		Address:           remoto.Address,
		Description:       remoto.Description,
		Email:             remoto.Email,
		ProfilePictureURL: remoto.ProfilePictureURL,
		Websites:          remoto.Websites,
		Vertical:          remoto.Vertical,
	}, nil
}

func (u *usecase) UpdateProfile(ctx context.Context, businessID uint, input BusinessProfileInput) (*BusinessProfile, error) {
	if err := validarPerfil(input); err != nil {
		return nil, err
	}

	phoneNumberID, token, baseURL, err := u.perfilTarget(ctx, businessID)
	if err != nil {
		return nil, err
	}

	update := ports.WhatsAppBusinessProfileUpdate{
		About:       input.About,
		Address:     input.Address,
		Description: input.Description,
		Email:       input.Email,
		Vertical:    input.Vertical,
		Websites:    input.Websites,
	}

	if err := u.profileFactory(baseURL).UpdateProfile(ctx, phoneNumberID, token, update); err != nil {
		return nil, err
	}

	u.log.Info(ctx).Uint("business_id", businessID).Msg("perfil de WhatsApp actualizado")

	return u.GetProfile(ctx, businessID)
}

func (u *usecase) UpdateProfilePhoto(ctx context.Context, businessID uint, contentType string, data []byte) (*BusinessProfile, error) {
	if !strings.HasPrefix(contentType, "image/") {
		return nil, fmt.Errorf("el archivo debe ser una imagen")
	}

	phoneNumberID, token, baseURL, err := u.perfilTarget(ctx, businessID)
	if err != nil {
		return nil, err
	}

	appID, err := u.appIDPlataforma(ctx)
	if err != nil {
		return nil, err
	}

	api := u.profileFactory(baseURL)

	handle, err := api.UploadProfilePhoto(ctx, appID, token, contentType, data)
	if err != nil {
		return nil, err
	}

	if err := api.UpdateProfile(ctx, phoneNumberID, token, ports.WhatsAppBusinessProfileUpdate{
		ProfilePictureHandle: &handle,
	}); err != nil {
		return nil, err
	}

	u.log.Info(ctx).Uint("business_id", businessID).Msg("foto del perfil de WhatsApp actualizada")

	return u.GetProfile(ctx, businessID)
}

func (u *usecase) appIDPlataforma(ctx context.Context) (string, error) {
	if u.resolver == nil {
		return "", fmt.Errorf("el modulo de integraciones no esta disponible")
	}

	creds, err := u.resolver.GetCachedPlatformCredentials(ctx, whatsAppTypeID)
	if err != nil {
		return "", fmt.Errorf("no se pudieron leer las credenciales de la plataforma: %w", err)
	}

	appID := stringValue(creds["app_id"])
	if appID == "" {
		return "", fmt.Errorf("la plataforma no tiene app_id configurado: sin el Meta no acepta la foto")
	}

	return appID, nil
}

func (u *usecase) perfilTarget(ctx context.Context, businessID uint) (string, string, string, error) {
	if u.profileFactory == nil {
		return "", "", "", fmt.Errorf("la edicion del perfil no esta disponible")
	}

	_, config, credentials, err := u.load(ctx, businessID)
	if err != nil {
		return "", "", "", err
	}

	phoneNumberID := stringValue(config["phone_number_id"])
	if phoneNumberID == "" {
		return "", "", "", fmt.Errorf("el negocio todavia no tiene un numero propio conectado")
	}

	platform, err := u.credentialsCache.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return "", "", "", fmt.Errorf("no se pudo leer la configuracion de la plataforma: %w", err)
	}

	return phoneNumberID, tokenFor(credentials, platform.AccessToken), platform.WhatsAppURL, nil
}

func validarPerfil(input BusinessProfileInput) error {
	if input.About != nil && len([]rune(*input.About)) > maxAbout {
		return fmt.Errorf("el texto de informacion no puede pasar de %d caracteres", maxAbout)
	}
	if input.Description != nil && len([]rune(*input.Description)) > maxDescription {
		return fmt.Errorf("la descripcion no puede pasar de %d caracteres", maxDescription)
	}
	if input.Websites != nil {
		if len(*input.Websites) > maxWebsites {
			return fmt.Errorf("WhatsApp acepta maximo %d sitios web", maxWebsites)
		}
		for _, sitio := range *input.Websites {
			if sitio == "" {
				continue
			}
			if !strings.HasPrefix(sitio, "http://") && !strings.HasPrefix(sitio, "https://") {
				return fmt.Errorf("el sitio web %q debe empezar por http:// o https://", sitio)
			}
		}
	}
	if input.Email != nil && *input.Email != "" && !strings.Contains(*input.Email, "@") {
		return fmt.Errorf("el correo no es valido")
	}
	if input.Vertical != nil && *input.Vertical != "" && !verticalesValidas[strings.ToUpper(*input.Vertical)] {
		return fmt.Errorf("la categoria %q no es una de las que acepta WhatsApp", *input.Vertical)
	}
	return nil
}

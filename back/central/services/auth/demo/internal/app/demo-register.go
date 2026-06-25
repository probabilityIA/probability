package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

const emailVerificationTTL = 24 * time.Hour

func (uc *UseCase) DemoRegister(ctx context.Context, request domain.DemoRegisterRequest) (*domain.DemoRegisterResponse, error) {
	fullName := strings.TrimSpace(request.FullName)
	businessName := strings.TrimSpace(request.BusinessName)
	email := strings.TrimSpace(strings.ToLower(request.Email))
	password := request.Password

	if fullName == "" || businessName == "" || email == "" {
		return nil, fmt.Errorf("nombre, negocio y correo son obligatorios")
	}
	if len(password) < 6 {
		return nil, fmt.Errorf("la contrasena debe tener al menos 6 caracteres")
	}

	exists, err := uc.repository.EmailExists(ctx, email)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error verificando email en registro demo")
		return nil, fmt.Errorf("error interno del servidor")
	}
	if exists {
		return nil, fmt.Errorf("el correo ya esta registrado")
	}

	roleID, err := uc.repository.GetDemoRoleID(ctx)
	if err != nil || roleID == 0 {
		uc.log.Error().Err(err).Msg("Rol demo no encontrado")
		return nil, fmt.Errorf("error interno del servidor")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error procesando la contrasena")
	}

	businessCode, err := uc.uniqueBusinessCode(ctx, businessName)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error generando codigo de negocio")
		return nil, fmt.Errorf("error interno del servidor")
	}

	rawToken, tokenHash, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("error interno del servidor")
	}

	params := domain.CreateDemoAccountParams{
		FullName:     fullName,
		BusinessName: businessName,
		BusinessCode: businessCode,
		OrderPrefix:  derivePrefix(businessName),
		Email:        email,
		PasswordHash: string(hashed),
		RoleID:       roleID,
		TokenHash:    tokenHash,
		ExpiresAt:    time.Now().Add(emailVerificationTTL),
	}

	if _, err := uc.repository.CreateDemoAccount(ctx, params); err != nil {
		uc.log.Error().Err(err).Str("email", email).Msg("Error creando cuenta demo")
		return nil, fmt.Errorf("no se pudo crear la cuenta demo")
	}

	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", uc.frontendBaseURL(), rawToken)
	html := buildVerificationEmail(fullName, businessName, verifyURL)
	if err := uc.emailSender.SendHTML(ctx, email, "Verifica tu cuenta demo de Probability", html); err != nil {
		uc.log.Error().Err(err).Str("email", email).Msg("Error enviando correo de verificacion demo")
	}

	uc.log.Info().Str("email", email).Str("business", businessName).Msg("Cuenta demo creada, correo de verificacion enviado")
	return &domain.DemoRegisterResponse{
		Success: true,
		Message: "Cuenta creada. Revisa tu correo para verificar tu cuenta y empezar.",
	}, nil
}

func (uc *UseCase) uniqueBusinessCode(ctx context.Context, businessName string) (string, error) {
	base := slugify(businessName)
	if base == "" {
		base = "demo"
	}
	for i := 0; i < 6; i++ {
		suffix, err := randomHex(3)
		if err != nil {
			return "", err
		}
		code := base + "-" + suffix
		if len(code) > 50 {
			code = code[:50]
		}
		exists, err := uc.repository.BusinessCodeExists(ctx, code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", fmt.Errorf("no se pudo generar un codigo de negocio unico")
}

func (uc *UseCase) frontendBaseURL() string {
	base := strings.TrimRight(uc.env.Get("FRONTEND_BASE_URL"), "/")
	if base == "" {
		base = "http://localhost:3000"
	}
	return base
}

func generateToken() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	raw := hex.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	return raw, hex.EncodeToString(sum[:]), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func slugify(s string) string {
	var sb strings.Builder
	prevDash := false
	for _, r := range strings.ToLower(s) {
		switch {
		case unicode.IsLetter(r) && r < 128, unicode.IsDigit(r) && r < 128:
			sb.WriteRune(r)
			prevDash = false
		case r == ' ' || r == '-' || r == '_':
			if !prevDash && sb.Len() > 0 {
				sb.WriteRune('-')
				prevDash = true
			}
		}
	}
	return strings.Trim(sb.String(), "-")
}

func derivePrefix(name string) string {
	letters := make([]rune, 0, 3)
	for _, r := range name {
		if unicode.IsLetter(r) && r < 128 {
			letters = append(letters, unicode.ToUpper(r))
			if len(letters) == 3 {
				break
			}
		}
	}
	if len(letters) == 0 {
		return "DEM"
	}
	for len(letters) < 3 {
		letters = append(letters, 'X')
	}
	return string(letters)
}

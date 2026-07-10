package domain

import "errors"

var (
	ErrEmailAlreadyRegistered   = errors.New("el correo ya esta registrado")
	ErrEmailPendingVerification = errors.New("la cuenta existe pero no ha sido verificada")
)

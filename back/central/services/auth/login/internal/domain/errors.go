package domain

import "errors"

var (
	ErrUserNameRequired  = errors.New("el nombre del usuario es requerido")
	ErrUserEmailRequired = errors.New("el email del usuario es requerido")
	ErrUserEmailInvalid  = errors.New("el email no tiene un formato válido")
	ErrUserPhoneInvalid  = errors.New("el teléfono debe tener exactamente 10 dígitos")
	ErrUserPasswordError = errors.New("error al generar o procesar la contraseña")

	ErrUserNotFound            = errors.New("usuario no encontrado")
	ErrUserEmailExists         = errors.New("el email ya está registrado")
	ErrUserInactive            = errors.New("usuario inactivo")
	ErrUserPendingVerification = errors.New("tu cuenta aun no ha sido verificada")
	ErrUserCannotBeDeleted     = errors.New("no se puede eliminar el usuario")
	ErrBusinessesNotFound      = errors.New("algunos businesses no existen")
	ErrRolesNotFound           = errors.New("algunos roles no existen")
	ErrUserAvatarUploadFailed  = errors.New("error al subir imagen de avatar")
	ErrUserAvatarDeleteFailed  = errors.New("error al eliminar imagen de avatar")

	ErrInvalidCredentials    = errors.New("credenciales inválidas")
	ErrEmailPasswordRequired = errors.New("email y contraseña son requeridos")
)

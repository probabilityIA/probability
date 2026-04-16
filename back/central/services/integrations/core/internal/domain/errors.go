package domain

import "errors"

var (
	// Errores de validación de integración
	ErrIntegrationNameRequired    = errors.New("el nombre de la integración es obligatorio")
	ErrIntegrationCodeRequired    = errors.New("el código de la integración es obligatorio")
	ErrIntegrationTypeRequired    = errors.New("el tipo de integración es obligatorio")
	ErrIntegrationCategoryInvalid = errors.New("categoría de integración inválida")

	// Errores de negocio de integración
	ErrIntegrationNotFound             = errors.New("integración no encontrada")
	ErrIntegrationCodeExists           = errors.New("ya existe una integración con el código")
	ErrIntegrationCannotDeleteWhatsApp = errors.New("no se puede eliminar la integración de WhatsApp. Solo se puede desactivar")
	ErrIntegrationTypeNotFound         = errors.New("tipo de integración no encontrado")
	ErrIntegrationCredentialsDecrypt   = errors.New("error al desencriptar credenciales")
	ErrIntegrationConfigDeserialize    = errors.New("error al deserializar configuración")
	ErrIntegrationCredentialsSerialize = errors.New("error al serializar credenciales")
	ErrIntegrationConfigSerialize      = errors.New("error al serializar configuración")
	ErrIntegrationTestFailed           = errors.New("test de conexión falló")
	ErrIntegrationAccessTokenNotFound  = errors.New("access_token no encontrado o inválido en las credenciales")

	// Errores de validación de tipo de integración
	ErrIntegrationTypeNameRequired = errors.New("el nombre del tipo de integración es obligatorio")
	ErrIntegrationTypeCodeRequired = errors.New("el código del tipo de integración es obligatorio")

	// Errores de negocio de tipo de integración
	ErrIntegrationTypeNameExists        = errors.New("el nombre del tipo de integración ya está en uso")
	ErrIntegrationTypeCodeExists        = errors.New("el código del tipo de integración ya está en uso")
	ErrIntegrationTypeHasIntegrations   = errors.New("no se puede eliminar un tipo de integración que tiene integraciones asociadas")
	ErrIntegrationTypeImageUploadFailed = errors.New("error al subir imagen del tipo de integración")
	ErrIntegrationTypeImageDeleteFailed = errors.New("error al eliminar imagen del tipo de integración")

)

package errors

import "errors"

var (
	ErrEventNotValid     = errors.New("evento no válido")
	ErrNoConfigsFound    = errors.New("no se encontraron configuraciones de notificación")
	ErrPublishFailed     = errors.New("error al publicar evento")
	ErrSerializeFailed   = errors.New("error al serializar evento")
	ErrDeserializeFailed = errors.New("error al deserializar evento")
)

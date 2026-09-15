package errors

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrEmptyConversation = errors.New("la conversacion no tiene un mensaje del usuario")
	ErrMessageTooLong    = errors.New("el mensaje es demasiado largo")
	ErrModelUnavailable  = errors.New("el asistente no esta disponible")
)

type RateLimitedError struct {
	Limit   int
	ResetAt *time.Time
}

func (e *RateLimitedError) Error() string {
	return fmt.Sprintf("limite de %d mensajes por hora alcanzado", e.Limit)
}

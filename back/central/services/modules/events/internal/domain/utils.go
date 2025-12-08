package domain

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateEventID genera un ID único para eventos
func GenerateEventID() string {
	// Generar 12 bytes aleatorios
	b := make([]byte, 12)
	rand.Read(b)
	// Codificar en base64 URL-safe y tomar los primeros 16 caracteres
	return base64.URLEncoding.EncodeToString(b)[:16]
}

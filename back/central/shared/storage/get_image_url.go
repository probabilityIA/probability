package storage

import (
	"fmt"
	"strings"
)

// GetImageURL genera la URL pública de la imagen
func (s *S3Uploader) GetImageURL(filename string) string {
	// Si hay URL_BASE_DOMAIN_S3 configurada, usarla
	// Esto se maneja en los casos de uso que usan URL_BASE_DOMAIN_S3
	// Por defecto, usar formato AWS S3
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, strings.TrimLeft(filename, "/"))
}

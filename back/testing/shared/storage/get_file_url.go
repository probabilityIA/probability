package storage

import (
	"context"
	"fmt"
)

// GetFileURL genera la URL pública del archivo
func (s *S3Uploader) GetFileURL(ctx context.Context, filename string) (string, error) {
	// Para simplicidad, asumimos bucket público o política que permite acceso
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, filename), nil
}

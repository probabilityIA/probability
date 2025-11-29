package storage

import "fmt"

// GetImageURL genera la URL pública de la imagen
func (s *S3Uploader) GetImageURL(filename string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, filename)
}

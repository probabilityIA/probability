package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ImageExists verifica si una imagen existe en S3
func (s *S3Uploader) ImageExists(ctx context.Context, filename string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		// Si hay error, asumimos que el archivo no existe
		return false, nil
	}

	return true, nil
}

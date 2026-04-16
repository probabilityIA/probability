package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// UploadFile mantiene la funcionalidad original para archivos generales
func (s *S3Uploader) UploadFile(ctx context.Context, file io.ReadSeeker, filename string) (string, error) {
	if file == nil {
		return "", fmt.Errorf("archivo es nulo")
	}

	// Detect content type from extension
	contentType := "application/octet-stream"
	if strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		contentType = "application/pdf"
	}

	// ServerSideEncryption removido: requiere KMS configurado, no compatible con MinIO local
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:             aws.String(s.bucket),
		Key:                aws.String(filename),
		Body:               file,
		ContentType:        aws.String(contentType),
		ContentDisposition: aws.String("inline"),
		StorageClass:       types.StorageClassIntelligentTiering,
	})
	if err != nil {
		s.log.Error().Err(err).Msg("error subiendo archivo a S3")
		return "", err
	}

	url := s.GetImageURL(filename)
	return url, nil
}

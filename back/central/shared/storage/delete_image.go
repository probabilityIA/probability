package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// DeleteImage elimina una imagen de S3
func (s *S3Uploader) DeleteImage(ctx context.Context, filename string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		s.log.Error(ctx).Err(err).Str("filename", filename).Msg("error eliminando imagen de S3")
		return err
	}

	s.log.Info(ctx).Str("filename", filename).Msg("imagen eliminada exitosamente")
	return nil
}

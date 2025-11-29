package storage

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

// FileExists verifica si un archivo existe en S3 usando HeadObject
func (s *S3Uploader) FileExists(ctx context.Context, filename string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.bucket,
		Key:    &filename,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()
			if code == "NotFound" || code == "NoSuchKey" {
				return false, nil
			}
		}
		s.log.Error(ctx).Err(err).Str("filename", filename).Msg("error verificando existencia de archivo en S3")
		return false, err
	}
	return true, nil
}

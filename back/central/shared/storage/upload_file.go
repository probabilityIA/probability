package storage

import (
	"central/shared/errs"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// UploadFile mantiene la funcionalidad original para archivos generales
func (s *S3Uploader) UploadFile(ctx context.Context, file io.ReadSeeker, filename string) (string, error) {
	if file == nil {
		return "", errs.New("archivo es nulo")
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(s.bucket),
		Key:                  aws.String(filename),
		Body:                 file,
		ContentDisposition:   aws.String("inline"),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
		StorageClass:         types.StorageClassIntelligentTiering,
	})
	if err != nil {
		s.log.Error(ctx).Err(err).Msg("error subiendo archivo a S3")
		return "", err
	}

	url := s.GetImageURL(filename)
	return url, nil
}

package storage

import (
	"context"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/secamc93/probability/back/central/shared/errs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (s *S3Uploader) UploadFile(ctx context.Context, file io.ReadSeeker, filename string) (string, error) {
	if file == nil {
		return "", errs.New("archivo es nulo")
	}

	input := &s3.PutObjectInput{
		Bucket:             aws.String(s.bucket),
		Key:                aws.String(filename),
		Body:               file,
		ContentDisposition: aws.String("inline"),
		StorageClass:       types.StorageClassIntelligentTiering,
	}

	if tipo := contentTypePorExtension(filename); tipo != "" {
		input.ContentType = aws.String(tipo)
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		s.log.Error(ctx).Err(err).Msg("error subiendo archivo a S3")
		return "", err
	}

	url := s.GetImageURL(filename)
	return url, nil
}

func contentTypePorExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return ""
	}
	return mime.TypeByExtension(ext)
}

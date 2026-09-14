package storage

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

const defaultChatAttachmentsBucket = "probability-chat-attachments"

type IChatMediaStorage interface {
	Put(ctx context.Context, key string, body []byte, contentType string) error
	Delete(ctx context.Context, key string) error
	PresignGet(ctx context.Context, key, filename string, ttl time.Duration) (string, error)
}

type chatMediaStorage struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
	logger  log.ILogger
}

func NewChatMedia(environment env.IConfig, logger log.ILogger) IChatMediaStorage {
	key := environment.Get("S3_KEY")
	secret := environment.Get("S3_SECRET")
	region := environment.Get("S3_REGION")
	bucket := strings.TrimSpace(environment.Get("S3_CHAT_ATTACHMENTS_BUCKET"))
	if bucket == "" {
		bucket = defaultChatAttachmentsBucket
	}

	if key == "" || secret == "" {
		logger.Warn(context.Background()).Msg("Sin credenciales de S3: los adjuntos de chat quedan deshabilitados")
		return nil
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, "")),
	)
	if err != nil {
		logger.Warn(context.Background()).Err(err).Msg("No se pudo configurar S3 para adjuntos de chat")
		return nil
	}

	endpoint := environment.Get("S3_ENDPOINT")
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = region
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		}
	})

	return &chatMediaStorage{
		client:  client,
		presign: s3.NewPresignClient(client),
		bucket:  bucket,
		logger:  logger.WithModule("chat_media_storage"),
	}
}

func (s *chatMediaStorage) Put(ctx context.Context, key string, body []byte, contentType string) error {
	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(body),
		ContentLength: aws.Int64(int64(len(body))),
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}
	if _, err := s.client.PutObject(ctx, input); err != nil {
		s.logger.Error(ctx).Err(err).Str("key", key).Msg("Error guardando adjunto de chat en S3")
		return err
	}
	return nil
}

func (s *chatMediaStorage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *chatMediaStorage) PresignGet(ctx context.Context, key, filename string, ttl time.Duration) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	if name := safeDownloadName(filename); name != "" {
		input.ResponseContentDisposition = aws.String(fmt.Sprintf("inline; filename=\"%s\"", name))
	}
	request, err := s.presign.PresignGetObject(ctx, input, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", err
	}
	return request.URL, nil
}

func safeDownloadName(filename string) string {
	name := path.Base(strings.ReplaceAll(strings.TrimSpace(filename), "\\", "/"))
	if name == "." || name == "/" {
		return ""
	}
	cleaned := strings.Map(func(r rune) rune {
		if r < 32 || r == '"' || r > 126 {
			return '_'
		}
		return r
	}, name)
	if len(cleaned) > 150 {
		cleaned = cleaned[:150]
	}
	return cleaned
}

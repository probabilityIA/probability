package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/shared/errs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// UploadImage sube una imagen con validaciones y optimizaciones específicas
// Retorna el path relativo del archivo (ej: "avatars/1234567890_imagen.jpg")
func (s *S3Uploader) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	// Validar archivo
	if file == nil {
		return "", errs.New("archivo es nulo")
	}

	// Validar tamaño
	if file.Size > maxFileSize {
		return "", errs.New("archivo demasiado grande, máximo 10MB")
	}

	// Validar tipo de archivo
	contentType := file.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return "", errs.New("tipo de archivo no permitido, solo imágenes (jpeg, jpg, png, gif, webp)")
	}

	// Abrir archivo
	src, err := file.Open()
	if err != nil {
		s.log.Error(ctx).Err(err).Msg("error abriendo archivo")
		return "", err
	}
	defer src.Close()

	// Generar nombre único para el archivo
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s/%d_%s", folder, timestamp, file.Filename)

	// Limpiar nombre de archivo
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ToLower(filename)

	// Subir a S3
	// ServerSideEncryption removido: requiere KMS configurado, no compatible con MinIO local
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:             aws.String(s.bucket),
		Key:                aws.String(filename),
		Body:               src,
		ContentType:        aws.String(contentType),
		ContentDisposition: aws.String("inline"),
		StorageClass:       types.StorageClassStandard,             // Mejor para acceso frecuente
		CacheControl:       aws.String("public, max-age=31536000"), // Cache por 1 año
	})
	if err != nil {
		s.log.Error(ctx).Err(err).Msg("error subiendo imagen a S3")
		return "", err
	}

	// Retornar solo el path relativo, no la URL completa
	s.log.Info(ctx).Str("filename", filename).Msg("imagen subida exitosamente")
	return filename, nil
}

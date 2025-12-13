package storage

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// maskString oculta parte de una cadena para logging seguro
func maskString(s string) string {
	if s == "" {
		return "<empty>"
	}
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}

// Tipos de archivo permitidos para imágenes
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// Tamaño máximo de archivo (10MB)
const maxFileSize = 10 * 1024 * 1024

// S3Uploader estructura principal para operaciones S3
// Implementa la interfaz IS3Service del dominio
type S3Uploader struct {
	client *s3.Client
	bucket string
	log    log.ILogger
}

// IS3Service define las operaciones de S3
type IS3Service interface {
	GetImageURL(filename string) string
	DeleteImage(ctx context.Context, filename string) error
	ImageExists(ctx context.Context, filename string) (bool, error)
	UploadFile(ctx context.Context, file io.ReadSeeker, filename string) (string, error)
	DownloadFile(ctx context.Context, filename string) (io.ReadSeeker, error)
	FileExists(ctx context.Context, filename string) (bool, error)
	GetFileURL(ctx context.Context, filename string) (string, error)
	UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
}

// New crea una nueva instancia de S3Uploader y retorna la interfaz IS3Service
func New(env env.IConfig, logger log.ILogger) IS3Service {
	s3Key := env.Get("S3_KEY")
	s3Secret := env.Get("S3_SECRET")
	s3Region := env.Get("S3_REGION")
	s3Bucket := env.Get("S3_BUCKET")

	// Debug: verificar qué valores se están leyendo
	logger.Debug(context.Background()).
		Str("s3_key", maskString(s3Key)).
		Str("s3_secret", maskString(s3Secret)).
		Str("s3_region", s3Region).
		Str("s3_bucket", s3Bucket).
		Msg("S3 configuration loaded")

	// Validar que las credenciales no estén vacías
	if s3Key == "" || s3Secret == "" {
		logger.Fatal(context.Background()).
			Bool("has_key", s3Key != "").
			Bool("has_secret", s3Secret != "").
			Msg("❌ S3 credentials are empty - check S3_KEY and S3_SECRET environment variables")
		panic("S3 credentials are empty")
	}

	// Intentar conectar a S3
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3Key, s3Secret, "")),
	)
	if err != nil {
		logger.Fatal(context.Background()).
			Err(err).
			Msg("❌ No se pudo conectar a S3 - verifica las credenciales")
		panic("Error conectando a S3: " + err.Error())
	}

	// Configurar endpoint personalizado si está disponible (útil para MinIO en desarrollo)
	// Si S3_ENDPOINT está vacío, usa AWS S3 estándar (producción)
	s3Endpoint := env.Get("S3_ENDPOINT")

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.Region = s3Region
		// Si hay endpoint personalizado (MinIO), configurarlo
		if s3Endpoint != "" {
			o.BaseEndpoint = aws.String(s3Endpoint)
			o.UsePathStyle = true // Necesario para MinIO y S3-compatibles
		}
		// Si no hay endpoint, usa AWS S3 estándar (producción)
	})

	// Probar la conexión
	_, err = s3Client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: &s3Bucket,
	})
	if err != nil {
		logger.Warn(context.Background()).
			Err(err).
			Str("bucket", s3Bucket).
			Msg("⚠️ Bucket not found or not accessible. Attempting to create...")

		// Try to create the bucket
		_, errCreate := s3Client.CreateBucket(context.Background(), &s3.CreateBucketInput{
			Bucket: &s3Bucket,
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(s3Region),
			},
		})

		if errCreate != nil {
			logger.Fatal(context.Background()).
				Err(errCreate).
				Msg("❌ Failed to create S3 bucket. Please create it manually or check permissions.")
			panic("Error creating S3 bucket: " + errCreate.Error())
		}

		logger.Info(context.Background()).Str("bucket", s3Bucket).Msg("✅ S3 Bucket created successfully")
	}

	return &S3Uploader{
		client: s3Client,
		bucket: s3Bucket,
		log:    logger,
	}
}

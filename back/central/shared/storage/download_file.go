package storage

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// DownloadFile descarga un archivo de S3 y retorna un ReadSeeker en memoria
func (s *S3Uploader) DownloadFile(ctx context.Context, filename string) (io.ReadSeeker, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &filename,
	})
	if err != nil {
		s.log.Error(ctx).Err(err).Str("filename", filename).Msg("error descargando archivo de S3")
		return nil, err
	}
	defer out.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, out.Body); err != nil {
		s.log.Error(ctx).Err(err).Str("filename", filename).Msg("error leyendo cuerpo de S3")
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

package storage

import (
	"context"
	stderrors "errors"
	"io"
	"mime/multipart"
	"testing"

	"github.com/secamc93/probability/back/central/shared/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type s3Falso struct {
	UploadFileFn  func(ctx context.Context, file io.ReadSeeker, filename string) (string, error)
	DeleteImageFn func(ctx context.Context, filename string) error

	ClavesSubidas  []string
	ClavesBorradas []string
}

var _ storage.IS3Service = (*s3Falso)(nil)

func (s *s3Falso) GetImageURL(filename string) string { return "" }

func (s *s3Falso) DeleteImage(ctx context.Context, filename string) error {
	s.ClavesBorradas = append(s.ClavesBorradas, filename)
	if s.DeleteImageFn != nil {
		return s.DeleteImageFn(ctx, filename)
	}
	return nil
}

func (s *s3Falso) ImageExists(ctx context.Context, filename string) (bool, error) {
	return false, nil
}

func (s *s3Falso) UploadFile(ctx context.Context, file io.ReadSeeker, filename string) (string, error) {
	s.ClavesSubidas = append(s.ClavesSubidas, filename)
	if s.UploadFileFn != nil {
		return s.UploadFileFn(ctx, file, filename)
	}
	return "https://cdn/" + filename, nil
}

func (s *s3Falso) DownloadFile(ctx context.Context, filename string) (io.ReadSeeker, error) {
	return nil, nil
}

func (s *s3Falso) FileExists(ctx context.Context, filename string) (bool, error) { return false, nil }

func (s *s3Falso) GetFileURL(ctx context.Context, filename string) (string, error) { return "", nil }

func (s *s3Falso) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	return "", nil
}

func TestUploadFile_ArmaLaClaveComoCarpetaBarraArchivo(t *testing.T) {
	s3 := &s3Falso{}

	url, err := New(s3).UploadFile(context.Background(), "tickets/7", "abc.png", []byte("datos"), "image/png")

	require.NoError(t, err)
	assert.Equal(t, []string{"tickets/7/abc.png"}, s3.ClavesSubidas)
	assert.Equal(t, "https://cdn/tickets/7/abc.png", url)
}

func TestUploadFile_EntregaElContenidoIntacto(t *testing.T) {
	var recibido []byte
	s3 := &s3Falso{
		UploadFileFn: func(ctx context.Context, file io.ReadSeeker, filename string) (string, error) {
			b, err := io.ReadAll(file)
			require.NoError(t, err)
			recibido = b
			return "https://cdn/" + filename, nil
		},
	}

	_, err := New(s3).UploadFile(context.Background(), "tickets/7", "abc.png", []byte("bytes reales"), "image/png")

	require.NoError(t, err)
	assert.Equal(t, "bytes reales", string(recibido))
}

func TestUploadFile_ErrorDeS3_SePropaga(t *testing.T) {
	fallo := stderrors.New("access denied")
	s3 := &s3Falso{
		UploadFileFn: func(ctx context.Context, file io.ReadSeeker, filename string) (string, error) {
			return "", fallo
		},
	}

	url, err := New(s3).UploadFile(context.Background(), "tickets/7", "abc.png", nil, "image/png")

	assert.ErrorIs(t, err, fallo)
	assert.Empty(t, url)
}

func TestDeleteFile_BorraPorLaURLCompleta(t *testing.T) {
	s3 := &s3Falso{}

	err := New(s3).DeleteFile(context.Background(), "https://cdn/tickets/7/abc.png")

	require.NoError(t, err)
	assert.Equal(t, []string{"https://cdn/tickets/7/abc.png"}, s3.ClavesBorradas)
}

func TestDeleteFile_ErrorDeS3_SePropaga(t *testing.T) {
	fallo := stderrors.New("no such key")
	s3 := &s3Falso{
		DeleteImageFn: func(ctx context.Context, filename string) error { return fallo },
	}

	err := New(s3).DeleteFile(context.Background(), "https://cdn/x.png")

	assert.ErrorIs(t, err, fallo)
}

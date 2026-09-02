package app

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func archivoEnDiscoYaBorrado(t *testing.T, nombre string) *multipart.FileHeader {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+nombre+`"`)
	h.Set("Content-Type", "image/png")
	part, err := w.CreatePart(h)
	require.NoError(t, err)
	_, err = part.Write([]byte("contenido que va a parar a disco"))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	r := multipart.NewReader(&buf, w.Boundary())
	form, err := r.ReadForm(1)
	require.NoError(t, err)
	require.NoError(t, form.RemoveAll())
	return form.File["file"][0]
}

func TestUploadAttachment_ArchivoIlegible_NoSubeNiRegistra(t *testing.T) {
	storage := &mocks.StorageServiceMock{}
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, storage).UploadAttachment(context.Background(),
		7, nil, 42, archivoEnDiscoYaBorrado(t, "captura.png"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "open file")
	assert.Empty(t, storage.UploadedFolders, "no se sube nada si no se pudo abrir el archivo")
	assert.Nil(t, repo.AddedAttachment)
}

func TestUploadAttachment_ArchivoVacio_NoSubeNiRegistra(t *testing.T) {
	storage := &mocks.StorageServiceMock{}
	repo := &mocks.RepositoryMock{}

	_, err := newTicketsUseCase(repo, storage).UploadAttachment(context.Background(),
		7, nil, 42, archivoDePrueba(t, "vacio.png", "", "image/png"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read file")
	assert.Empty(t, storage.UploadedFolders)
	assert.Nil(t, repo.AddedAttachment, "un adjunto de cero bytes no llega a la base")
}

func TestUploadAttachment_SinExtension_NoInventaUna(t *testing.T) {
	var nombreEnS3 string
	storage := &mocks.StorageServiceMock{
		UploadFileFn: func(ctx context.Context, folder, filename string, data []byte, contentType string) (string, error) {
			nombreEnS3 = filename
			return "s3://" + filename, nil
		},
	}

	_, err := newTicketsUseCase(&mocks.RepositoryMock{}, storage).UploadAttachment(
		context.Background(), 7, nil, 42, archivoDePrueba(t, "informe", "datos", "text/plain"))

	require.NoError(t, err)
	assert.NotContains(t, nombreEnS3, ".")
}

func TestUploadAttachment_SubeElContenidoExacto(t *testing.T) {
	var subido []byte
	storage := &mocks.StorageServiceMock{
		UploadFileFn: func(ctx context.Context, folder, filename string, data []byte, contentType string) (string, error) {
			subido = append([]byte(nil), data...)
			return "s3://x.png", nil
		},
	}

	_, err := newTicketsUseCase(&mocks.RepositoryMock{}, storage).UploadAttachment(
		context.Background(), 7, nil, 42, archivoDePrueba(t, "captura.png", "bytes reales", "image/png"))

	require.NoError(t, err)
	assert.Equal(t, "bytes reales", string(subido))
}

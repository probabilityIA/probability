package storage

import (
	"bytes"
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	sharedstorage "github.com/secamc93/probability/back/central/shared/storage"
)

type PDFUploader struct {
	s3 sharedstorage.IS3Service
}

func New(s3 sharedstorage.IS3Service) domain.IPDFUploader {
	return &PDFUploader{s3: s3}
}

func (u *PDFUploader) UploadPDF(ctx context.Context, key string, content []byte) (string, error) {
	reader := bytes.NewReader(content)
	return u.s3.UploadFile(ctx, reader, key)
}

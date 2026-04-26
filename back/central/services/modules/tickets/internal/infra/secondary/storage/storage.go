package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/storage"
)

type Adapter struct {
	s3 storage.IS3Service
}

func New(s3 storage.IS3Service) ports.IStorageService {
	return &Adapter{s3: s3}
}

func (a *Adapter) UploadFile(ctx context.Context, folder, filename string, data []byte, _ string) (string, error) {
	key := fmt.Sprintf("%s/%s", folder, filename)
	reader := bytes.NewReader(data)
	return a.s3.UploadFile(ctx, reader, key)
}

func (a *Adapter) DeleteFile(ctx context.Context, fileURL string) error {
	return a.s3.DeleteImage(ctx, fileURL)
}

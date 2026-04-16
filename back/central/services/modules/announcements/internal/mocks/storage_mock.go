package mocks

import "context"

type StorageMock struct {
	UploadFileFn func(ctx context.Context, folder string, filename string, data []byte, contentType string) (string, error)
	DeleteFileFn func(ctx context.Context, fileURL string) error
}

func (m *StorageMock) UploadFile(ctx context.Context, folder string, filename string, data []byte, contentType string) (string, error) {
	if m.UploadFileFn != nil {
		return m.UploadFileFn(ctx, folder, filename, data, contentType)
	}
	return "", nil
}

func (m *StorageMock) DeleteFile(ctx context.Context, fileURL string) error {
	if m.DeleteFileFn != nil {
		return m.DeleteFileFn(ctx, fileURL)
	}
	return nil
}

package mocks

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

// S3Mock es un mock de domain.IS3Service usando testify/mock
type S3Mock struct {
	mock.Mock
}

func (m *S3Mock) GetImageURL(filename string) string {
	args := m.Called(filename)
	return args.String(0)
}

func (m *S3Mock) DeleteImage(ctx context.Context, filename string) error {
	args := m.Called(ctx, filename)
	return args.Error(0)
}

func (m *S3Mock) ImageExists(ctx context.Context, filename string) (bool, error) {
	args := m.Called(ctx, filename)
	return args.Bool(0), args.Error(1)
}

func (m *S3Mock) UploadFile(ctx context.Context, file io.ReadSeeker, filename string) (string, error) {
	args := m.Called(ctx, file, filename)
	return args.String(0), args.Error(1)
}

func (m *S3Mock) DownloadFile(ctx context.Context, filename string) (io.ReadSeeker, error) {
	args := m.Called(ctx, filename)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadSeeker), args.Error(1)
}

func (m *S3Mock) FileExists(ctx context.Context, filename string) (bool, error) {
	args := m.Called(ctx, filename)
	return args.Bool(0), args.Error(1)
}

func (m *S3Mock) GetFileURL(ctx context.Context, filename string) (string, error) {
	args := m.Called(ctx, filename)
	return args.String(0), args.Error(1)
}

func (m *S3Mock) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	args := m.Called(ctx, file, folder)
	return args.String(0), args.Error(1)
}

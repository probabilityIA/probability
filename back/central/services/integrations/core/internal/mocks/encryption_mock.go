package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// EncryptionMock es un mock de domain.IEncryptionService usando testify/mock
type EncryptionMock struct {
	mock.Mock
}

func (m *EncryptionMock) EncryptCredentials(ctx context.Context, credentials map[string]interface{}) ([]byte, error) {
	args := m.Called(ctx, credentials)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *EncryptionMock) DecryptCredentials(ctx context.Context, encryptedData []byte) (map[string]interface{}, error) {
	args := m.Called(ctx, encryptedData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *EncryptionMock) EncryptValue(ctx context.Context, value string) (string, error) {
	args := m.Called(ctx, value)
	return args.String(0), args.Error(1)
}

func (m *EncryptionMock) DecryptValue(ctx context.Context, encryptedValue string) (string, error) {
	args := m.Called(ctx, encryptedValue)
	return args.String(0), args.Error(1)
}

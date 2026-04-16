package mocks

import "github.com/stretchr/testify/mock"

// ConfigMock es un mock de env.IConfig usando testify/mock
type ConfigMock struct {
	mock.Mock
}

func (m *ConfigMock) Get(key string) string {
	args := m.Called(key)
	return args.String(0)
}

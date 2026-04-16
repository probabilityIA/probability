package mocks

// ConfigMock implementa env.IConfig para tests unitarios
type ConfigMock struct {
	Values map[string]string
}

func (m *ConfigMock) Get(key string) string {
	if m.Values != nil {
		return m.Values[key]
	}
	return ""
}

package mocks

// ConfigMock implementa env.IConfig para tests unitarios
type ConfigMock struct {
	Values map[string]string
}

func (c *ConfigMock) Get(key string) string {
	if c.Values != nil {
		return c.Values[key]
	}
	return ""
}

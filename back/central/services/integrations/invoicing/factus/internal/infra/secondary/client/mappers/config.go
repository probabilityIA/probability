package mappers

import "strconv"

// GetConfigString extrae un string del mapa de config con fallback al valor por defecto
func GetConfigString(config map[string]interface{}, key, defaultValue string) string {
	if config == nil {
		return defaultValue
	}
	val, ok := config[key]
	if !ok || val == nil {
		return defaultValue
	}
	switch v := val.(type) {
	case string:
		if v != "" {
			return v
		}
	case float64:
		return strconv.Itoa(int(v))
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	}
	return defaultValue
}

// GetConfigInt extrae un int del mapa de config con fallback al valor por defecto
func GetConfigInt(config map[string]interface{}, key string, defaultValue int) int {
	if config == nil {
		return defaultValue
	}
	val, ok := config[key]
	if !ok || val == nil {
		return defaultValue
	}
	switch v := val.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case string:
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultValue
}

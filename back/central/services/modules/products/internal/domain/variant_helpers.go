package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"gorm.io/datatypes"
)

// CanonicalizeVariantAttributes normaliza los atributos de variante para comparación y deduplicación.
func CanonicalizeVariantAttributes(raw datatypes.JSON) (datatypes.JSON, string, error) {
	if len(raw) == 0 {
		return nil, "", nil
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()

	var payload map[string]interface{}
	if err := decoder.Decode(&payload); err != nil {
		return nil, "", fmt.Errorf("variant_attributes must be a valid JSON object: %w", err)
	}

	if len(payload) == 0 {
		return datatypes.JSON([]byte("{}")), "", nil
	}

	normalized := normalizeMap(payload)
	canonicalBytes, err := json.Marshal(normalized)
	if err != nil {
		return nil, "", fmt.Errorf("failed to canonicalize variant_attributes: %w", err)
	}

	return datatypes.JSON(canonicalBytes), string(canonicalBytes), nil
}

func normalizeMap(payload map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0, len(payload))
	for key := range payload {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	normalized := make(map[string]interface{}, len(payload))
	for _, key := range keys {
		normalized[strings.TrimSpace(strings.ToLower(key))] = normalizeValue(payload[key])
	}

	return normalized
}

func normalizeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case map[string]interface{}:
		return normalizeMap(v)
	case []interface{}:
		normalized := make([]interface{}, len(v))
		for i, item := range v {
			normalized[i] = normalizeValue(item)
		}
		return normalized
	default:
		return v
	}
}

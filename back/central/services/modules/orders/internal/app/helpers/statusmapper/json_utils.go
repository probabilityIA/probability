package statusmapper

import "encoding/json"

// EqualJSON compara dos valores JSONB
func EqualJSON(a, b []byte) bool {
	var aMap, bMap map[string]interface{}
	if err := json.Unmarshal(a, &aMap); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &bMap); err != nil {
		return false
	}

	aBytes, _ := json.Marshal(aMap)
	bBytes, _ := json.Marshal(bMap)
	return string(aBytes) == string(bBytes)
}

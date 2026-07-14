package repository

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func marshalModuleCodes(codes []string) (datatypes.JSON, error) {
	raw, err := json.Marshal(codes)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(raw), nil
}

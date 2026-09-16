package businessdata

import (
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
)

type Reader struct {
	db db.IDatabase
}

func New(database db.IDatabase) ports.IBusinessDataReader {
	if database == nil {
		return nil
	}
	return &Reader{db: database}
}

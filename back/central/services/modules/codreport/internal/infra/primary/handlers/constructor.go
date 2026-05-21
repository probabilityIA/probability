package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Handlers struct {
	uc  app.Iapp
	log log.ILogger
}

func New(uc app.Iapp, logger log.ILogger) *Handlers {
	return &Handlers{uc: uc, log: logger}
}

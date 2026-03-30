package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/monitoring/internal/domain/ports"
)

type IHandler interface {
	RegisterRoutes(router *gin.Engine)
}

type handler struct {
	useCase   ports.IUseCase
	jwtSecret string
}

func New(useCase ports.IUseCase, jwtSecret string) IHandler {
	return &handler{
		useCase:   useCase,
		jwtSecret: jwtSecret,
	}
}

package repository

import (
	"github.com/secamc93/probability/back/monitoring/internal/domain/ports"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) ports.IUserRepository {
	return &UserRepository{db: db}
}

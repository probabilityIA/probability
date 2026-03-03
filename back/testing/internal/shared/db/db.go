package db

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/testing/shared/env"
	"github.com/secamc93/probability/back/testing/shared/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type IDatabase interface {
	Conn(ctx context.Context) *gorm.DB
	Close() error
}

type database struct {
	conn   *gorm.DB
	log    log.ILogger
	config env.IConfig
}

func New(logger log.ILogger, config env.IConfig) IDatabase {
	d := &database{log: logger, config: config}

	if err := d.connect(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	return d
}

func (d *database) connect() error {
	sslmode := d.config.GetWithDefault("PGSSLMODE", "require")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		d.config.Get("DB_HOST"),
		d.config.Get("DB_USER"),
		d.config.Get("DB_PASS"),
		d.config.Get("DB_NAME"),
		d.config.Get("DB_PORT"),
		sslmode,
	)

	var err error
	d.conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:              false,
		DisableNestedTransaction: true,
		Logger:                   logger.Default.LogMode(logger.Error),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}

	sqlDB, err := d.conn.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	d.log.Info().Msg("Database connected successfully (read-only mode)")
	return nil
}

func (d *database) Conn(ctx context.Context) *gorm.DB {
	return d.conn.WithContext(ctx)
}

func (d *database) Close() error {
	if d.conn != nil {
		sqlDB, err := d.conn.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

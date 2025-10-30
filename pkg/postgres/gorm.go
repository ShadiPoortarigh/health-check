package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnOptions struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Schema   string
}

func (cfg DBConnOptions) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.Schema)
}

func NewPsqlGormConnection(cfg DBConnOptions) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{
		Logger: logger.Discard,
	})
}

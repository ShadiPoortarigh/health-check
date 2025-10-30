package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"health-check/config"
)

func PostgresDSN(cfg config.DBConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.Schema)
}

func NewPsqlGormConnection(cfg config.DBConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(PostgresDSN(cfg)), &gorm.Config{
		Logger: logger.Discard,
	})
}

func SetDB(cfg config.Config) (*gorm.DB, error) {
	db, err := NewPsqlGormConnection(config.DBConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Database: cfg.DB.Database,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		Schema:   cfg.DB.Schema,
	})
	if err != nil {
		panic(err)
	}
	return db, nil
}

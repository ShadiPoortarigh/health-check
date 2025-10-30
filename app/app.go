package app

import (
	"gorm.io/gorm"
	"health-check/config"
	"health-check/pkg/postgres"
)

type app struct {
	db  *gorm.DB
	cfg config.Config
}

func (a *app) DB() *gorm.DB {
	return a.db
}

func (a *app) Config() config.Config {
	return a.cfg
}

func (a *app) SetDB() error {
	db, err := postgres.NewPsqlGormConnection(postgres.DBConnOptions{
		Host:     a.cfg.DB.Host,
		Port:     a.cfg.DB.Port,
		Database: a.cfg.DB.Database,
		Username: a.cfg.DB.Username,
		Password: a.cfg.DB.Password,
		Schema:   a.cfg.DB.Schema,
	})
	if err != nil {
		panic(err)
	}
	a.db = db
	return nil
}

func NewApp(cfg config.Config) (App, error) {
	a := &app{
		cfg: cfg,
	}
	if err := a.SetDB(); err != nil {
		return nil, err
	}
	return a, nil
}

func MustNewApp(cfg config.Config) App {
	a, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return a
}

package app

import (
	"gorm.io/gorm"
	"health-check/config"
)

type App interface {
	DB() *gorm.DB
	Config() config.Config
}

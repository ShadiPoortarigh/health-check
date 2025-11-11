package app

import (
	"gorm.io/gorm"
	"health-check/config"
	"health-check/internal/monitor_service/port"
)

type App interface {
	DB() *gorm.DB
	Config() config.Config
	HealthCheck() port.Service
}

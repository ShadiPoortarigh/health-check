package app

import (
	"context"
	"gorm.io/gorm"
	"health-check/config"
	"health-check/internal/monitor_service/port"
)

type App interface {
	DB() *gorm.DB
	Config() config.Config
	HealthCheck(ctx context.Context) port.Service
}

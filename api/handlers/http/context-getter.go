package http

import (
	"context"
	"health-check/api/service"
	"health-check/app"
)

type ContextGetter[T any] func(ctx context.Context) T

func SetContext(appContainer app.App) ContextGetter[*service.MonitorService] {

	return func(ctx context.Context) *service.MonitorService {

		return service.NewMonitorService(appContainer.HealthCheck(ctx))

	}
}

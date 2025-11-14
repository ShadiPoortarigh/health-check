package port

import (
	"context"
	"health-check/internal/common"
	"health-check/internal/monitor_service/domain"
)

type Service interface {
	RegisterApi(ctx context.Context, api domain.MonitoredAPI) (domain.ApiID, error)
	common.SchedulerHandler[domain.MonitoredAPI]
	ListAll(ctx context.Context) ([]domain.MonitoredAPI, error)
	DeleteApi(ctx context.Context, id domain.ApiID) error
}

package port

import (
	"context"
	"health-check/internal/monitor_service/domain"
)

type Repo interface {
	Create(ctx context.Context, api domain.MonitoredAPI) (domain.ApiID, error)
	GetByID(ctx context.Context, id domain.ApiID) (*domain.MonitoredAPI, error)
	SaveCheckResult(ctx context.Context, result domain.CheckResult) error
	ListAll(ctx context.Context) ([]domain.MonitoredAPI, error)
}

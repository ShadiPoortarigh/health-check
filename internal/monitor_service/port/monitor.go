package port

import (
	"context"
	"health-check/internal/monitor_service/domain"
)

type Repo interface {
	Create(ctx context.Context, api domain.MonitoredAPI) (domain.ApiID, error)
}

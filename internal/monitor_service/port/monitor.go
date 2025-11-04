package port

import "health-check/internal/monitor_service/domain"

type Repo interface {
	Create(api domain.MonitoredAPI) (domain.ApiID, error)
}

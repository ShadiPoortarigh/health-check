package port

import "health-check/internal/monitor_service/domain"

type Service interface {
	RegisterApi(api domain.MonitoredAPI) (domain.ApiID, error)
}

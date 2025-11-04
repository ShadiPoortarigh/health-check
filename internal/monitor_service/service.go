package monitor_service

import (
	"errors"
	"health-check/internal/monitor_service/domain"
	"health-check/internal/monitor_service/port"
)

type service struct {
	repo port.Repo
}

func NewService(repo port.Repo) port.Service {
	return &service{repo: repo}
}

func (s *service) RegisterApi(api domain.MonitoredAPI) (domain.ApiID, error) {
	if err := api.Validate(); err != nil {
		return domain.ApiID{}, err
	}
	if api.Interval == 0 {
		return domain.ApiID{}, errors.New("interval must be greater than zero")
	}
	return s.repo.Create(api)
}

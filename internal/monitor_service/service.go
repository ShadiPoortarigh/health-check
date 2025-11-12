package monitor_service

import (
	"context"
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

var (
	ErrAPIRegistrationValidation = errors.New("API validation failed")
)

func (s *service) RegisterApi(ctx context.Context, api domain.MonitoredAPI) (domain.ApiID, error) {
	if err := api.Validate(); err != nil {
		return domain.ApiID(0), err
	}
	if api.Interval == 0 {
		return domain.ApiID(0), errors.New("interval must be greater than zero")
	}
	return s.repo.Create(ctx, api)
}

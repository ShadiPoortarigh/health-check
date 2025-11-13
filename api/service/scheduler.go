package service

import (
	"context"
	"health-check/internal/common"
	"health-check/internal/monitor_service/domain"
	"health-check/internal/monitor_service/port"
	"log"
	"time"
)

type SchedulerService struct {
	svc             port.Service
	schedulerRunner *common.SchedulerRunner[domain.MonitoredAPI]
}

func NewSchedulerService(
	svc port.Service,
	schedulerRunner *common.SchedulerRunner[domain.MonitoredAPI],
) *SchedulerService {
	return &SchedulerService{
		svc:             svc,
		schedulerRunner: schedulerRunner,
	}
}

func (s *SchedulerService) Start(ctx context.Context, apiID domain.ApiID, duration time.Duration) error {

	api, err := s.svc.Query(ctx, apiID)
	if err != nil {
		return err
	}
	if api == nil {
		return nil
	}

	go s.schedulerRunner.RunFor(ctx, apiID, duration)

	log.Printf("Scheduler started for API %d for %.0f minutes", apiID, duration.Minutes())
	return nil
}

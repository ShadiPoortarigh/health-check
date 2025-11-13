package common

import (
	"context"
	"github.com/go-co-op/gocron"
	"health-check/internal/monitor_service/domain"
	"log"
	"time"
)

type SchedulerHandler[T any] interface {
	Handle(ctx context.Context, api *T) error
	Query(ctx context.Context, id domain.ApiID) (*T, error)
	Interval() time.Duration
}

type SchedulerRunner[T any] struct {
	handler   SchedulerHandler[T]
	scheduler *gocron.Scheduler
}

func NewSchedulerRunner[T any](handler SchedulerHandler[T]) *SchedulerRunner[T] {
	return &SchedulerRunner[T]{
		handler:   handler,
		scheduler: gocron.NewScheduler(time.UTC),
	}
}

func (s *SchedulerRunner[T]) RunFor(ctx context.Context, apiID domain.ApiID, duration time.Duration) {

	log.Printf("[Scheduler] Starting job for API %d every %v for %v minutes\n",
		apiID, s.handler.Interval(), duration.Minutes())

	ctx, cancel := context.WithTimeout(ctx, duration)

	var job *gocron.Job

	jobFunc := func() {
		select {
		case <-ctx.Done():
			log.Printf("[Scheduler] Stopping job for API %d\n", apiID)
			if job != nil {
				s.scheduler.RemoveByReference(job)
			}
			cancel()
			return
		default:
		}

		log.Printf("[Scheduler] Checking API %d at %s\n", apiID, time.Now().Format(time.RFC3339))

		api, err := s.handler.Query(ctx, apiID)
		if err != nil {
			log.Printf("[Scheduler] Failed to query API %d: %v\n", apiID, err)
			return
		}
		if api == nil {
			log.Printf("[Scheduler] API %d not found, stopping job\n", apiID)
			if job != nil {
				s.scheduler.RemoveByReference(job)
			}
			cancel()
			return
		}

		if err := s.handler.Handle(ctx, api); err != nil {
			log.Printf("[Scheduler] Failed to handle API %d: %v\n", apiID, err)
		}
	}

	var err error
	job, err = s.scheduler.Every(s.handler.Interval()).Do(jobFunc)
	if err != nil {
		log.Printf("[Scheduler] Failed to schedule job: %v\n", err)
		return
	}

	s.scheduler.StartAsync()
}

func (s *SchedulerRunner[T]) StopAll() {
	log.Println("[Scheduler] Stopping all jobs")
	s.scheduler.Clear()
}

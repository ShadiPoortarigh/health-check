package service

import (
	"context"
	"health-check/api/proto"
	monit "health-check/internal/monitor_service"
	"health-check/internal/monitor_service/domain"
	"health-check/internal/monitor_service/port"
	"time"
)

type MonitorService struct {
	svc port.Service
}

func NewMonitorService(svc port.Service) *MonitorService {
	return &MonitorService{svc: svc}
}

var (
	ErrAPIRegistrationValidation = monit.ErrAPIRegistrationValidation
)

func (u *MonitorService) RegisterAPI(ctx context.Context, req *proto.RegisterApiRequest) (*proto.RegisterApiResponse, error) {
	api := domain.MonitoredAPI{
		Name:     req.GetName(),
		URL:      req.GetUrl(),
		Method:   req.GetMethod(),
		Headers:  req.GetHeaders(),
		Body:     req.GetBody(),
		Interval: time.Duration(req.GetIntervalSeconds()) * time.Second,
		Enabled:  req.GetEnabled(),
		Webhook: domain.WebhookConfig{
			URL:     req.GetWebhook().GetUrl(),
			Headers: req.GetWebhook().GetHeaders(),
		},
	}
	id, err := u.svc.RegisterApi(api)
	if err != nil {
		return nil, err
	}
	return &proto.RegisterApiResponse{
		Id:              uint64(id),
		Url:             req.GetUrl(),
		Method:          req.GetMethod(),
		IntervalSeconds: req.GetIntervalSeconds(),
		Enabled:         req.GetEnabled(),
		CreatedAt:       time.Now().Format(time.RFC3339),
	}, nil
}

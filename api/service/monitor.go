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
	id, err := u.svc.RegisterApi(ctx, api)
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
func (m *MonitorService) Svc() port.Service {
	return m.svc
}

func (u *MonitorService) ListAPIs(ctx context.Context, req *proto.ListApisRequest) (*proto.ListApisResponse, error) {
	apis, err := u.svc.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := &proto.ListApisResponse{
		Apis: make([]*proto.MonitoredApi, 0, len(apis)),
	}

	for _, a := range apis {
		var lastCheckedAt string
		if a.LastCheckedAt != nil {
			lastCheckedAt = a.LastCheckedAt.Format(time.RFC3339)
		}

		resp.Apis = append(resp.Apis, &proto.MonitoredApi{
			Id:              uint64(a.ID),
			Name:            a.Name,
			Url:             a.URL,
			Method:          a.Method,
			Headers:         a.Headers,
			Body:            a.Body,
			IntervalSeconds: int64(a.Interval.Seconds()),
			Enabled:         a.Enabled,
			Webhook: &proto.Webhook{
				Url:     a.Webhook.URL,
				Headers: a.Webhook.Headers,
			},
			LastStatus:    a.LastStatus,
			LastCheckedAt: lastCheckedAt,
		})
	}

	return resp, nil
}

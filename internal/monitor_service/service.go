package monitor_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"health-check/internal/monitor_service/domain"
	"health-check/internal/monitor_service/port"
	"io"
	"net/http"
	"time"
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

func (s *service) Interval() time.Duration {
	return time.Second * 10
}

func (s *service) Query(ctx context.Context, id domain.ApiID) (*domain.MonitoredAPI, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Handle(ctx context.Context, api *domain.MonitoredAPI) error {
	start := time.Now()

	client := &http.Client{
		Timeout: 7 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
		},
	}
	// create request
	req, err := http.NewRequestWithContext(ctx, api.Method, api.URL, bytes.NewBufferString(api.Body))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	for k, v := range api.Headers {
		req.Header.Set(k, v)
	}

	// send request
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	checkResult := domain.CheckResult{
		ApiID:              api.ID,
		Timestamp:          time.Now(),
		ResponseTimeMillis: elapsed.Milliseconds(),
	}

	if err != nil {
		checkResult.Success = false
		checkResult.ErrorMessage = err.Error()
	} else {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		checkResult.StatusCode = resp.StatusCode
		checkResult.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
		checkResult.ResponseSnippet = string(bodyBytes)
	}

	if err := s.repo.SaveCheckResult(ctx, checkResult); err != nil {
		return fmt.Errorf("failed to save result: %w", err)
	}

	if !checkResult.Success && api.Webhook.URL != "" {
		return sendWebhookWithRetry(ctx, api.Webhook, checkResult)
	}

	return nil
}

func sendWebhookWithRetry(ctx context.Context, webhook domain.WebhookConfig, result domain.CheckResult) error {
	const maxRetries = 3
	const retryDelay = 2 * time.Second

	client := &http.Client{Timeout: 5 * time.Second}

	payload, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook.URL, bytes.NewReader(payload))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		for k, v := range webhook.Headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			return nil
		}

		if err != nil {
			fmt.Printf("Webhook attempt %d failed: %v\n", attempt, err)
		} else {
			fmt.Printf("Webhook attempt %d returned status %d\n", attempt, resp.StatusCode)
			resp.Body.Close()
		}

		time.Sleep(retryDelay)
	}

	return fmt.Errorf("failed to send webhook after %d retries", maxRetries)
}

func (s *service) ListAll(ctx context.Context) ([]domain.MonitoredAPI, error) {
	return s.repo.ListAll(ctx)
}

func (s *service) DeleteApi(ctx context.Context, id domain.ApiID) error {
	api, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if api == nil {
		return fmt.Errorf("api %d not found", id)
	}

	return s.repo.Delete(ctx, id)
}

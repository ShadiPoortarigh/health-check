package mapper

import (
	"encoding/json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"health-check/internal/monitor_service/domain"
	"health-check/pkg/adapters/types"
	"time"
)

func MonitorDomain2Storage(monitorDomain domain.MonitoredAPI) *types.MonitoredAPI {
	headersJSON, _ := json.Marshal(monitorDomain.Headers)
	webhookHeadersJSON, _ := json.Marshal(monitorDomain.Webhook.Headers)

	return &types.MonitoredAPI{
		Model: gorm.Model{
			ID:        uint(monitorDomain.ID),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
		URL:            monitorDomain.URL,
		Method:         monitorDomain.Method,
		Headers:        datatypes.JSON(headersJSON),
		Body:           monitorDomain.Body,
		Interval:       int64(monitorDomain.Interval.Seconds()),
		Enabled:        monitorDomain.Enabled,
		LastStatus:     nullableString(monitorDomain.LastStatus),
		LastCheckedAt:  monitorDomain.LastCheckedAt,
		WebhookURL:     monitorDomain.Webhook.URL,
		WebhookHeaders: datatypes.JSON(webhookHeadersJSON),
	}
}
func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func MonitorStorage2Domain(apiDB types.MonitoredAPI) (*domain.MonitoredAPI, error) {
	var headers, webhookHeaders map[string]string
	if err := json.Unmarshal(apiDB.Headers, &headers); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(apiDB.WebhookHeaders, &webhookHeaders); err != nil {
		return nil, err
	}

	var lastStatus string
	if apiDB.LastStatus != nil {
		lastStatus = *apiDB.LastStatus
	}

	return &domain.MonitoredAPI{
		ID:            domain.ApiID(apiDB.ID),
		URL:           apiDB.URL,
		Method:        apiDB.Method,
		Headers:       headers,
		Body:          apiDB.Body,
		Interval:      time.Duration(apiDB.Interval) * time.Second,
		Enabled:       apiDB.Enabled,
		LastStatus:    lastStatus,
		LastCheckedAt: apiDB.LastCheckedAt,
		Webhook: domain.WebhookConfig{
			URL:     apiDB.WebhookURL,
			Headers: webhookHeaders,
		},
	}, nil
}

func ResultDomain2Storage(result domain.CheckResult) *types.CheckResultDB {
	return &types.CheckResultDB{
		ApiID:              uint(result.ApiID),
		Timestamp:          result.Timestamp,
		StatusCode:         result.StatusCode,
		Success:            result.Success,
		ResponseTimeMillis: result.ResponseTimeMillis,
		ResponseSnippet:    result.ResponseSnippet,
		ErrorMessage:       result.ErrorMessage,
		CreatedAt:          time.Now(),
	}
}

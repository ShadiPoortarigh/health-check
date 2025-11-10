package mapper

import (
	"encoding/json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"health-check/internal/monitor_service/domain"
	"health-check/pkg/adapters/types"
	"time"
)

func MonitorDomain2Storage(monitorDomain domain.MonitoredAPI) *types.MonitoredAPIDB {
	headersJSON, _ := json.Marshal(monitorDomain.Headers)
	webhookHeadersJSON, _ := json.Marshal(monitorDomain.Webhook.Headers)

	return &types.MonitoredAPIDB{
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

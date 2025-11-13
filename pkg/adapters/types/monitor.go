package types

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type MonitoredAPIDB struct {
	gorm.Model
	URL            string         `gorm:"column:url" json:"url"`
	Method         string         `gorm:"column:method" json:"method"`
	Headers        datatypes.JSON `gorm:"column:headers" json:"headers"`
	Body           string         `gorm:"column:body" json:"body"`
	Interval       int64          `gorm:"column:interval_seconds" json:"interval"`
	Enabled        bool           `gorm:"column:enabled" json:"enabled"`
	LastStatus     *string        `gorm:"column:last_status" json:"last_status,omitempty"`
	LastCheckedAt  *time.Time     `gorm:"column:last_checked_at" json:"last_checked_at,omitempty"`
	WebhookURL     string         `gorm:"column:webhook_url" json:"webhook_url"`
	WebhookHeaders datatypes.JSON `gorm:"column:webhook_headers" json:"webhook_headers"`
}

type CheckResultDB struct {
	ID                 uint      `gorm:"primaryKey"`
	ApiID              uint      `gorm:"column:api_id"`
	Timestamp          time.Time `gorm:"column:timestamp"`
	StatusCode         int       `gorm:"column:status_code"`
	Success            bool      `gorm:"column:success"`
	ResponseTimeMillis int64     `gorm:"column:response_time_ms"`
	ResponseSnippet    string    `gorm:"column:response_snippet"`
	ErrorMessage       string    `gorm:"column:error_message"`
	CreatedAt          time.Time `gorm:"column:created_at"`
}

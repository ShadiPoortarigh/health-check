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

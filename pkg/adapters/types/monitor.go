package types

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type MonitoredAPIDB struct {
	gorm.Model
	URL            string         `db:"url" json:"url"`
	Method         string         `db:"method" json:"method"`
	Headers        datatypes.JSON `db:"headers" json:"headers"`
	Body           string         `db:"body" json:"body"`
	Interval       int64          `db:"interval_seconds" json:"interval"`
	Enabled        bool           `db:"enabled" json:"enabled"`
	LastStatus     *string        `db:"last_status" json:"last_status,omitempty"`
	LastCheckedAt  *time.Time     `db:"last_checked_at" json:"last_checked_at,omitempty"`
	WebhookURL     string         `db:"webhook_url" json:"webhook_url"`
	WebhookHeaders datatypes.JSON `db:"webhook_headers" json:"webhook_headers"`
}

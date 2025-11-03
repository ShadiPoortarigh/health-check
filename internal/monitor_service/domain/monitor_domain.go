package domain

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type (
	ApiID    string
	ResultID string
)

type MonitoredAPI struct {
	ID            ApiID             `json:"id"`                        // unique identifier
	Name          string            `json:"name,omitempty"`            // optional display name
	URL           string            `json:"url"`                       // target API URL
	Method        string            `json:"method"`                    // HTTP method (GET, POST, etc.)
	Headers       map[string]string `json:"headers,omitempty"`         // optional request headers
	Body          string            `json:"body,omitempty"`            // optional request body
	Interval      time.Duration     `json:"interval_seconds"`          // interval in seconds for health check
	Enabled       bool              `json:"enabled"`                   // whether the API is currently monitored
	Webhook       WebhookConfig     `json:"webhook"`                   // webhook configuration for notifications
	LastStatus    string            `json:"last_status,omitempty"`     // last known status (OK, FAIL, etc.)
	LastCheckedAt *time.Time        `json:"last_checked_at,omitempty"` // when the API was last checked
}

type WebhookConfig struct {
	URL     string            `json:"url"`               // webhook endpoint
	Headers map[string]string `json:"headers,omitempty"` // optional webhook headers
}

type CheckResult struct {
	ID                 ResultID  `json:"id"`                         // unique result ID
	ApiID              ApiID     `json:"api_id"`                     // related monitored API ID
	Timestamp          time.Time `json:"timestamp"`                  // time of the check
	StatusCode         int       `json:"status_code"`                // HTTP response code
	Success            bool      `json:"success"`                    // whether the check succeeded
	ResponseTimeMillis int64     `json:"response_time_ms"`           // response time in ms
	ResponseSnippet    string    `json:"response_snippet,omitempty"` // optional short response body
	ErrorMessage       string    `json:"error_message,omitempty"`    // optional error info
}

func (m *MonitoredAPI) Validate() error {
	if m.URL == "" {
		return errors.New("url cannot be empty")
	}

	parsedURL, err := url.ParseRequestURI(m.URL)
	if err != nil {
		return errors.New("invalid url format")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("url must start with http or https")
	}

	host := parsedURL.Hostname()

	if host == "localhost" || host == "0.0.0.0" || host == "127.0.0.1" {
		return fmt.Errorf("local or loopback addresses are not allowed: %s", host)
	}

	if strings.HasSuffix(host, ".internal") {
		return fmt.Errorf("internal hostnames are not allowed: %s", host)
	}

	ip := net.ParseIP(host)
	if ip == nil {

		ips, err := net.LookupIP(host)
		if err == nil {
			for _, resolvedIP := range ips {
				if isPrivateIP(resolvedIP) {
					return fmt.Errorf("private IP addresses are not allowed: %s", resolvedIP.String())
				}
			}
		}
	} else {
		if isPrivateIP(ip) {
			return fmt.Errorf("private IP addresses are not allowed: %s", ip.String())
		}
	}

	if m.Method == "" {
		return errors.New("http method is required")
	}

	method := strings.ToUpper(m.Method)
	validMethods := map[string]bool{"GET": true, "POST": true, "DELETE": true}
	if !validMethods[method] {
		return errors.New("invalid http method")
	}

	if m.Webhook.URL != "" {
		if _, err := url.ParseRequestURI(m.Webhook.URL); err != nil {
			return errors.New("invalid webhook url")
		}
	}

	return nil
}

func isPrivateIP(ip net.IP) bool {
	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",    // loopback
		"169.254.0.0/16", // link-local
	}

	for _, cidr := range privateBlocks {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

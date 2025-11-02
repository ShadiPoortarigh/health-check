package domain

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type MonitoredAPI struct {
	ID            string            `json:"id"`
	Name          string            `json:"name,omitempty"`
	URL           string            `json:"url"`
	Method        string            `json:"method"`
	Headers       map[string]string `json:"headers,omitempty"`
	Body          string            `json:"body,omitempty"`
	Interval      time.Duration     `json:"interval"`
	Enabled       bool              `json:"enabled"`
	Webhook       WebhookConfig     `json:"webhook"`
	LastStatus    string            `json:"last_status"`
	LastCheckedAt *time.Time        `json:"last_checked_at,omitempty"`
}

type WebhookConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}

type CheckResult struct {
	ID                 string    `json:"id"`
	ApiID              string    `json:"api_id"`
	Timestamp          time.Time `json:"timestamp"`
	StatusCode         int       `json:"status_code"`
	Success            bool      `json:"success"`
	ResponseTimeMillis int64     `json:"response_time_ms"`
	ResponseSnippet    string    `json:"response_snippet,omitempty"`
	ErrorMessage       string    `json:"error_message,omitempty"`
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

	// --- Block localhost, loopback, and internal hostnames ---
	if host == "localhost" || host == "0.0.0.0" || host == "127.0.0.1" {
		return fmt.Errorf("local or loopback addresses are not allowed: %s", host)
	}

	if strings.HasSuffix(host, ".internal") {
		return fmt.Errorf("internal hostnames are not allowed: %s", host)
	}

	// --- Resolve IPs and check private ranges ---
	ip := net.ParseIP(host)
	if ip == nil {
		// If host isn't an IP, try resolving DNS to check if it maps to a private IP
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

// --- Helper function to detect private/local IPs ---
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

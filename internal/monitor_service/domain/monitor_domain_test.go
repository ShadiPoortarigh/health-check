package domain

import "testing"

func TestMonitoredAPI_Validate(t *testing.T) {
	tests := []struct {
		name      string
		api       MonitoredAPI
		wantError bool
	}{
		{
			name: "valid GET API",
			api: MonitoredAPI{
				URL:    "https://example.com/health",
				Method: "GET",
			},
			wantError: false,
		},
		{
			name: "missing URL",
			api: MonitoredAPI{
				Method: "GET",
			},
			wantError: true,
		},
		{
			name: "invalid URL format",
			api: MonitoredAPI{
				URL:    "://bad-url",
				Method: "GET",
			},
			wantError: true,
		},
		{
			name: "missing HTTP method",
			api: MonitoredAPI{
				URL: "https://example.com",
			},
			wantError: true,
		},
		{
			name: "invalid HTTP method",
			api: MonitoredAPI{
				URL:    "https://example.com",
				Method: "FETCH",
			},
			wantError: true,
		},
		{
			name: "invalid webhook URL",
			api: MonitoredAPI{
				URL:     "https://example.com",
				Method:  "GET",
				Webhook: WebhookConfig{URL: "not-a-url"},
			},
			wantError: true,
		},
		// --- New validation cases ---
		{
			name: "rejects localhost",
			api: MonitoredAPI{
				URL:    "http://localhost/api",
				Method: "GET",
			},
			wantError: true,
		},
		{
			name: "rejects 127.0.0.1",
			api: MonitoredAPI{
				URL:    "http://127.0.0.1:8080/test",
				Method: "GET",
			},
			wantError: true,
		},
		{
			name: "rejects 0.0.0.0",
			api: MonitoredAPI{
				URL:    "http://0.0.0.0:8080/test",
				Method: "GET",
			},
			wantError: true,
		},
		{
			name: "rejects 192.168.x.x private IP",
			api: MonitoredAPI{
				URL:    "http://192.168.1.20/status",
				Method: "GET",
			},
			wantError: true,
		},
		{
			name: "rejects 10.x.x.x private IP",
			api: MonitoredAPI{
				URL:    "http://10.1.2.3/health",
				Method: "GET",
			},
			wantError: true,
		},
		{
			name: "rejects 172.16.x.x private IP",
			api: MonitoredAPI{
				URL:    "http://172.16.5.10/check",
				Method: "GET",
			},
			wantError: true,
		},
		{
			name: "rejects internal hostname",
			api: MonitoredAPI{
				URL:    "http://service.internal/api",
				Method: "GET",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.api.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("expected error=%v, got %v (err=%v)", tt.wantError, err != nil, err)
			}
		})
	}
}

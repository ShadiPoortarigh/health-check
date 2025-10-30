package app

import (
	"health-check/config"
	"os"
	"testing"
)

func makeValidConfig() config.Config {
	return config.Config{
		DB: config.DBConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "user",
			Password: "pass",
			Database: "testdb",
			Schema:   "public",
		},
		Server: config.ServerConfig{
			Port: 8080,
		},
	}
}

func TestNewApp(t *testing.T) {
	if os.Getenv("RUN_DB_TESTS") != "true" {
		t.Skip("Skipping DB integration test")
	}

	cfg := makeValidConfig()
	a, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if a.Config().DB.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", a.Config().DB.Host)
	}
}

func TestMustNewApp_PanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	var bad config.Config
	MustNewApp(bad)
}

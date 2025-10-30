package config

import (
	"os"
	"testing"
)

func TestReadConfig_ValidFile(t *testing.T) {
	jsonData := `{

			"db": {
				"host": "localhost",
				"port": 5432,
				"database": "testdb",
				"username": "user",
				"password": "pass",
				"schema": "public"
				},
			"server":{
				"port":8080
				}
	}`
	tmpFile := "test_config.json"
	_ = os.WriteFile(tmpFile, []byte(jsonData), 0644)
	defer os.Remove(tmpFile)

	cfg, err := ReadConfig(tmpFile)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.DB.Database != "testdb" {
		t.Errorf("expected database testdb, got %s", cfg.DB.Database)
	}
}

func TestReadConfig_InvalidFile(t *testing.T) {
	_, err := ReadConfig("nonexistent.json")

	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestReadConfig_InvalidJson(t *testing.T) {
	tmpFile := "bad.json"
	_ = os.WriteFile(tmpFile, []byte(`{invalid jason}`), 0644)
	defer os.Remove(tmpFile)

	_, err := ReadConfig(tmpFile)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

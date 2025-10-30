package postgres

import "testing"

func TestPostgresDSN(t *testing.T) {
	cfg := DBConnOptions{
		Host:     "localhost",
		Port:     5432,
		Username: "user",
		Password: "pass",
		Database: "testdb",
		Schema:   "public",
	}
	want := "host=localhost port=5432 user=user password=pass dbname=testdb search_path=public sslmode=disable"

	if got := cfg.PostgresDSN(); got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}

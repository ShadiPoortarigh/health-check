package storage

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"health-check/internal/monitor_service/domain"
	"testing"
	"time"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	connector := postgres.New(postgres.Config{
		Conn: db,
	})

	gdb, err := gorm.Open(connector, &gorm.Config{})
	assert.NoError(t, err)

	return gdb, mock
}

func TestMonitorRepo_Create_Success(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := NewDomainRepo(gdb)

	ctx := context.Background()

	api := domain.MonitoredAPI{
		ID:       0,
		URL:      "https://example.com",
		Method:   "GET",
		Headers:  map[string]string{"Authorization": "Bearer token"},
		Body:     "",
		Interval: 5 * time.Second,
		Enabled:  true,
		Webhook: domain.WebhookConfig{
			URL:     "https://webhook.site/test",
			Headers: map[string]string{"Content-Type": "application/json"},
		},
		LastStatus: "ok",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "monitored_apis"`).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			api.URL,
			api.Method,
			sqlmock.AnyArg(), // headers
			api.Body,
			int64(api.Interval.Seconds()),
			api.Enabled,
			sqlmock.AnyArg(), // last_status
			api.LastCheckedAt,
			api.Webhook.URL,
			sqlmock.AnyArg(), // webhook_headers
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	id, err := repo.Create(ctx, api)

	assert.NoError(t, err)
	assert.Equal(t, domain.ApiID(1), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMonitorRepo_Create_Failure(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := NewDomainRepo(gdb)

	ctx := context.Background()

	api := domain.MonitoredAPI{
		URL:      "https://example.com",
		Method:   "POST",
		Interval: 10 * time.Second,
		Enabled:  true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "monitored_apis"`).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	id, err := repo.Create(ctx, api)

	assert.Error(t, err)
	assert.Equal(t, domain.ApiID(0), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

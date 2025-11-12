package monitor_service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"health-check/internal/monitor_service/domain"
	"testing"
	"time"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(ctx context.Context, api domain.MonitoredAPI) (domain.ApiID, error) {
	args := m.Called(ctx, api)
	return args.Get(0).(domain.ApiID), args.Error(1)
}

func TestRegisterApi_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewService(mockRepo)

	ctx := context.Background()

	api := domain.MonitoredAPI{
		Name:     "Example API",
		URL:      "https://example.com",
		Method:   "GET",
		Interval: 10 * time.Second,
		Enabled:  true,
	}

	expectedID := domain.ApiID(1)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("domain.MonitoredAPI")).
		Return(expectedID, nil)

	id, err := service.RegisterApi(ctx, api)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
	mockRepo.AssertExpectations(t)
}

func TestRegisterApi_InvalidInterval(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewService(mockRepo)

	ctx := context.Background()

	api := domain.MonitoredAPI{
		Name:     "Invalid API",
		URL:      "https://example.com",
		Method:   "GET",
		Interval: 0,
	}

	id, err := service.RegisterApi(ctx, api)

	assert.Error(t, err)
	assert.Equal(t, domain.ApiID(0), id)
	assert.Contains(t, err.Error(), "interval must be greater than zero")
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRegisterApi_ValidationFails(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewService(mockRepo)

	ctx := context.Background()

	api := domain.MonitoredAPI{
		Name:     "",
		URL:      "",
		Method:   "GET",
		Interval: 10 * time.Second,
	}

	id, err := service.RegisterApi(ctx, api)

	assert.Error(t, err)
	assert.Equal(t, domain.ApiID(0), id)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRegisterApi_RepoError(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewService(mockRepo)

	ctx := context.Background()

	api := domain.MonitoredAPI{
		Name:     "Example API",
		URL:      "https://example.com",
		Method:   "GET",
		Interval: 10 * time.Second,
	}

	expectedError := errors.New("database unavailable")

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("domain.MonitoredAPI")).
		Return(domain.ApiID(0), expectedError)

	id, err := service.RegisterApi(ctx, api)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, domain.ApiID(0), id)
	mockRepo.AssertExpectations(t)
}

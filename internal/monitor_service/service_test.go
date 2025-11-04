package monitor_service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"health-check/internal/monitor_service/domain"
	"testing"
	"time"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(api domain.MonitoredAPI) (domain.ApiID, error) {
	args := m.Called(api)
	return args.Get(0).(domain.ApiID), args.Error(1)
}

// happy path
func TestRegisterApi_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewService(mockRepo)

	api := domain.MonitoredAPI{
		Name:     "Example API",
		URL:      "https://example.com",
		Method:   "GET",
		Interval: 10 * time.Second,
		Enabled:  true,
	}

	expectedID := domain.ApiID(uuid.New())

	mockRepo.On("Create", mock.AnythingOfType("domain.MonitoredAPI")).Return(expectedID, nil)

	id, err := service.RegisterApi(api)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
	mockRepo.AssertExpectations(t)

}

func TestRegisterApi_InvalidInterval(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewService(mockRepo)

	api := domain.MonitoredAPI{
		Name:     "Invalid url",
		URL:      "https://example.com",
		Method:   "GET",
		Interval: 0,
	}
	id, err := service.RegisterApi(api)

	assert.Error(t, err)
	assert.Equal(t, domain.ApiID{}, id)
	assert.Contains(t, err.Error(), "interval must be greater than zero")
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestRegisterApi_ValidationFails(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewService(mockRepo)

	api := domain.MonitoredAPI{
		Name:     "",
		URL:      "",
		Method:   "GET",
		Interval: 10 * time.Second,
	}

	id, err := service.RegisterApi(api)
	assert.Error(t, err)
	assert.Equal(t, domain.ApiID{}, id)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestRegisterApi_RepoError(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewService(mockRepo)

	api := domain.MonitoredAPI{
		Name:     "example api",
		URL:      "https://example.com",
		Method:   "GET",
		Interval: 10 * time.Second,
	}
	expectedError := errors.New("database unavailable")

	mockRepo.
		On("Create", mock.AnythingOfType("domain.MonitoredAPI")).
		Return(domain.ApiID{}, expectedError)

	id, err := service.RegisterApi(api)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, domain.ApiID{}, id)
	mockRepo.AssertExpectations(t)
}

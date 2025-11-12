package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	apihttp "health-check/api/handlers/http"
	"health-check/api/proto"
	"health-check/api/service"
	"health-check/internal/monitor_service/domain"
	"health-check/internal/monitor_service/port"
	"net/http/httptest"
	"testing"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) RegisterApi(ctx context.Context, api domain.MonitoredAPI) (domain.ApiID, error) {
	args := m.Called(ctx, api)
	return args.Get(0).(domain.ApiID), args.Error(1)
}

func setupTestApp(svc port.Service) *fiber.App {
	app := fiber.New()
	getSvc := func(ctx context.Context) *service.MonitorService {
		return service.NewMonitorService(svc)
	}
	app.Post("/api/v1/register", apihttp.RegisterAPI(getSvc))
	return app
}

func TestRegisterAPI_Success(t *testing.T) {
	mockSvc := new(mockService)
	app := setupTestApp(mockSvc)

	reqBody := proto.RegisterApiRequest{
		Name:            "User Service Health",
		Url:             "https://example.com/health",
		Method:          "GET",
		IntervalSeconds: 60,
		Enabled:         true,
		Webhook: &proto.Webhook{
			Url: "https://webhook.site/test",
			Headers: map[string]string{
				"Auth": "token",
			},
		},
	}

	mockSvc.On("RegisterApi", mock.Anything, mock.AnythingOfType("domain.MonitoredAPI")).
		Return(domain.ApiID(1), nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var parsed proto.RegisterApiResponse
	_ = json.NewDecoder(resp.Body).Decode(&parsed)
	assert.Equal(t, "https://example.com/health", parsed.Url)
	assert.Equal(t, int64(60), parsed.IntervalSeconds)
	assert.True(t, parsed.Enabled)

	mockSvc.AssertExpectations(t)
}

func TestRegisterAPI_BadRequest(t *testing.T) {
	mockSvc := new(mockService)
	app := setupTestApp(mockSvc)

	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer([]byte(`invalid-json`)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRegisterAPI_ValidationError(t *testing.T) {
	mockSvc := new(mockService)
	app := setupTestApp(mockSvc)

	reqBody := proto.RegisterApiRequest{
		Name: "Broken API",
	}

	mockSvc.On("RegisterApi", mock.Anything, mock.AnythingOfType("domain.MonitoredAPI")).
		Return(domain.ApiID(0), errors.New("interval must be greater than zero"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

package http

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"health-check/app"
	"health-check/config"
	"health-check/internal/common"
	"health-check/internal/monitor_service/domain"
)

func Run(appContainer app.App, cfg config.ServerConfig) error {
	router := fiber.New(fiber.Config{
		StrictRouting: false,
	})

	api := router.Group("/api/v1")
	getMonitorService := SetContext(appContainer)

	schedulerRunner := common.NewSchedulerRunner[domain.MonitoredAPI](
		appContainer.HealthCheck(context.Background()),
	)

	api.Post("/register", RegisterAPI(getMonitorService))
	api.Post("/:api_id/start", StartAPIHandler(getMonitorService, schedulerRunner))
	api.Get("/apis", ListAPIs(getMonitorService))
	api.Delete("/delete/:api_id", DeleteAPI(getMonitorService))

	return router.Listen(fmt.Sprintf(":%d", cfg.HttpPort))
}

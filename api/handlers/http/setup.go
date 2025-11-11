package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"health-check/api/service"
	"health-check/app"
	"health-check/config"
)

func Run(appContainer app.App, cfg config.ServerConfig) error {

	router := fiber.New()
	api := router.Group("/api/v1")

	api.Post("/register", RegisterAPI(service.NewMonitorService(appContainer.HealthCheck())))

	return router.Listen(fmt.Sprintf(":%d", cfg.HttpPort))
}

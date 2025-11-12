package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"health-check/app"
	"health-check/config"
)

func Run(appContainer app.App, cfg config.ServerConfig) error {

	router := fiber.New()
	api := router.Group("/api/v1")

	getMonitorService := SetContext(appContainer)

	api.Post("/register", RegisterAPI(getMonitorService))

	return router.Listen(fmt.Sprintf(":%d", cfg.HttpPort))
}

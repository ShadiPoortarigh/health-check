package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"health-check/api/proto"
	"health-check/api/service"
	"health-check/internal/common"
	"health-check/internal/monitor_service/domain"
	"strconv"
	"time"
)

func RegisterAPI(getSvc ContextGetter[*service.MonitorService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req proto.RegisterApiRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		ctx := c.UserContext()
		svc := getSvc(ctx)
		resp, err := svc.RegisterAPI(c.UserContext(), &req)
		if err != nil {
			if errors.Is(err, service.ErrAPIRegistrationValidation) {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.JSON(resp)
	}
}

func StartAPIHandler(
	getSvc ContextGetter[*service.MonitorService],
	schedulerRunner *common.SchedulerRunner[domain.MonitoredAPI],
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("api_id")
		apiIDUint, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid api id")
		}
		apiID := domain.ApiID(apiIDUint)

		duration := 5 * time.Minute
		if durParam := c.Query("duration"); durParam != "" {
			if durSeconds, err := strconv.Atoi(durParam); err == nil {
				duration = time.Duration(durSeconds) * time.Second
			}
		}

		ctx := c.UserContext()
		monitorSvc := getSvc(ctx)

		schedulerSvc := service.NewSchedulerService(monitorSvc.Svc(), schedulerRunner)

		if err := schedulerSvc.Start(ctx, apiID, duration); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"message":  "scheduler started",
			"api_id":   apiID,
			"duration": duration.Seconds(),
		})
	}
}

func ListAPIs(getSvc ContextGetter[*service.MonitorService]) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()
		svc := getSvc(ctx)

		req := &proto.ListApisRequest{}

		resp, err := svc.ListAPIs(ctx, req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(resp)
	}
}

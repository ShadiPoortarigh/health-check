package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"health-check/api/proto"
	"health-check/api/service"
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

package route

import (
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/http/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupConfigRoutes(r fiber.Router, h *handler.ConfigHandler, authMiddleware fiber.Handler) {
	group := r.Group("/configs", authMiddleware)

	group.Get("/", h.GetAllConfigs)
	group.Put("/:key", h.UpdateConfig)
}

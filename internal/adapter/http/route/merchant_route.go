package route

import (
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/http/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupMerchantRoutes(r fiber.Router, h *handler.MerchantHandler, authMiddleware fiber.Handler) {
	group := r.Group("/merchants", authMiddleware)

	group.Get("/", h.GetMerchants)
	group.Get("/:id", h.GetMerchantDetail)
	group.Post("/:id/suspend", h.SuspendMerchant)
	group.Post("/:id/activate", h.ActivateMerchant)
}

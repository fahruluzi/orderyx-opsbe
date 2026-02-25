package route

import (
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/http/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupDashboardRoutes(r fiber.Router, h *handler.DashboardHandler, authMiddleware fiber.Handler) {
	group := r.Group("/dashboard", authMiddleware)

	group.Get("/summary", h.GetSummary)
	group.Get("/growth", h.GetGrowth)
	group.Get("/revenue", h.GetRevenue)
}

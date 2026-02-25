package route

import (
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/http/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupAuditLogRoutes(r fiber.Router, h *handler.AuditLogHandler, authMiddleware fiber.Handler) {
	group := r.Group("/audit-logs", authMiddleware)
	group.Get("/", h.GetLogs)
}

package route

import (
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/http/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupPaymentRoutes(r fiber.Router, h *handler.PaymentHandler, authMiddleware fiber.Handler) {
	group := r.Group("/payments", authMiddleware)
	group.Get("/", h.GetPayments)
}

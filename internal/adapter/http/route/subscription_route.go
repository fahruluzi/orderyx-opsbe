package route

import (
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/http/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupSubscriptionRoutes(r fiber.Router, h *handler.SubscriptionHandler, authMiddleware fiber.Handler) {
	group := r.Group("/subscriptions", authMiddleware)

	group.Get("/:merchant_id/history", h.GetMerchantSubscriptions)
	group.Put("/:merchant_id/extend-trial", h.ExtendTrial)
	group.Put("/:merchant_id/change-plan", h.ChangePlan)
}

package route

import (
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/http/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(r fiber.Router, authHandler *handler.AuthHandler, authMiddleware fiber.Handler) {
	group := r.Group("/auth")

	group.Post("/login", authHandler.Login)
	group.Post("/logout", authMiddleware, authHandler.Logout)
	group.Get("/me", authMiddleware, authHandler.Me)
}

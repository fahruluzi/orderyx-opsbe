package handler

import (
	"github.com/fahruluzi/orderyx-opsbe/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	uc usecase.DashboardUsecase
}

func NewDashboardHandler(uc usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{uc: uc}
}

func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	summary, err := h.uc.GetSummary(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": summary})
}

func (h *DashboardHandler) GetGrowth(c *fiber.Ctx) error {
	days := c.QueryInt("days", 30) // Default 30 days
	growth, err := h.uc.GetGrowth(c.Context(), days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": growth})
}

func (h *DashboardHandler) GetRevenue(c *fiber.Ctx) error {
	days := c.QueryInt("days", 30)
	revenue, err := h.uc.GetRevenue(c.Context(), days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": revenue})
}

package handler

import (
	"github.com/fahruluzi/orderyx-opsbe/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	uc usecase.PaymentUsecase
}

func NewPaymentHandler(uc usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

func (h *PaymentHandler) GetPayments(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	payments, err := h.uc.GetPayments(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(payments)
}

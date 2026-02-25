package handler

import (
	"github.com/fahruluzi/orderyx-opsbe/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type AuditLogHandler struct {
	uc usecase.AuditLogUsecase
}

func NewAuditLogHandler(uc usecase.AuditLogUsecase) *AuditLogHandler {
	return &AuditLogHandler{uc: uc}
}

func (h *AuditLogHandler) GetLogs(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	logs, err := h.uc.GetLogs(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(logs)
}

package handler

import (
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/dto"
	"github.com/fahruluzi/orderyx-opsbe/internal/usecase"
	"github.com/fahruluzi/orderyx-opsbe/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

type ConfigHandler struct {
	uc usecase.ConfigUsecase
}

func NewConfigHandler(uc usecase.ConfigUsecase) *ConfigHandler {
	return &ConfigHandler{uc: uc}
}

func (h *ConfigHandler) GetAllConfigs(c *fiber.Ctx) error {
	configs, err := h.uc.GetAllConfigs(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": configs})
}

func (h *ConfigHandler) UpdateConfig(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(*jwt.OpsClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Only Super Admin can change system configs
	if actor.Role != "super_admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
	}

	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing config key"})
	}

	var req dto.UpdateConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Value == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "value is required"})
	}

	err := h.uc.UpdateConfig(c.Context(), actor, key, req.Value)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Configuration updated successfully"})
}

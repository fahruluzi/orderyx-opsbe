package handler

import (
	"strconv"

	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/dto"
	"github.com/fahruluzi/orderyx-opsbe/internal/usecase"
	"github.com/fahruluzi/orderyx-opsbe/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

type SubscriptionHandler struct {
	uc usecase.SubscriptionUsecase
}

func NewSubscriptionHandler(uc usecase.SubscriptionUsecase) *SubscriptionHandler {
	return &SubscriptionHandler{uc: uc}
}

func (h *SubscriptionHandler) GetMerchantSubscriptions(c *fiber.Ctx) error {
	merchantID, err := strconv.ParseInt(c.Params("merchant_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid merchant id"})
	}

	subs, err := h.uc.GetMerchantSubscriptions(c.Context(), merchantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": subs})
}

func (h *SubscriptionHandler) ExtendTrial(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(*jwt.OpsClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	merchantID, err := strconv.ParseInt(c.Params("merchant_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid merchant id"})
	}

	var req dto.ExtendTrialRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err = h.uc.ExtendTrial(c.Context(), actor, merchantID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Trial extended successfully"})
}

func (h *SubscriptionHandler) ChangePlan(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(*jwt.OpsClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	merchantID, err := strconv.ParseInt(c.Params("merchant_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid merchant id"})
	}

	var req dto.ChangePlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err = h.uc.ChangePlan(c.Context(), actor, merchantID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Subscription plan changed successfully"})
}

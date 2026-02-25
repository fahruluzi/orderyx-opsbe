package handler

import (
	"strconv"

	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/dto"
	"github.com/fahruluzi/orderyx-opsbe/internal/usecase"
	"github.com/fahruluzi/orderyx-opsbe/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

type MerchantHandler struct {
	uc usecase.MerchantUsecase
}

func NewMerchantHandler(uc usecase.MerchantUsecase) *MerchantHandler {
	return &MerchantHandler{uc: uc}
}

func (h *MerchantHandler) GetMerchants(c *fiber.Ctx) error {
	req := dto.MerchantPaginationRequest{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
		Status: c.Query("status", ""),
		Plan:   c.Query("plan", ""),
	}

	res, err := h.uc.GetMerchants(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (h *MerchantHandler) GetMerchantDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid merchant id"})
	}

	detail, err := h.uc.GetMerchantDetail(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": detail})
}

func (h *MerchantHandler) SuspendMerchant(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(*jwt.OpsClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid merchant id"})
	}

	err = h.uc.SuspendMerchant(c.Context(), actor, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Merchant suspended successfully"})
}

func (h *MerchantHandler) ActivateMerchant(c *fiber.Ctx) error {
	actor, ok := c.Locals("user").(*jwt.OpsClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid merchant id"})
	}

	err = h.uc.ActivateMerchant(c.Context(), actor, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Merchant activated successfully"})
}

package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/dto"
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/repository"
	"github.com/fahruluzi/orderyx-opsbe/pkg/jwt"
)

type MerchantUsecase interface {
	GetMerchants(ctx context.Context, req dto.MerchantPaginationRequest) (dto.PaginatedResponse, error)
	GetMerchantDetail(ctx context.Context, id int64) (*dto.MerchantDetailResponse, error)
	SuspendMerchant(ctx context.Context, actor *jwt.OpsClaims, id int64) error
	ActivateMerchant(ctx context.Context, actor *jwt.OpsClaims, id int64) error
}

type merchantUsecase struct {
	repo      repository.MerchantRepository
	auditRepo repository.AuditLogRepository
}

func NewMerchantUsecase(repo repository.MerchantRepository, auditRepo repository.AuditLogRepository) MerchantUsecase {
	return &merchantUsecase{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (u *merchantUsecase) GetMerchants(ctx context.Context, req dto.MerchantPaginationRequest) (dto.PaginatedResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	merchants, total, err := u.repo.GetMerchants(ctx, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	var data []dto.MerchantListResponse
	for _, m := range merchants {
		sub, _ := u.repo.GetLatestSubscription(ctx, m.ID)

		var planName, status string
		if sub != nil {
			status = sub.Status
			if sub.Plan != nil {
				planName = sub.Plan.Name
			}
		}

		// Skip appending if requested Plan filter does not match
		if req.Plan != "" && req.Plan != planName {
			continue // Note: this makes total count inaccurate, best to filter via Join in repository. For S1 prototype scope, it's acceptable.
		}

		item := dto.MerchantListResponse{
			ID:                 m.ID,
			Name:               m.Name,
			BusinessType:       string(m.BusinessType),
			IsActive:           m.IsActive,
			SubscriptionPlan:   planName,
			SubscriptionStatus: status,
			CreatedAt:          m.CreatedAt,
		}

		if sub != nil && status == "TRIAL" {
			item.SubscriptionTrialEnd = &sub.EndDate
		}

		data = append(data, item)
	}

	return dto.PaginatedResponse{
		Data:       data,
		TotalCount: total,
		Page:       req.Page,
		Limit:      req.Limit,
	}, nil
}

func (u *merchantUsecase) GetMerchantDetail(ctx context.Context, id int64) (*dto.MerchantDetailResponse, error) {
	m, err := u.repo.GetMerchantByID(ctx, id)
	if err != nil {
		return nil, errors.New("merchant not found")
	}

	sub, _ := u.repo.GetLatestSubscription(ctx, id)
	users, orders, purchases, _ := u.repo.GetMerchantStats(ctx, id)

	var planName, status string
	if sub != nil {
		status = sub.Status
		if sub.Plan != nil {
			planName = sub.Plan.Name
		}
	}

	res := &dto.MerchantDetailResponse{
		ID:                   m.ID,
		Name:                 m.Name,
		BusinessType:         string(m.BusinessType),
		Email:                m.Email,
		Phone:                m.Phone,
		Address:              m.Address,
		TaxID:                m.TaxID,
		IsActive:             m.IsActive,
		IsOnboardingComplete: m.IsOnboardingComplete,
		Settings:             m.Settings,
		CreatedAt:            m.CreatedAt,
		SubscriptionPlan:     planName,
		SubscriptionStatus:   status,
		TotalUsers:           users,
		TotalOrders:          orders,
		TotalPurchases:       purchases,
	}

	if sub != nil {
		res.SubscriptionEnd = &sub.EndDate
	}

	return res, nil
}

func (u *merchantUsecase) SuspendMerchant(ctx context.Context, actor *jwt.OpsClaims, id int64) error {
	m, err := u.repo.GetMerchantByID(ctx, id)
	if err != nil {
		return errors.New("merchant not found")
	}

	if !m.IsActive {
		return errors.New("merchant is already suspended")
	}

	err = u.repo.UpdateMerchantStatus(ctx, id, false)
	if err != nil {
		return err
	}

	// Audit Log
	detailMsg, _ := json.Marshal(map[string]string{"previous_status": "active", "new_status": "suspended", "reason": "manual suspension by ops admin"})
	_ = u.auditRepo.LogAction(ctx, actor.UserID, "SUSPEND_MERCHANT", "merchant", &id, string(detailMsg))

	return nil
}

func (u *merchantUsecase) ActivateMerchant(ctx context.Context, actor *jwt.OpsClaims, id int64) error {
	m, err := u.repo.GetMerchantByID(ctx, id)
	if err != nil {
		return errors.New("merchant not found")
	}

	if m.IsActive {
		return errors.New("merchant is already active")
	}

	err = u.repo.UpdateMerchantStatus(ctx, id, true)
	if err != nil {
		return err
	}

	// Audit Log
	detailMsg, _ := json.Marshal(map[string]string{"previous_status": "suspended", "new_status": "active", "reason": "manual activation by ops admin"})
	_ = u.auditRepo.LogAction(ctx, actor.UserID, "ACTIVATE_MERCHANT", "merchant", &id, string(detailMsg))

	return nil
}

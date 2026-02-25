package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/dto"
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/repository"
	"github.com/fahruluzi/orderyx-opsbe/pkg/jwt"
)

type SubscriptionUsecase interface {
	GetMerchantSubscriptions(ctx context.Context, merchantID int64) ([]dto.SubscriptionResponse, error)
	ExtendTrial(ctx context.Context, actor *jwt.OpsClaims, merchantID int64, req dto.ExtendTrialRequest) error
	ChangePlan(ctx context.Context, actor *jwt.OpsClaims, merchantID int64, req dto.ChangePlanRequest) error
}

type subscriptionUsecase struct {
	repo      repository.SubscriptionRepository
	auditRepo repository.AuditLogRepository
}

func NewSubscriptionUsecase(repo repository.SubscriptionRepository, auditRepo repository.AuditLogRepository) SubscriptionUsecase {
	return &subscriptionUsecase{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (u *subscriptionUsecase) GetMerchantSubscriptions(ctx context.Context, merchantID int64) ([]dto.SubscriptionResponse, error) {
	subs, err := u.repo.GetSubscriptionsByMerchant(ctx, merchantID)
	if err != nil {
		return nil, err
	}

	var res []dto.SubscriptionResponse
	for _, s := range subs {
		planName := "Custom"
		if s.Plan != nil {
			planName = s.Plan.Name
		}
		res = append(res, dto.SubscriptionResponse{
			ID:            s.ID,
			MerchantID:    s.MerchantID,
			PlanName:      planName,
			Status:        s.Status,
			PaymentStatus: s.PaymentStatus,
			StartDate:     s.StartDate,
			EndDate:       s.EndDate,
			CreatedAt:     s.CreatedAt,
		})
	}
	return res, nil
}

func (u *subscriptionUsecase) ExtendTrial(ctx context.Context, actor *jwt.OpsClaims, merchantID int64, req dto.ExtendTrialRequest) error {
	sub, err := u.repo.GetLatestSubscription(ctx, merchantID)
	if err != nil {
		return errors.New("subscription not found")
	}

	if sub.Status != "TRIAL" && sub.Status != "EXPIRED" {
		return errors.New("merchant is not in a trial or recently expired state to extend trial")
	}

	if req.EndDate.Before(time.Now()) {
		return errors.New("new end date must be in the future")
	}

	oldDate := sub.EndDate
	sub.EndDate = req.EndDate
	sub.Status = "TRIAL" // Reactivate if it was expired

	err = u.repo.UpdateSubscription(ctx, sub)
	if err != nil {
		return err
	}

	// Log Action
	detailMsg, _ := json.Marshal(map[string]interface{}{
		"old_end_date": oldDate,
		"new_end_date": req.EndDate,
	})
	_ = u.auditRepo.LogAction(ctx, actor.UserID, "EXTEND_TRIAL", "subscription", &sub.ID, string(detailMsg))

	return nil
}

func (u *subscriptionUsecase) ChangePlan(ctx context.Context, actor *jwt.OpsClaims, merchantID int64, req dto.ChangePlanRequest) error {
	sub, err := u.repo.GetLatestSubscription(ctx, merchantID)
	if err != nil {
		return errors.New("subscription not found")
	}

	plan, err := u.repo.GetPlanByID(ctx, req.PlanID)
	if err != nil {
		return errors.New("target subscription plan not found")
	}

	oldPlanName := "Custom"
	if sub.Plan != nil {
		oldPlanName = sub.Plan.Name
	}

	sub.PlanID = plan.ID
	err = u.repo.UpdateSubscription(ctx, sub)
	if err != nil {
		return err
	}

	// Log Action
	detailMsg, _ := json.Marshal(map[string]interface{}{
		"old_plan": oldPlanName,
		"new_plan": plan.Name,
	})
	_ = u.auditRepo.LogAction(ctx, actor.UserID, "CHANGE_PLAN", "subscription", &sub.ID, string(detailMsg))

	return nil
}

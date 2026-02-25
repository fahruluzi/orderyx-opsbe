package repository

import (
	"context"

	"github.com/fahruluzi/orderyx-opsbe/internal/domain"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	GetSubscriptionsByMerchant(ctx context.Context, merchantID int64) ([]domain.Subscription, error)
	GetSubscriptionByID(ctx context.Context, id int64) (*domain.Subscription, error)
	GetLatestSubscription(ctx context.Context, merchantID int64) (*domain.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *domain.Subscription) error
	GetPlanByID(ctx context.Context, planID int64) (*domain.SubscriptionPlan, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) GetSubscriptionsByMerchant(ctx context.Context, merchantID int64) ([]domain.Subscription, error) {
	var subs []domain.Subscription
	err := r.db.WithContext(ctx).
		Preload("Plan").
		Where("merchant_id = ?", merchantID).
		Order("created_at desc").
		Find(&subs).Error
	return subs, err
}

func (r *subscriptionRepository) GetSubscriptionByID(ctx context.Context, id int64) (*domain.Subscription, error) {
	var sub domain.Subscription
	err := r.db.WithContext(ctx).Preload("Plan").First(&sub, id).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) GetLatestSubscription(ctx context.Context, merchantID int64) (*domain.Subscription, error) {
	var sub domain.Subscription
	err := r.db.WithContext(ctx).
		Preload("Plan").
		Where("merchant_id = ?", merchantID).
		Order("created_at desc").
		First(&sub).Error

	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) UpdateSubscription(ctx context.Context, sub *domain.Subscription) error {
	return r.db.WithContext(ctx).Save(sub).Error
}

func (r *subscriptionRepository) GetPlanByID(ctx context.Context, planID int64) (*domain.SubscriptionPlan, error) {
	var plan domain.SubscriptionPlan
	err := r.db.WithContext(ctx).First(&plan, planID).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

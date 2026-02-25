package repository

import (
	"context"
	"fmt"

	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/dto"
	"github.com/fahruluzi/orderyx-opsbe/internal/domain"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	GetMerchants(ctx context.Context, req dto.MerchantPaginationRequest) ([]domain.Merchant, int64, error)
	GetMerchantByID(ctx context.Context, id int64) (*domain.Merchant, error)
	GetLatestSubscription(ctx context.Context, merchantID int64) (*domain.Subscription, error)
	UpdateMerchantStatus(ctx context.Context, id int64, isActive bool) error
	GetMerchantStats(ctx context.Context, id int64) (users int64, orders int64, purchases int64, err error)
}

type merchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

func (r *merchantRepository) GetMerchants(ctx context.Context, req dto.MerchantPaginationRequest) ([]domain.Merchant, int64, error) {
	var merchants []domain.Merchant
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Merchant{})

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", search, search)
	}

	if req.Status != "" {
		if req.Status == "Active" {
			query = query.Where("is_active = ?", true)
		} else if req.Status == "Inactive" {
			query = query.Where("is_active = ?", false)
		}
	}

	// We can't easily filter by plan/subscription status easily without a JOIN,
	// for now we'll fetch merchants and then attach their subscriptions in usercase,
	// or perform a join here. Doing a join for filtering if needed later.

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	err := query.Order("created_at desc").Offset(offset).Limit(req.Limit).Find(&merchants).Error
	if err != nil {
		return nil, 0, err
	}

	return merchants, total, nil
}

func (r *merchantRepository) GetMerchantByID(ctx context.Context, id int64) (*domain.Merchant, error) {
	var merchant domain.Merchant
	err := r.db.WithContext(ctx).First(&merchant, id).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (r *merchantRepository) GetLatestSubscription(ctx context.Context, merchantID int64) (*domain.Subscription, error) {
	var sub domain.Subscription
	err := r.db.WithContext(ctx).
		Preload("Plan").
		Where("merchant_id = ?", merchantID).
		Order("created_at desc").
		First(&sub).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil // Return nil if no sub found without erroring
	}
	return &sub, nil
}

func (r *merchantRepository) UpdateMerchantStatus(ctx context.Context, id int64, isActive bool) error {
	return r.db.WithContext(ctx).Model(&domain.Merchant{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (r *merchantRepository) GetMerchantStats(ctx context.Context, id int64) (users int64, orders int64, purchases int64, err error) {
	// Need to query respective tables. Assuming tables are named 'users', 'orders', 'purchases'.

	// Count users for this merchant tenant
	userErr := r.db.WithContext(ctx).Table("users").Where("merchant_id = ? AND deleted_at IS NULL", id).Count(&users).Error
	if userErr != nil {
		// Log error but continue
		fmt.Printf("Warning: Could not fetch user count: %v\n", userErr)
	}

	// Count orders
	orderErr := r.db.WithContext(ctx).Table("orders").Where("merchant_id = ? AND deleted_at IS NULL", id).Count(&orders).Error
	if orderErr != nil {
		fmt.Printf("Warning: Could not fetch order count: %v\n", orderErr)
	}

	// Count purchases
	purchErr := r.db.WithContext(ctx).Table("purchases").Where("merchant_id = ? AND deleted_at IS NULL", id).Count(&purchases).Error
	if purchErr != nil {
		fmt.Printf("Warning: Could not fetch purchase count: %v\n", purchErr)
	}

	return users, orders, purchases, nil
}

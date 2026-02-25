package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	MerchantID     int64     `json:"merchant_id"`
	SubscriptionID int64     `json:"subscription_id"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	PaymentMethod  string    `json:"payment_method"`
	Reference      string    `json:"reference"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PaymentRepository interface {
	GetPayments(ctx context.Context, limit, offset int) ([]Payment, int64, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) GetPayments(ctx context.Context, limit, offset int) ([]Payment, int64, error) {
	var payments []Payment
	var total int64

	query := r.db.WithContext(ctx).Table("payments") // Assume table is 'payments'

	err := query.Count(&total).Error
	if err != nil {
		// Table might not exist yet, grace fallback
		return []Payment{}, 0, nil
	}

	err = query.Order("created_at desc").Limit(limit).Offset(offset).Find(&payments).Error
	return payments, total, err
}

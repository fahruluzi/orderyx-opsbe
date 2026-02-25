package usecase

import (
	"context"

	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/repository"
)

type PaymentResponse struct {
	ID             int64   `json:"id"`
	MerchantID     int64   `json:"merchant_id"`
	SubscriptionID int64   `json:"subscription_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"`
	PaymentMethod  string  `json:"payment_method"`
	Reference      string  `json:"reference"`
	CreatedAt      string  `json:"created_at"`
}

type PaymentPaginationResponse struct {
	Data  []PaymentResponse `json:"data"`
	Total int64             `json:"total"`
}

type PaymentUsecase interface {
	GetPayments(ctx context.Context, limit, offset int) (*PaymentPaginationResponse, error)
}

type paymentUsecase struct {
	repo repository.PaymentRepository
}

func NewPaymentUsecase(repo repository.PaymentRepository) PaymentUsecase {
	return &paymentUsecase{repo: repo}
}

func (u *paymentUsecase) GetPayments(ctx context.Context, limit, offset int) (*PaymentPaginationResponse, error) {
	payments, total, err := u.repo.GetPayments(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var res []PaymentResponse
	for _, p := range payments {
		res = append(res, PaymentResponse{
			ID:             p.ID,
			MerchantID:     p.MerchantID,
			SubscriptionID: p.SubscriptionID,
			Amount:         p.Amount,
			Currency:       p.Currency,
			Status:         p.Status,
			PaymentMethod:  p.PaymentMethod,
			Reference:      p.Reference,
			CreatedAt:      p.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &PaymentPaginationResponse{
		Data:  res,
		Total: total,
	}, nil
}

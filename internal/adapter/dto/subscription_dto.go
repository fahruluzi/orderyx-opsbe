package dto

import "time"

type SubscriptionResponse struct {
	ID            int64     `json:"id"`
	MerchantID    int64     `json:"merchant_id"`
	PlanName      string    `json:"plan_name"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	CreatedAt     time.Time `json:"created_at"`
}

type ExtendTrialRequest struct {
	EndDate time.Time `json:"end_date"`
}

type ChangePlanRequest struct {
	PlanID int64 `json:"plan_id"`
}

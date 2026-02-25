package dto

import "time"

type MerchantListResponse struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	BusinessType         string     `json:"business_type"`
	IsActive             bool       `json:"is_active"`
	SubscriptionPlan     string     `json:"subscription_plan"`
	SubscriptionStatus   string     `json:"subscription_status"`
	SubscriptionTrialEnd *time.Time `json:"subscription_trial_end"`
	CreatedAt            time.Time  `json:"created_at"`
}

type MerchantDetailResponse struct {
	ID                   int64                  `json:"id"`
	Name                 string                 `json:"name"`
	BusinessType         string                 `json:"business_type"`
	Email                *string                `json:"email"`
	Phone                *string                `json:"phone"`
	Address              *string                `json:"address"`
	TaxID                *string                `json:"tax_id"`
	IsActive             bool                   `json:"is_active"`
	IsOnboardingComplete bool                   `json:"is_onboarding_complete"`
	Settings             map[string]interface{} `json:"settings"`
	CreatedAt            time.Time              `json:"created_at"`

	// Subscription summary
	SubscriptionPlan   string     `json:"subscription_plan"`
	SubscriptionStatus string     `json:"subscription_status"`
	SubscriptionEnd    *time.Time `json:"subscription_end"`

	// Stats
	TotalUsers     int64 `json:"total_users"`
	TotalOrders    int64 `json:"total_orders"`
	TotalPurchases int64 `json:"total_purchases"`
}

type MerchantPaginationRequest struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
	Status string `query:"status"` // Active, Inactive, Trial, Expired
	Plan   string `query:"plan"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"total_count"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
}

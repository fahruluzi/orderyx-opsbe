package domain

import "time"

// BusinessType represents the type of business
type BusinessType string

const (
	BusinessTypeManufacturing BusinessType = "MANUFACTURING"
	BusinessTypeRetail        BusinessType = "RETAIL"
	BusinessTypeDistribution  BusinessType = "DISTRIBUTION"
	BusinessTypeService       BusinessType = "SERVICE"
	BusinessTypeFoodBeverage  BusinessType = "FOOD_BEVERAGE"
)

// Merchant represents a business entity read from the main Orderyx DB
type Merchant struct {
	ID                   int64                  `json:"id" gorm:"primaryKey"`
	Name                 string                 `json:"name"`
	BusinessType         BusinessType           `json:"business_type"`
	Address              *string                `json:"address"`
	Phone                *string                `json:"phone"`
	Email                *string                `json:"email"`
	TaxID                *string                `json:"tax_id"`
	IsActive             bool                   `json:"is_active"`
	IsOnboardingComplete bool                   `json:"is_onboarding_complete"`
	Settings             map[string]interface{} `json:"settings,omitempty" gorm:"type:jsonb;serializer:json"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	DeletedAt            *time.Time             `json:"deleted_at,omitempty"`
}

type SubscriptionPlan struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Subscription represents the merchant's subscription state
type Subscription struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	MerchantID    int64     `json:"merchant_id"`
	PlanID        int64     `json:"plan_id"`
	Status        string    `json:"status"` // expected: "TRIAL", "ACTIVE", "EXPIRED", "CANCELLED", "PENDING"
	PaymentStatus string    `json:"payment_status"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Merchant *Merchant         `json:"merchant,omitempty" gorm:"foreignKey:MerchantID"`
	Plan     *SubscriptionPlan `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
}

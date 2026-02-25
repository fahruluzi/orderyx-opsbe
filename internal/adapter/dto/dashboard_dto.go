package dto

type DashboardSummaryResponse struct {
	TotalMerchants   int64   `json:"total_merchants"`
	ActiveMerchants  int64   `json:"active_merchants"`
	TrialMerchants   int64   `json:"trial_merchants"`
	ExpiredMerchants int64   `json:"expired_merchants"`
	RevenueMTD       float64 `json:"revenue_mtd"`
}

type GrowthData struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type RevenueData struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

type DashboardGrowthResponse struct {
	Growth []GrowthData `json:"growth"`
}

type DashboardRevenueResponse struct {
	Revenue []RevenueData `json:"revenue"`
}

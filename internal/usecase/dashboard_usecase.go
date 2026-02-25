package usecase

import (
	"context"

	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/dto"
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/repository"
)

type DashboardUsecase interface {
	GetSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error)
	GetGrowth(ctx context.Context, days int) (*dto.DashboardGrowthResponse, error)
	GetRevenue(ctx context.Context, days int) (*dto.DashboardRevenueResponse, error)
}

type dashboardUsecase struct {
	repo repository.DashboardRepository
}

func NewDashboardUsecase(repo repository.DashboardRepository) DashboardUsecase {
	return &dashboardUsecase{repo: repo}
}

func (u *dashboardUsecase) GetSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error) {
	return u.repo.GetSummary(ctx)
}

func (u *dashboardUsecase) GetGrowth(ctx context.Context, days int) (*dto.DashboardGrowthResponse, error) {
	growthData, err := u.repo.GetGrowth(ctx, days)
	if err != nil {
		return nil, err
	}
	return &dto.DashboardGrowthResponse{Growth: growthData}, nil
}

func (u *dashboardUsecase) GetRevenue(ctx context.Context, days int) (*dto.DashboardRevenueResponse, error) {
	revenueData, err := u.repo.GetRevenue(ctx, days)
	if err != nil {
		return nil, err
	}
	return &dto.DashboardRevenueResponse{Revenue: revenueData}, nil
}

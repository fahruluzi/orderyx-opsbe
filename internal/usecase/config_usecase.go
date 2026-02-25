package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/repository"
	"github.com/fahruluzi/orderyx-opsbe/internal/domain"
	"github.com/fahruluzi/orderyx-opsbe/pkg/jwt"
)

type ConfigUsecase interface {
	GetAllConfigs(ctx context.Context) ([]domain.SystemConfig, error)
	UpdateConfig(ctx context.Context, actor *jwt.OpsClaims, key string, value interface{}) error
}

type configUsecase struct {
	repo      repository.ConfigRepository
	auditRepo repository.AuditLogRepository
}

func NewConfigUsecase(repo repository.ConfigRepository, auditRepo repository.AuditLogRepository) ConfigUsecase {
	return &configUsecase{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (u *configUsecase) GetAllConfigs(ctx context.Context) ([]domain.SystemConfig, error) {
	return u.repo.GetAllConfigs(ctx)
}

func (u *configUsecase) UpdateConfig(ctx context.Context, actor *jwt.OpsClaims, key string, value interface{}) error {
	// Let's verify if config exists
	oldConf, err := u.repo.GetConfigByKey(ctx, key)
	if err != nil {
		return errors.New("configuration key not found")
	}

	err = u.repo.UpdateConfig(ctx, key, value)
	if err != nil {
		return err
	}

	// Audit Log
	detailMsg, _ := json.Marshal(map[string]interface{}{
		"key":       key,
		"old_value": oldConf.Value,
		"new_value": value,
	})
	_ = u.auditRepo.LogAction(ctx, actor.UserID, "UPDATE_SYSTEM_CONFIG", "config", nil, string(detailMsg))

	return nil
}

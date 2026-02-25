package repository

import (
	"context"

	"github.com/fahruluzi/orderyx-opsbe/internal/domain"
	"gorm.io/gorm"
)

type ConfigRepository interface {
	GetAllConfigs(ctx context.Context) ([]domain.SystemConfig, error)
	GetConfigByKey(ctx context.Context, key string) (*domain.SystemConfig, error)
	UpdateConfig(ctx context.Context, key string, value interface{}) error
}

type configRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepository{db: db}
}

func (r *configRepository) GetAllConfigs(ctx context.Context) ([]domain.SystemConfig, error) {
	var configs []domain.SystemConfig
	err := r.db.WithContext(ctx).Find(&configs).Error
	return configs, err
}

func (r *configRepository) GetConfigByKey(ctx context.Context, key string) (*domain.SystemConfig, error) {
	var config domain.SystemConfig
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *configRepository) UpdateConfig(ctx context.Context, key string, value interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.SystemConfig{}).Where("key = ?", key).Update("value", value).Error
}

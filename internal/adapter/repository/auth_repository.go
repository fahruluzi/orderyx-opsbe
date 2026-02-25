package repository

import (
	"context"

	"github.com/fahruluzi/orderyx-opsbe/internal/domain"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.OpsUser, error)
	FindByID(ctx context.Context, id int64) (*domain.OpsUser, error)
	UpdateLastLogin(ctx context.Context, id int64) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*domain.OpsUser, error) {
	var user domain.OpsUser
	err := r.db.WithContext(ctx).Where("email = ? AND is_active = ?", email, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindByID(ctx context.Context, id int64) (*domain.OpsUser, error) {
	var user domain.OpsUser
	err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&domain.OpsUser{}).
		Where("id = ?", id).
		Update("last_login_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

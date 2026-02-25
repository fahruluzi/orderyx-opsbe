package usecase

import (
	"context"
	"errors"

	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/dto"
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/repository"
	"github.com/fahruluzi/orderyx-opsbe/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	Me(ctx context.Context, userID int64) (dto.OpsUserMeta, error)
}

type authUsecase struct {
	authRepo   repository.AuthRepository
	jwtService *jwt.JWTService
}

func NewAuthUsecase(authRepo repository.AuthRepository, jwtService *jwt.JWTService) AuthUsecase {
	return &authUsecase{
		authRepo:   authRepo,
		jwtService: jwtService,
	}
}

func (u *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := u.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return dto.LoginResponse{}, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return dto.LoginResponse{}, errors.New("invalid email or password")
	}

	token, err := u.jwtService.GenerateToken(user.ID, user.Email, string(user.Role), user.FullName)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	_ = u.authRepo.UpdateLastLogin(ctx, user.ID)

	return dto.LoginResponse{
		Token: token,
		User: dto.OpsUserMeta{
			ID:          user.ID,
			FullName:    user.FullName,
			Email:       user.Email,
			Role:        string(user.Role),
			LastLoginAt: user.LastLoginAt,
		},
	}, nil
}

func (u *authUsecase) Me(ctx context.Context, userID int64) (dto.OpsUserMeta, error) {
	user, err := u.authRepo.FindByID(ctx, userID)
	if err != nil {
		return dto.OpsUserMeta{}, errors.New("user not found")
	}

	return dto.OpsUserMeta{
		ID:          user.ID,
		FullName:    user.FullName,
		Email:       user.Email,
		Role:        string(user.Role),
		LastLoginAt: user.LastLoginAt,
	}, nil
}

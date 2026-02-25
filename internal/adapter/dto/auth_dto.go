package dto

import "time"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  OpsUserMeta `json:"user"`
}

type OpsUserMeta struct {
	ID          int64      `json:"id"`
	FullName    string     `json:"full_name"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

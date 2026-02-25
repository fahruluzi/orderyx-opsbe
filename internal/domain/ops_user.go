package domain

import (
	"time"

	"gorm.io/gorm"
)

type OpsRole string

const (
	OpsRoleSuperAdmin OpsRole = "super_admin"
	OpsRoleSupport    OpsRole = "support"
	OpsRoleViewer     OpsRole = "viewer"
)

type OpsUser struct {
	ID           int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	FullName     string         `json:"full_name" gorm:"type:varchar(255);not null"`
	Email        string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string         `json:"-" gorm:"type:varchar(255);not null"` // JSON exclude
	Role         OpsRole        `json:"role" gorm:"type:varchar(50);not null;default:'viewer'"`
	IsActive     bool           `json:"is_active" gorm:"not null;default:true"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type OpsAuditLog struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	OpsUserID  int64     `json:"ops_user_id" gorm:"not null;index"`
	OpsUser    OpsUser   `json:"-" gorm:"foreignKey:OpsUserID"`
	Action     string    `json:"action" gorm:"type:varchar(100);not null"`
	TargetType string    `json:"target_type" gorm:"type:varchar(50);not null"`
	TargetID   *int64    `json:"target_id" gorm:"index"`
	Details    string    `json:"details" gorm:"type:jsonb"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

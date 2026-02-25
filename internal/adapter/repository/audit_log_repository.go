package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	ActorID    int64     `json:"actor_id"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   *int64    `json:"entity_id"`
	Details    string    `json:"details" gorm:"type:jsonb"`
	CreatedAt  time.Time `json:"created_at"`
}

type AuditLogRepository interface {
	LogAction(ctx context.Context, actorID int64, action, entityType string, entityID *int64, details string) error
	GetLogs(ctx context.Context, limit, offset int) ([]AuditLog, int64, error)
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) LogAction(ctx context.Context, actorID int64, action, entityType string, entityID *int64, details string) error {
	log := AuditLog{
		ActorID:    actorID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    details,
	}
	return r.db.WithContext(ctx).Table("ops_audit_logs").Create(&log).Error
}

func (r *auditLogRepository) GetLogs(ctx context.Context, limit, offset int) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	query := r.db.WithContext(ctx).Table("ops_audit_logs")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at desc").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

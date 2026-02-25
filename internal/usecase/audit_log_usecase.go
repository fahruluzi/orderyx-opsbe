package usecase

import (
	"context"

	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/repository"
)

type AuditLogResponse struct {
	ID         int64  `json:"id"`
	ActorID    int64  `json:"actor_id"`
	Action     string `json:"action"`
	EntityType string `json:"entity_type"`
	EntityID   *int64 `json:"entity_id"`
	Details    string `json:"details"`
	CreatedAt  string `json:"created_at"`
}

type AuditLogPaginationResponse struct {
	Data  []AuditLogResponse `json:"data"`
	Total int64              `json:"total"`
}

type AuditLogUsecase interface {
	GetLogs(ctx context.Context, limit, offset int) (*AuditLogPaginationResponse, error)
}

type auditLogUsecase struct {
	repo repository.AuditLogRepository
}

func NewAuditLogUsecase(repo repository.AuditLogRepository) AuditLogUsecase {
	return &auditLogUsecase{repo: repo}
}

func (u *auditLogUsecase) GetLogs(ctx context.Context, limit, offset int) (*AuditLogPaginationResponse, error) {
	logs, total, err := u.repo.GetLogs(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var res []AuditLogResponse
	for _, l := range logs {
		res = append(res, AuditLogResponse{
			ID:         l.ID,
			ActorID:    l.ActorID,
			Action:     l.Action,
			EntityType: l.EntityType,
			EntityID:   l.EntityID,
			Details:    l.Details,
			CreatedAt:  l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &AuditLogPaginationResponse{
		Data:  res,
		Total: total,
	}, nil
}

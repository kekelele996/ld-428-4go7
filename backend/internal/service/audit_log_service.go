package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/model"
	"github.com/artvault/artvault/internal/repository"
	"github.com/artvault/artvault/internal/util"
)

// AuditLogService 操作日志。
type AuditLogService struct {
	repo   *repository.AuditLogRepository
	logger *slog.Logger
}

func NewAuditLogService(repo *repository.AuditLogRepository, logger *slog.Logger) *AuditLogService {
	return &AuditLogService{repo: repo, logger: logger}
}

func (s *AuditLogService) Write(ctx context.Context, userID, action, entity, entityID, detail string) {
	entry := &model.AuditLog{ID: util.NewID("al"), UserID: userID, Action: action, Entity: entity, EntityID: entityID, Detail: detail, CreatedAt: time.Now()}
	if err := s.repo.Create(ctx, entry); err != nil {
		s.logger.Error("write audit log failed", "error", err)
		return
	}
	s.logger.Info(constants.LogAuditWritten, "action", action)
}

func (s *AuditLogService) List(ctx context.Context, limit int64) ([]model.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	list, err := s.repo.List(ctx, limit)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return list, nil
}

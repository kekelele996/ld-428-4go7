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

// ReviewService 内容审核。
type ReviewService struct {
	repo   *repository.ReviewRepository
	logger *slog.Logger
}

func NewReviewService(repo *repository.ReviewRepository, logger *slog.Logger) *ReviewService {
	return &ReviewService{repo: repo, logger: logger}
}

func (s *ReviewService) Create(ctx context.Context, targetType, targetID, reviewerID, result, comment string) (*model.ReviewLog, error) {
	entry := &model.ReviewLog{ID: util.NewID("rv"), TargetType: targetType, TargetID: targetID, ReviewerID: reviewerID, Result: result, Comment: comment, ReviewedAt: time.Now()}
	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogReviewSubmitted, "target", targetID, "result", result)
	return entry, nil
}

func (s *ReviewService) List(ctx context.Context, targetType string) ([]model.ReviewLog, error) {
	list, err := s.repo.List(ctx, targetType)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return list, nil
}

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

// InteractionService 互动业务逻辑。
type InteractionService struct {
	repo        *repository.InteractionRepository
	artworkRepo *repository.ArtworkRepository
	logger      *slog.Logger
}

func NewInteractionService(repo *repository.InteractionRepository, artworkRepo *repository.ArtworkRepository, logger *slog.Logger) *InteractionService {
	return &InteractionService{repo: repo, artworkRepo: artworkRepo, logger: logger}
}

func (s *InteractionService) Create(ctx context.Context, userID, targetType, targetID, itype, comment string) (*model.Interaction, error) {
	i := &model.Interaction{ID: util.NewID("i"), UserID: userID, TargetType: targetType, TargetID: targetID, Type: itype, Comment: comment, CreatedAt: time.Now()}
	if err := s.repo.Create(ctx, i); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	// 更新作品计数
	if targetType == "Artwork" {
		var views, likes, bookmarks int64
		switch itype {
		case constants.InteractionLike:
			likes = 1
		case constants.InteractionBookmark:
			bookmarks = 1
		}
		_ = s.artworkRepo.IncCounters(ctx, targetID, views, likes, bookmarks)
	}
	s.logger.Info(constants.LogInteraction, "type", itype, "target", targetID)
	return i, nil
}

func (s *InteractionService) List(ctx context.Context, targetType, targetID string) ([]model.Interaction, error) {
	list, err := s.repo.List(ctx, targetType, targetID)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return list, nil
}

// Cancel 取消点赞/收藏。
func (s *InteractionService) Cancel(ctx context.Context, userID, targetType, targetID, itype string) error {
	if err := s.repo.DeleteByUserTargetType(ctx, userID, targetType, targetID, itype); err != nil {
		return util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return nil
}

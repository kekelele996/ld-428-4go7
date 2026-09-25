package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/model"
	"github.com/artvault/artvault/internal/repository"
	"github.com/artvault/artvault/internal/util"
)

// ArtworkService 作品业务逻辑。
type ArtworkService struct {
	repo   *repository.ArtworkRepository
	logger *slog.Logger
}

func NewArtworkService(repo *repository.ArtworkRepository, logger *slog.Logger) *ArtworkService {
	return &ArtworkService{repo: repo, logger: logger}
}

func (s *ArtworkService) Create(ctx context.Context, a *model.Artwork) (*model.Artwork, error) {
	a.ID = util.NewID("art")
	if a.Status == "" {
		a.Status = constants.ArtworkDraft
	}
	a.ReviewStatus = "Pending"
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogArtworkCreated, "artworkID", a.ID)
	return a, nil
}

func (s *ArtworkService) Get(ctx context.Context, id string) (*model.Artwork, error) {
	a, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(constants.CodeNotFound, "作品不存在", nil)
	}
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return a, nil
}

func (s *ArtworkService) List(ctx context.Context, onlyPublished bool, artistID string) ([]model.Artwork, error) {
	list, err := s.repo.List(ctx, onlyPublished, artistID)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return list, nil
}

func (s *ArtworkService) ChangeStatus(ctx context.Context, id, status string) (*model.Artwork, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	a.Status = status
	a.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if status == constants.ArtworkPublished {
		s.logger.Info(constants.LogArtworkPublished, "artworkID", id)
	}
	return a, nil
}

func (s *ArtworkService) Review(ctx context.Context, id string, result, comment, reviewerID string) (*model.Artwork, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	a.ReviewStatus = result
	a.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogReviewDecided, "artworkID", id, "result", result)
	_ = reviewerID
	return a, nil
}

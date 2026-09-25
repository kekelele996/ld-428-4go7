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

// ExhibitionService 展览业务逻辑。
type ExhibitionService struct {
	repo        *repository.ExhibitionRepository
	artworkRepo *repository.ArtworkRepository
	logger      *slog.Logger
}

func NewExhibitionService(repo *repository.ExhibitionRepository, artworkRepo *repository.ArtworkRepository, logger *slog.Logger) *ExhibitionService {
	return &ExhibitionService{repo: repo, artworkRepo: artworkRepo, logger: logger}
}

func (s *ExhibitionService) Create(ctx context.Context, e *model.Exhibition) (*model.Exhibition, error) {
	e.ID = util.NewID("exh")
	if e.Status == "" {
		e.Status = constants.ExhibitionPlanning
	}
	e.ReviewStatus = "Pending"
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogExhibitionCreate, "exhibitionID", e.ID)
	return e, nil
}

func (s *ExhibitionService) Get(ctx context.Context, id string) (*model.Exhibition, error) {
	e, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(constants.CodeNotFound, "展览不存在", nil)
	}
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return e, nil
}

func (s *ExhibitionService) List(ctx context.Context, onlyPublished bool) ([]model.Exhibition, error) {
	list, err := s.repo.List(ctx, onlyPublished)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return list, nil
}

func (s *ExhibitionService) ChangeStatus(ctx context.Context, id, status string) (*model.Exhibition, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	e.Status = status
	e.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, e); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if status == constants.ExhibitionActive {
		s.logger.Info(constants.LogExhibitionPublish, "exhibitionID", id)
	}
	return e, nil
}

// AddArtwork 展览收录作品。
func (s *ExhibitionService) AddArtwork(ctx context.Context, id, artworkID string) (*model.Exhibition, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	for _, aid := range e.ArtworkIDs {
		if aid == artworkID {
			return e, nil
		}
	}
	e.ArtworkIDs = append(e.ArtworkIDs, artworkID)
	e.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, e); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return e, nil
}

package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/dto"
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

// Readiness 实时统计展览关联作品的公开就绪情况：可公开展出 = 已发布 + 审核通过 + 有图。
// 所有异常（未发布/待审核/审核未通过/已售/下架/无图/关联缺失）都会计入阻断原因。
func (s *ExhibitionService) Readiness(ctx context.Context, id string) (*dto.ExhibitionReadinessResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(constants.CodeNotFound, "展览不存在", nil)
	}
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	artworks, err := s.artworkRepo.ListByIDs(ctx, e.ArtworkIDs)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}

	res := evaluateReadiness(e.ID, e.ArtworkIDs, artworks)
	s.logger.Info(constants.LogExhibitionReadiness, "exhibitionID", id, "total", res.TotalArtworks, "ready", res.Ready, "readyRatio", res.ReadyRatio)
	return res, nil
}

// evaluateReadiness 纯函数：按关联作品实时统计就绪指标，所有异常都计入阻断原因。
func evaluateReadiness(exhibitionID string, artworkIDs []string, artworks []model.Artwork) *dto.ExhibitionReadinessResponse {
	res := &dto.ExhibitionReadinessResponse{
		ExhibitionID:    exhibitionID,
		TotalArtworks:   len(artworkIDs),
		BlockingReasons: []string{},
	}
	for _, a := range artworks {
		hasImage := len(a.ImageURLs) > 0
		switch a.ReviewStatus {
		case constants.ReviewApproved:
		case "Pending":
			res.PendingReviewCount++
		default:
			res.RejectedCount++
		}
		switch a.Status {
		case constants.ArtworkDraft:
			res.DraftCount++
		case constants.ArtworkSold:
			res.SoldCount++
		case constants.ArtworkArchived:
			res.ArchivedCount++
		}
		if !hasImage {
			res.NoImageCount++
		}
		if a.Status == constants.ArtworkPublished && a.ReviewStatus == constants.ReviewApproved && hasImage {
			res.PublicReadyCount++
		}
	}
	res.MissingCount = res.TotalArtworks - len(artworks)
	if res.TotalArtworks > 0 {
		res.ReadyRatio = float64(res.PublicReadyCount) / float64(res.TotalArtworks)
	}
	res.Ready = res.TotalArtworks > 0 && res.PublicReadyCount == res.TotalArtworks

	if res.TotalArtworks == 0 {
		res.BlockingReasons = append(res.BlockingReasons, constants.MsgReadinessEmptyExhibition)
	}
	if res.MissingCount > 0 {
		res.BlockingReasons = append(res.BlockingReasons, fmt.Sprintf(constants.MsgReadinessMissingArtworks, res.MissingCount))
	}
	if res.DraftCount > 0 {
		res.BlockingReasons = append(res.BlockingReasons, fmt.Sprintf(constants.MsgReadinessDraftArtworks, res.DraftCount))
	}
	if res.PendingReviewCount > 0 {
		res.BlockingReasons = append(res.BlockingReasons, fmt.Sprintf(constants.MsgReadinessPendingReview, res.PendingReviewCount))
	}
	if res.RejectedCount > 0 {
		res.BlockingReasons = append(res.BlockingReasons, fmt.Sprintf(constants.MsgReadinessRejectedReview, res.RejectedCount))
	}
	if res.SoldCount > 0 {
		res.BlockingReasons = append(res.BlockingReasons, fmt.Sprintf(constants.MsgReadinessSoldArtworks, res.SoldCount))
	}
	if res.ArchivedCount > 0 {
		res.BlockingReasons = append(res.BlockingReasons, fmt.Sprintf(constants.MsgReadinessArchivedArtworks, res.ArchivedCount))
	}
	if res.NoImageCount > 0 {
		res.BlockingReasons = append(res.BlockingReasons, fmt.Sprintf(constants.MsgReadinessNoImage, res.NoImageCount))
	}
	return res
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

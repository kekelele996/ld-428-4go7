package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/model"
	"github.com/artvault/artvault/internal/repository"
	"github.com/artvault/artvault/internal/util"
)

// readinessAgg 阻断原因累加器。
type readinessAgg struct {
	count      int
	artworkIDs []string
	message    string
}

func (a *readinessAgg) add(artworkID, message string) {
	if a.count == 0 {
		// 首条命中的详细文案作为该原因的说明。
		a.message = message
	}
	a.count++
	a.artworkIDs = append(a.artworkIDs, artworkID)
}

// GetReadiness 按关联作品实时统计展览的公开展出就绪情况。
// 任何异常（缺失引用、草稿、待审核/驳回/标记、已售、下架、无图、未知状态）都作为阻断原因。
func (s *ExhibitionService) GetReadiness(ctx context.Context, id string) (*model.ExhibitionReadiness, error) {
	e, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(constants.CodeNotFound, fmt.Sprintf("Exhibition[id=%s] readiness failed: exhibition not found", id), nil)
	}
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, fmt.Errorf("Exhibition[id=%s] readiness load exhibition: %w", id, err))
	}

	if len(e.ArtworkIDs) == 0 {
		r := evaluateReadiness(e, nil)
		s.logger.Info(constants.LogExhibitionReadiness, "exhibitionID", id, "ready", false, "total", 0, "empty", true)
		return r, nil
	}

	artworks, err := s.artworkRepo.FindByIDs(ctx, e.ArtworkIDs)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, fmt.Errorf("Exhibition[id=%s] readiness load artworks: %w", id, err))
	}
	r := evaluateReadiness(e, artworks)

	s.logger.Info(constants.LogExhibitionReadiness,
		"exhibitionID", id, "ready", r.Ready, "total", r.Total,
		"readyCount", r.ReadyCount, "pending", r.PendingCount,
		"sold", r.SoldCount, "archived", r.ArchivedCount, "noImage", r.NoImageCount,
		"reasonCount", len(r.Reasons))
	return r, nil
}

// evaluateReadiness 按关联作品实时统计就绪情况（纯函数，便于测试）。
// 一个作品可同时命中多个原因（如已售且无图）。
func evaluateReadiness(e *model.Exhibition, artworks []model.Artwork) *model.ExhibitionReadiness {
	r := &model.ExhibitionReadiness{ExhibitionID: e.ID, Exhibition: e, Reasons: []model.ReadinessBlockReason{}}
	aggs := map[model.ReadinessReasonCode]*readinessAgg{}
	addReason := func(code model.ReadinessReasonCode, artworkID, message string) {
		a, ok := aggs[code]
		if !ok {
			a = &readinessAgg{}
			aggs[code] = a
		}
		a.add(artworkID, message)
	}

	r.Total = len(e.ArtworkIDs)
	if r.Total == 0 {
		r.IsEmpty = true
		r.Ready = false
		r.ReadyRatio = 0
		return r
	}

	found := make(map[string]*model.Artwork, len(artworks))
	for i := range artworks {
		found[artworks[i].ID] = &artworks[i]
	}

	// 保留展览收录顺序逐个评估；一个作品可同时命中多个原因。
	readyCount := 0
	for _, artworkID := range e.ArtworkIDs {
		a := found[artworkID]
		blocked := false
		if a == nil {
			addReason(constants.ReadinessReasonMissingArtwork, artworkID, fmt.Sprintf(constants.MsgReadinessMissingFmt, artworkID))
			blocked = true
		} else {
			if a.Status == constants.ArtworkDraft {
				addReason(constants.ReadinessReasonDraft, artworkID, fmt.Sprintf(constants.MsgReadinessDraftFmt, artworkID))
				blocked = true
			} else if a.Status == constants.ArtworkSold {
				addReason(constants.ReadinessReasonSold, artworkID, fmt.Sprintf(constants.MsgReadinessSoldFmt, artworkID))
				blocked = true
				r.SoldCount++
			} else if a.Status == constants.ArtworkArchived {
				addReason(constants.ReadinessReasonArchived, artworkID, fmt.Sprintf(constants.MsgReadinessArchFmt, artworkID))
				blocked = true
				r.ArchivedCount++
			} else if a.Status != constants.ArtworkPublished {
				addReason(constants.ReadinessReasonStatusUnknown, artworkID, fmt.Sprintf(constants.MsgReadinessStsUnkFmt, artworkID, a.Status))
				blocked = true
			}

			switch a.ReviewStatus {
			case constants.ReviewApproved:
				// 审核通过，不阻断
			case constants.ReviewPending:
				addReason(constants.ReadinessReasonPendingReview, artworkID, fmt.Sprintf(constants.MsgReadinessPendingFmt, artworkID))
				blocked = true
				r.PendingCount++
			case constants.ReviewRejected:
				addReason(constants.ReadinessReasonRejected, artworkID, fmt.Sprintf(constants.MsgReadinessRejectFmt, artworkID))
				blocked = true
			case constants.ReviewFlagged:
				addReason(constants.ReadinessReasonFlagged, artworkID, fmt.Sprintf(constants.MsgReadinessFlagFmt, artworkID))
				blocked = true
			default:
				addReason(constants.ReadinessReasonReviewUnknown, artworkID, fmt.Sprintf(constants.MsgReadinessRevUnkFmt, artworkID, a.ReviewStatus))
				blocked = true
			}

			if len(a.ImageURLs) == 0 || allImagesBlank(a.ImageURLs) {
				addReason(constants.ReadinessReasonNoImage, artworkID, fmt.Sprintf(constants.MsgReadinessNoImgFmt, artworkID))
				blocked = true
				r.NoImageCount++
			}
		}
		if !blocked {
			readyCount++
		}
	}

	r.ReadyCount = readyCount
	r.ReadyRatio = float64(readyCount) / float64(r.Total)
	r.Ready = readyCount == r.Total
	r.Reasons = buildReadinessReasons(aggs)
	return r
}

// buildReadinessReasons 按固定顺序输出阻断原因，避免 map 遍历导致响应顺序不稳定。
func buildReadinessReasons(aggs map[model.ReadinessReasonCode]*readinessAgg) []model.ReadinessBlockReason {
	order := []model.ReadinessReasonCode{
		constants.ReadinessReasonMissingArtwork,
		constants.ReadinessReasonDraft,
		constants.ReadinessReasonPendingReview,
		constants.ReadinessReasonRejected,
		constants.ReadinessReasonFlagged,
		constants.ReadinessReasonReviewUnknown,
		constants.ReadinessReasonSold,
		constants.ReadinessReasonArchived,
		constants.ReadinessReasonStatusUnknown,
		constants.ReadinessReasonNoImage,
	}
	reasons := []model.ReadinessBlockReason{}
	for _, code := range order {
		a, ok := aggs[code]
		if !ok {
			continue
		}
		reasons = append(reasons, model.ReadinessBlockReason{
			Code:       code,
			Message:    a.message,
			Count:      a.count,
			ArtworkIDs: a.artworkIDs,
		})
	}
	return reasons
}

func allImagesBlank(urls []string) bool {
	for _, u := range urls {
		if strings.TrimSpace(u) != "" {
			return false
		}
	}
	return true
}

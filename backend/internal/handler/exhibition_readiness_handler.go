package handler

import (
	"context"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/dto"
	"github.com/artvault/artvault/internal/model"
	"github.com/artvault/artvault/internal/service"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// ExhibitionReadinessHandler 展览就绪接口（单独文件，避免与 CRUD 职责混杂）。
type ExhibitionReadinessHandler struct {
	svc *service.ExhibitionService
}

func NewExhibitionReadinessHandler(svc *service.ExhibitionService) *ExhibitionReadinessHandler {
	return &ExhibitionReadinessHandler{svc: svc}
}

// Get GET /api/v1/exhibitions/:id/readiness
func (h *ExhibitionReadinessHandler) Get(c *gin.Context) {
	r, err := h.svc.GetReadiness(context.Background(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, toExhibitionReadinessResponse(r))
}

func toExhibitionReadinessResponse(r *model.ExhibitionReadiness) dto.ExhibitionReadinessResponse {
	conclusion := constants.MsgReadinessReady
	if r.IsEmpty {
		conclusion = constants.MsgReadinessEmpty
	} else if !r.Ready {
		conclusion = constants.MsgReadinessBlocked
	}

	resp := dto.ExhibitionReadinessResponse{
		ExhibitionID:  r.ExhibitionID,
		Conclusion:    conclusion,
		Ready:         r.Ready,
		IsEmpty:       r.IsEmpty,
		Total:         r.Total,
		ReadyCount:    r.ReadyCount,
		PendingCount:  r.PendingCount,
		SoldCount:     r.SoldCount,
		ArchivedCount: r.ArchivedCount,
		NoImageCount:  r.NoImageCount,
		ReadyRatio:    r.ReadyRatio,
		Reasons:       []dto.ExhibitionReadinessReasonItem{},
	}
	if r.Exhibition != nil {
		resp.Exhibition = &dto.ExhibitionBriefResponse{
			ID:           r.Exhibition.ID,
			Title:        r.Exhibition.Title,
			Status:       r.Exhibition.Status,
			ReviewStatus: r.Exhibition.ReviewStatus,
		}
	}
	for _, reason := range r.Reasons {
		resp.Reasons = append(resp.Reasons, dto.ExhibitionReadinessReasonItem{
			Code:       string(reason.Code),
			Message:    reason.Message,
			Count:      reason.Count,
			ArtworkIDs: reason.ArtworkIDs,
		})
	}
	return resp
}

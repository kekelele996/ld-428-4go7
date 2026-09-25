package handler

import (
	"context"
	"net/http"

	"github.com/artvault/artvault/internal/dto"
	"github.com/artvault/artvault/internal/service"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// ReviewHandler 内容审核接口。
type ReviewHandler struct {
	svc *service.ReviewService
}

func NewReviewHandler(svc *service.ReviewService) *ReviewHandler { return &ReviewHandler{svc: svc} }

func (h *ReviewHandler) List(c *gin.Context) {
	targetType := c.DefaultQuery("targetType", "")
	list, err := h.svc.List(context.Background(), targetType)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, list)
}

func (h *ReviewHandler) Create(c *gin.Context) {
	var req dto.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	targetType := c.DefaultQuery("targetType", "Artwork")
	targetID := c.Param("targetId")
	userID, _ := c.Get("userID")
	entry, err := h.svc.Create(context.Background(), targetType, targetID, userID.(string), req.Result, req.Comment)
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, entry)
}

package handler

import (
	"context"
	"net/http"

	"github.com/artvault/artvault/internal/dto"
	"github.com/artvault/artvault/internal/service"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// InteractionHandler 互动接口。
type InteractionHandler struct {
	svc *service.InteractionService
}

func NewInteractionHandler(svc *service.InteractionService) *InteractionHandler {
	return &InteractionHandler{svc: svc}
}

func (h *InteractionHandler) List(c *gin.Context) {
	targetType := c.DefaultQuery("targetType", "")
	targetID := c.DefaultQuery("targetId", "")
	list, err := h.svc.List(context.Background(), targetType, targetID)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, list)
}

func (h *InteractionHandler) Create(c *gin.Context) {
	var req dto.InteractionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	userID, _ := c.Get("userID")
	i, err := h.svc.Create(context.Background(), userID.(string), req.TargetType, req.TargetID, req.Type, req.Comment)
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, i)
}

func (h *InteractionHandler) Cancel(c *gin.Context) {
	userID, _ := c.Get("userID")
	targetType := c.DefaultQuery("targetType", "")
	targetID := c.DefaultQuery("targetId", "")
	itype := c.DefaultQuery("type", "")
	if err := h.svc.Cancel(context.Background(), userID.(string), targetType, targetID, itype); err != nil {
		c.Error(err)
		return
	}
	util.OK(c, gin.H{"cancelled": true})
}

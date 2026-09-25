package handler

import (
	"context"
	"strconv"

	"github.com/artvault/artvault/internal/service"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// AuditHandler 操作日志接口。
type AuditHandler struct {
	svc *service.AuditLogService
}

func NewAuditHandler(svc *service.AuditLogService) *AuditHandler { return &AuditHandler{svc: svc} }

func (h *AuditHandler) List(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	list, err := h.svc.List(context.Background(), limit)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, list)
}

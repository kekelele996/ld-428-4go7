package handler

import (
	"context"
	"net/http"

	"github.com/artvault/artvault/internal/dto"
	"github.com/artvault/artvault/internal/service"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// ArtistHandler 艺术家接口。
type ArtistHandler struct {
	svc *service.ArtistService
}

func NewArtistHandler(svc *service.ArtistService) *ArtistHandler { return &ArtistHandler{svc: svc} }

func (h *ArtistHandler) List(c *gin.Context) {
	list, err := h.svc.List(context.Background())
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, list)
}

func (h *ArtistHandler) Get(c *gin.Context) {
	a, err := h.svc.Get(context.Background(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}

func (h *ArtistHandler) Upsert(c *gin.Context) {
	var req dto.ArtistUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	userID, _ := c.Get("userID")
	a, err := h.svc.UpsertForUser(context.Background(), userID.(string), req.ArtistName, req.Bio, req.AvatarURL, req.Mediums)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}

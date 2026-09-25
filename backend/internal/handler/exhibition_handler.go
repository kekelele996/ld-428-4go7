package handler

import (
	"context"
	"net/http"

	"github.com/artvault/artvault/internal/dto"
	"github.com/artvault/artvault/internal/model"
	"github.com/artvault/artvault/internal/service"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// ExhibitionHandler 展览接口。
type ExhibitionHandler struct {
	svc *service.ExhibitionService
}

func NewExhibitionHandler(svc *service.ExhibitionService) *ExhibitionHandler {
	return &ExhibitionHandler{svc: svc}
}

func (h *ExhibitionHandler) List(c *gin.Context) {
	onlyPublished := c.DefaultQuery("published", "true") == "true"
	list, err := h.svc.List(context.Background(), onlyPublished)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, list)
}

func (h *ExhibitionHandler) Get(c *gin.Context) {
	e, err := h.svc.Get(context.Background(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, e)
}

func (h *ExhibitionHandler) Create(c *gin.Context) {
	var req dto.ExhibitionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	userID, _ := c.Get("userID")
	e := &model.Exhibition{Title: req.Title, Description: req.Description, CuratorID: userID.(string), StartDate: req.StartDate, EndDate: req.EndDate, Type: req.Type, CoverURL: req.CoverURL, ArtworkIDs: req.ArtworkIDs}
	created, err := h.svc.Create(context.Background(), e)
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, created)
}

func (h *ExhibitionHandler) ChangeStatus(c *gin.Context) {
	var req dto.ExhibitionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	e, err := h.svc.ChangeStatus(context.Background(), c.Param("id"), req.Status)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, e)
}

func (h *ExhibitionHandler) AddArtwork(c *gin.Context) {
	e, err := h.svc.AddArtwork(context.Background(), c.Param("id"), c.Param("artworkId"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, e)
}

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

// ArtworkHandler 作品接口。
type ArtworkHandler struct {
	svc *service.ArtworkService
}

func NewArtworkHandler(svc *service.ArtworkService) *ArtworkHandler { return &ArtworkHandler{svc: svc} }

func (h *ArtworkHandler) List(c *gin.Context) {
	onlyPublished := c.DefaultQuery("published", "true") == "true"
	artistID := c.DefaultQuery("artistId", "")
	list, err := h.svc.List(context.Background(), onlyPublished, artistID)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, list)
}

func (h *ArtworkHandler) Get(c *gin.Context) {
	a, err := h.svc.Get(context.Background(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}

func (h *ArtworkHandler) Create(c *gin.Context) {
	var req dto.ArtworkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	userID, _ := c.Get("userID")
	a := &model.Artwork{
		Title: req.Title, Description: req.Description, Year: req.Year, Medium: req.Medium,
		Materials: req.Materials, Size: model.ArtworkSize{Length: req.Size.Length, Width: req.Size.Width, Height: req.Size.Height},
		ImageURLs: req.ImageURLs, VideoURL: req.VideoURL, Tags: req.Tags, Price: req.Price,
		ArtistID: userID.(string),
	}
	created, err := h.svc.Create(context.Background(), a)
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, created)
}

func (h *ArtworkHandler) ChangeStatus(c *gin.Context) {
	var req dto.ArtworkStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	a, err := h.svc.ChangeStatus(context.Background(), c.Param("id"), req.Status)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}

func (h *ArtworkHandler) Review(c *gin.Context) {
	var req dto.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	userID, _ := c.Get("userID")
	a, err := h.svc.Review(context.Background(), c.Param("id"), req.Result, req.Comment, userID.(string))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}

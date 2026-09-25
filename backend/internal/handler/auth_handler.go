package handler

import (
	"context"
	"net/http"

	"github.com/artvault/artvault/internal/dto"
	"github.com/artvault/artvault/internal/service"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证接口。
type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	token, user, err := h.svc.Login(context.Background(), req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, dto.LoginResponse{Token: token, User: dto.UserSummary{ID: user.ID, Username: user.Username, Name: user.Name, Role: user.Role}})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	user, err := h.svc.Register(context.Background(), req.Username, req.Password, req.Name, req.Role)
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, dto.UserSummary{ID: user.ID, Username: user.Username, Name: user.Name, Role: user.Role})
}

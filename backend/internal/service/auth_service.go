package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/artvault/artvault/internal/config"
	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/model"
	"github.com/artvault/artvault/internal/repository"
	"github.com/artvault/artvault/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 注册/登录。
type AuthService struct {
	userRepo *repository.UserRepository
	logger   *slog.Logger
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, logger *slog.Logger, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, logger: logger, cfg: cfg}
}

func (s *AuthService) Register(ctx context.Context, username, password, name, role string) (*model.User, error) {
	if _, err := s.userRepo.FindByUsername(ctx, username); err == nil {
		return nil, util.NewAppError(constants.CodeUserExists, "用户名已存在", nil)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	user := &model.User{ID: util.NewID("user"), Username: username, Password: string(hash), Name: name, Role: role, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, *model.User, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if errors.Is(err, repository.ErrNotFound) {
		return "", nil, util.NewAppError(constants.CodeInvalidCredentials, "用户名或密码错误", nil)
	}
	if err != nil {
		return "", nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", nil, util.NewAppError(constants.CodeInvalidCredentials, "用户名或密码错误", nil)
	}
	token, err := util.GenerateToken(s.cfg.JWTSecret, user.ID, user.Role, s.cfg.TokenExpire)
	if err != nil {
		return "", nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogUserLoggedIn, "userID", user.ID)
	return token, user, nil
}

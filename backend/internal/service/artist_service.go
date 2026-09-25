package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/model"
	"github.com/artvault/artvault/internal/repository"
	"github.com/artvault/artvault/internal/util"
)

// ArtistService 艺术家业务逻辑。
type ArtistService struct {
	repo   *repository.ArtistRepository
	logger *slog.Logger
}

func NewArtistService(repo *repository.ArtistRepository, logger *slog.Logger) *ArtistService {
	return &ArtistService{repo: repo, logger: logger}
}

func (s *ArtistService) Get(ctx context.Context, id string) (*model.Artist, error) {
	a, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(constants.CodeNotFound, "艺术家不存在", nil)
	}
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return a, nil
}

func (s *ArtistService) List(ctx context.Context) ([]model.Artist, error) {
	list, err := s.repo.List(ctx, "Active")
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return list, nil
}

// UpsertForUser 为用户创建/更新艺术家主页。
func (s *ArtistService) UpsertForUser(ctx context.Context, userID string, artistName, bio, avatarURL string, mediums []string) (*model.Artist, error) {
	a, err := s.repo.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		a = &model.Artist{ID: util.NewID("artist"), UserID: userID, ArtistName: artistName, Bio: bio, AvatarURL: avatarURL, Mediums: mediums, SocialLinks: map[string]string{}, Status: "Active", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		if err := s.repo.Create(ctx, a); err != nil {
			return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		return a, nil
	}
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if artistName != "" {
		a.ArtistName = artistName
	}
	if bio != "" {
		a.Bio = bio
	}
	if avatarURL != "" {
		a.AvatarURL = avatarURL
	}
	if len(mediums) > 0 {
		a.Mediums = mediums
	}
	a.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogArtistProfile, "artistID", a.ID)
	return a, nil
}

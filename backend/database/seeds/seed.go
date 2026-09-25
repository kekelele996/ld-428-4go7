package seeds

import (
	"context"
	"log/slog"
	"time"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Seed 写入演示种子数据。
func Seed(ctx context.Context, db *mongo.Database, logger *slog.Logger) {
	hash := func(p string) string {
		h, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		return string(h)
	}
	users := []model.User{
		{ID: "user-admin", Username: "admin", Password: hash("Admin@123"), Name: "平台管理员", Role: constants.RoleAdmin},
		{ID: "user-lin", Username: "lin", Password: hash("Artist@123"), Name: "林知微", Role: constants.RoleArtist},
		{ID: "user-chen", Username: "chen", Password: hash("Artist@123"), Name: "陈序", Role: constants.RoleArtist},
		{ID: "user-viewer", Username: "viewer", Password: hash("Viewer@123"), Name: "观众小艺", Role: constants.RoleViewer},
	}
	for i := range users {
		users[i].CreatedAt = time.Now()
		users[i].UpdatedAt = time.Now()
		n, _ := db.Collection("users").CountDocuments(ctx, bson.M{"_id": users[i].ID})
		if n == 0 {
			db.Collection("users").InsertOne(ctx, users[i])
		}
	}

	artistCount, _ := db.Collection("artists").CountDocuments(ctx, bson.M{})
	if artistCount == 0 {
		now := time.Now()
		db.Collection("artists").InsertOne(ctx, model.Artist{ID: "artist-lin", UserID: "user-lin", ArtistName: "林知微", Bio: "以东方纸本纹理和城市拆迁现场为线索，长期记录空间记忆在材料上的残留。", AvatarURL: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=crop&w=400&q=80", Mediums: []string{constants.MediumInk, constants.MediumMixed, constants.MediumInstallation}, FeaturedWorkURLs: []string{"https://images.unsplash.com/photo-1547891654-e66ed7ebb968?auto=format&fit=crop&w=900&q=80"}, SocialLinks: map[string]string{"website": "https://example.com/lin"}, FollowerCount: 1280, Status: "Active", CreatedAt: now, UpdatedAt: now})
		db.Collection("artists").InsertOne(ctx, model.Artist{ID: "artist-chen", UserID: "user-chen", ArtistName: "陈序", Bio: "关注夜间人工光源与人造景观，作品横跨摄影、数字绘画与空间影像。", AvatarURL: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=400&q=80", Mediums: []string{constants.MediumPhotography, constants.MediumDigital}, FeaturedWorkURLs: []string{"https://images.unsplash.com/photo-1515405295579-ba7b45403062?auto=format&fit=crop&w=900&q=80"}, SocialLinks: map[string]string{"instagram": "https://example.com/chen"}, FollowerCount: 2460, Status: "Active", CreatedAt: now, UpdatedAt: now})
	}

	artCount, _ := db.Collection("artworks").CountDocuments(ctx, bson.M{})
	if artCount == 0 {
		now := time.Now()
		db.Collection("artworks").InsertOne(ctx, model.Artwork{ID: "art-001", Title: "墙体记忆 No.7", Description: "纸本、矿物颜料与拆除墙皮拓印组成的混合材料作品。", Year: 2025, Medium: constants.MediumMixed, Materials: "纸本、矿物颜料、石灰墙皮", Size: model.ArtworkSize{Length: 120, Width: 90, Height: 4}, ImageURLs: []string{"https://images.unsplash.com/photo-1579783900882-c0d3dad7b119?auto=format&fit=crop&w=1200&q=80"}, Tags: []string{"城市记忆", "混合材料", "纸本"}, ArtistID: "artist-lin", ExhibitionIDs: []string{"exh-001"}, Status: constants.ArtworkPublished, Price: 32000, Views: 4820, Likes: 318, Bookmarks: 92, ReviewStatus: "Approved", CreatedAt: now, UpdatedAt: now})
		db.Collection("artworks").InsertOne(ctx, model.Artwork{ID: "art-002", Title: "低照度花园", Description: "以长曝光摄影和数字上色重建城市边缘的夜间绿地。", Year: 2024, Medium: constants.MediumPhotography, Materials: "艺术微喷、铝塑板", Size: model.ArtworkSize{Length: 80, Width: 120}, ImageURLs: []string{"https://images.unsplash.com/photo-1531913764164-f85c52e6e654?auto=format&fit=crop&w=1200&q=80"}, Tags: []string{"摄影", "夜景", "人工光"}, ArtistID: "artist-chen", ExhibitionIDs: []string{"exh-001"}, Status: constants.ArtworkPublished, Price: 18000, Views: 3910, Likes: 274, Bookmarks: 71, ReviewStatus: "Approved", CreatedAt: now, UpdatedAt: now})
		db.Collection("artworks").InsertOne(ctx, model.Artwork{ID: "art-003", Title: "临时入口", Description: "由回收门框和半透明织物构成，可根据展厅动线重组。", Year: 2026, Medium: constants.MediumInstallation, Materials: "回收木门框、织物、冷光源", Size: model.ArtworkSize{Length: 260, Width: 180, Height: 230}, ImageURLs: []string{"https://images.unsplash.com/photo-1561214115-f2f134cc4912?auto=format&fit=crop&w=1200&q=80"}, Tags: []string{"装置", "空间", "回收材料"}, ArtistID: "artist-lin", ExhibitionIDs: []string{"exh-002"}, Status: constants.ArtworkDraft, Views: 870, Likes: 54, Bookmarks: 18, ReviewStatus: "Pending", CreatedAt: now, UpdatedAt: now})
	}

	exhCount, _ := db.Collection("exhibitions").CountDocuments(ctx, bson.M{})
	if exhCount == 0 {
		now := time.Now()
		db.Collection("exhibitions").InsertOne(ctx, model.Exhibition{ID: "exh-001", Title: "材料仍在说话", Description: "一组围绕城市更新、人工光和材料记忆展开的小型群展。", CuratorID: "artist-lin", StartDate: "2026-05-01", EndDate: "2026-08-15", Type: constants.ExhibitionTypeGroup, CoverURL: "https://images.unsplash.com/photo-1545987796-200677ee1011?auto=format&fit=crop&w=1400&q=80", ArtworkIDs: []string{"art-001", "art-002"}, Status: constants.ExhibitionActive, Visitors: 12840, ReviewStatus: "Approved", CreatedAt: now, UpdatedAt: now})
		db.Collection("exhibitions").InsertOne(ctx, model.Exhibition{ID: "exh-002", Title: "可拆卸的房间", Description: "林知微个展，讨论临时结构如何改变观看关系。", CuratorID: "artist-lin", StartDate: "2026-09-10", EndDate: "2026-11-30", Type: constants.ExhibitionTypeSolo, CoverURL: "https://images.unsplash.com/photo-1564399579883-451a5d44ec08?auto=format&fit=crop&w=1400&q=80", ArtworkIDs: []string{"art-003"}, Status: constants.ExhibitionPlanning, ReviewStatus: "Approved", CreatedAt: now, UpdatedAt: now})
	}

	logger.Info(constants.LogSeedDone, "project", "artvault")
}

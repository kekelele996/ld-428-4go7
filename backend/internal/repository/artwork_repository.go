package repository

import (
	"context"
	"fmt"

	"github.com/artvault/artvault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ArtworkRepository 作品数据访问。
type ArtworkRepository struct {
	coll *mongo.Collection
}

func NewArtworkRepository(db *mongo.Database) *ArtworkRepository {
	return &ArtworkRepository{coll: db.Collection("artworks")}
}

func (r *ArtworkRepository) Create(ctx context.Context, a *model.Artwork) error {
	_, err := r.coll.InsertOne(ctx, a)
	if err != nil {
		return fmt.Errorf("insert artwork: %w", err)
	}
	return nil
}

func (r *ArtworkRepository) FindByID(ctx context.Context, id string) (*model.Artwork, error) {
	var a model.Artwork
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&a)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find artwork: %w", err)
	}
	return &a, nil
}

// List 列表；onlyPublished 时仅返回已发布且审核通过的作品。
func (r *ArtworkRepository) List(ctx context.Context, onlyPublished bool, artistID string) ([]model.Artwork, error) {
	filter := bson.M{}
	if onlyPublished {
		filter["status"] = "Published"
		filter["review_status"] = "Approved"
	}
	if artistID != "" {
		filter["artist_id"] = artistID
	}
	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("list artworks: %w", err)
	}
	defer cursor.Close(ctx)
	list := []model.Artwork{}
	if err := cursor.All(ctx, &list); err != nil {
		return nil, fmt.Errorf("decode artworks: %w", err)
	}
	return list, nil
}

func (r *ArtworkRepository) Update(ctx context.Context, a *model.Artwork) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": a.ID}, a)
	if err != nil {
		return fmt.Errorf("update artwork: %w", err)
	}
	return nil
}

// IncCounters 增减浏览/点赞/收藏计数。
func (r *ArtworkRepository) IncCounters(ctx context.Context, id string, views, likes, bookmarks int64) error {
	update := bson.M{"$inc": bson.M{}}
	if views != 0 {
		update["$inc"].(bson.M)["views"] = views
	}
	if likes != 0 {
		update["$inc"].(bson.M)["likes"] = likes
	}
	if bookmarks != 0 {
		update["$inc"].(bson.M)["bookmarks"] = bookmarks
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("inc artwork counters: %w", err)
	}
	return nil
}

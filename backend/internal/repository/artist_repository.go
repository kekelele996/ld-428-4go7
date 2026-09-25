package repository

import (
	"context"
	"fmt"

	"github.com/artvault/artvault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ArtistRepository 艺术家数据访问。
type ArtistRepository struct {
	coll *mongo.Collection
}

func NewArtistRepository(db *mongo.Database) *ArtistRepository {
	return &ArtistRepository{coll: db.Collection("artists")}
}

func (r *ArtistRepository) Create(ctx context.Context, a *model.Artist) error {
	_, err := r.coll.InsertOne(ctx, a)
	if err != nil {
		return fmt.Errorf("insert artist: %w", err)
	}
	return nil
}

func (r *ArtistRepository) FindByID(ctx context.Context, id string) (*model.Artist, error) {
	var a model.Artist
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&a)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find artist: %w", err)
	}
	return &a, nil
}

func (r *ArtistRepository) FindByUserID(ctx context.Context, userID string) (*model.Artist, error) {
	var a model.Artist
	err := r.coll.FindOne(ctx, bson.M{"user_id": userID}).Decode(&a)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find artist by user: %w", err)
	}
	return &a, nil
}

func (r *ArtistRepository) List(ctx context.Context, status string) ([]model.Artist, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "follower_count", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("list artists: %w", err)
	}
	defer cursor.Close(ctx)
	list := []model.Artist{}
	if err := cursor.All(ctx, &list); err != nil {
		return nil, fmt.Errorf("decode artists: %w", err)
	}
	return list, nil
}

func (r *ArtistRepository) Update(ctx context.Context, a *model.Artist) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": a.ID}, a)
	if err != nil {
		return fmt.Errorf("update artist: %w", err)
	}
	return nil
}

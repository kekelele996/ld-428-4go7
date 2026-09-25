package repository

import (
	"context"
	"fmt"

	"github.com/artvault/artvault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ExhibitionRepository 展览数据访问。
type ExhibitionRepository struct {
	coll *mongo.Collection
}

func NewExhibitionRepository(db *mongo.Database) *ExhibitionRepository {
	return &ExhibitionRepository{coll: db.Collection("exhibitions")}
}

func (r *ExhibitionRepository) Create(ctx context.Context, e *model.Exhibition) error {
	_, err := r.coll.InsertOne(ctx, e)
	if err != nil {
		return fmt.Errorf("insert exhibition: %w", err)
	}
	return nil
}

func (r *ExhibitionRepository) FindByID(ctx context.Context, id string) (*model.Exhibition, error) {
	var e model.Exhibition
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&e)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find exhibition: %w", err)
	}
	return &e, nil
}

func (r *ExhibitionRepository) List(ctx context.Context, onlyPublished bool) ([]model.Exhibition, error) {
	filter := bson.M{}
	if onlyPublished {
		filter["status"] = bson.M{"$ne": "Archived"}
		filter["review_status"] = "Approved"
	}
	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("list exhibitions: %w", err)
	}
	defer cursor.Close(ctx)
	list := []model.Exhibition{}
	if err := cursor.All(ctx, &list); err != nil {
		return nil, fmt.Errorf("decode exhibitions: %w", err)
	}
	return list, nil
}

func (r *ExhibitionRepository) Update(ctx context.Context, e *model.Exhibition) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": e.ID}, e)
	if err != nil {
		return fmt.Errorf("update exhibition: %w", err)
	}
	return nil
}

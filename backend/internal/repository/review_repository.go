package repository

import (
	"context"
	"fmt"

	"github.com/artvault/artvault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReviewRepository 内容审核数据访问。
type ReviewRepository struct {
	coll *mongo.Collection
}

func NewReviewRepository(db *mongo.Database) *ReviewRepository {
	return &ReviewRepository{coll: db.Collection("review_logs")}
}

func (r *ReviewRepository) Create(ctx context.Context, log *model.ReviewLog) error {
	_, err := r.coll.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("insert review log: %w", err)
	}
	return nil
}

func (r *ReviewRepository) List(ctx context.Context, targetType string) ([]model.ReviewLog, error) {
	filter := bson.M{}
	if targetType != "" {
		filter["target_type"] = targetType
	}
	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "reviewed_at", Value: -1}}).SetLimit(200))
	if err != nil {
		return nil, fmt.Errorf("list review logs: %w", err)
	}
	defer cursor.Close(ctx)
	list := []model.ReviewLog{}
	if err := cursor.All(ctx, &list); err != nil {
		return nil, fmt.Errorf("decode review logs: %w", err)
	}
	return list, nil
}

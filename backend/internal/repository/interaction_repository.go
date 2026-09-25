package repository

import (
	"context"
	"fmt"

	"github.com/artvault/artvault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InteractionRepository 互动数据访问。
type InteractionRepository struct {
	coll *mongo.Collection
}

func NewInteractionRepository(db *mongo.Database) *InteractionRepository {
	return &InteractionRepository{coll: db.Collection("interactions")}
}

func (r *InteractionRepository) Create(ctx context.Context, i *model.Interaction) error {
	_, err := r.coll.InsertOne(ctx, i)
	if err != nil {
		return fmt.Errorf("insert interaction: %w", err)
	}
	return nil
}

func (r *InteractionRepository) List(ctx context.Context, targetType, targetID string) ([]model.Interaction, error) {
	filter := bson.M{}
	if targetType != "" {
		filter["target_type"] = targetType
	}
	if targetID != "" {
		filter["target_id"] = targetID
	}
	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(200))
	if err != nil {
		return nil, fmt.Errorf("list interactions: %w", err)
	}
	defer cursor.Close(ctx)
	list := []model.Interaction{}
	if err := cursor.All(ctx, &list); err != nil {
		return nil, fmt.Errorf("decode interactions: %w", err)
	}
	return list, nil
}

// DeleteByUserTargetType 取消点赞/收藏。
func (r *InteractionRepository) DeleteByUserTargetType(ctx context.Context, userID, targetType, targetID, itype string) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"user_id": userID, "target_type": targetType, "target_id": targetID, "type": itype})
	if err != nil {
		return fmt.Errorf("delete interaction: %w", err)
	}
	return nil
}

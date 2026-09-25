package repository

import (
	"context"
	"fmt"

	"github.com/artvault/artvault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuditLogRepository 操作日志数据访问。
type AuditLogRepository struct {
	coll *mongo.Collection
}

func NewAuditLogRepository(db *mongo.Database) *AuditLogRepository {
	return &AuditLogRepository{coll: db.Collection("audit_logs")}
}

func (r *AuditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	_, err := r.coll.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

func (r *AuditLogRepository) List(ctx context.Context, limit int64) ([]model.AuditLog, error) {
	cursor, err := r.coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit))
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	defer cursor.Close(ctx)
	list := []model.AuditLog{}
	if err := cursor.All(ctx, &list); err != nil {
		return nil, fmt.Errorf("decode audit logs: %w", err)
	}
	return list, nil
}

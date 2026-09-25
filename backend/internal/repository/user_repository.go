package repository

import (
	"context"
	"fmt"

	"github.com/artvault/artvault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrNotFound 哨兵错误。
var ErrNotFound = fmt.Errorf("not found")

// UserRepository 用户数据访问（mongo-driver 实现）。
type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection("users")}
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	_, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &u, nil
}

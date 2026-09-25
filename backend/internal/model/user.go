package model

import "time"

// User 用户。
type User struct {
	ID        string    `bson:"_id" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"-"`
	Name      string    `bson:"name" json:"name"`
	Role      string    `bson:"role" json:"role"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

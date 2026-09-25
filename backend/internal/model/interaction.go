package model

import "time"

// Interaction 互动记录。
type Interaction struct {
	ID         string    `bson:"_id" json:"id"`
	UserID     string    `bson:"user_id" json:"userId"`
	TargetType string    `bson:"target_type" json:"targetType"`
	TargetID   string    `bson:"target_id" json:"targetId"`
	Type       string    `bson:"type" json:"type"`
	Comment    string    `bson:"comment" json:"comment,omitempty"`
	CreatedAt  time.Time `bson:"created_at" json:"createdAt"`
}

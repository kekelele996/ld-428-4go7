package model

import "time"

// ReviewLog 内容审核日志。
type ReviewLog struct {
	ID         string    `bson:"_id" json:"id"`
	TargetType string    `bson:"target_type" json:"targetType"`
	TargetID   string    `bson:"target_id" json:"targetId"`
	ReviewerID string    `bson:"reviewer_id" json:"reviewerId"`
	Result     string    `bson:"result" json:"result"`
	Comment    string    `bson:"comment" json:"comment"`
	ReviewedAt time.Time `bson:"reviewed_at" json:"reviewedAt"`
}

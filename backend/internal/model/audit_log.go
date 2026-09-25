package model

import "time"

// AuditLog 操作日志。
type AuditLog struct {
	ID        string    `bson:"_id" json:"id"`
	UserID    string    `bson:"user_id" json:"userId"`
	Action    string    `bson:"action" json:"action"`
	Entity    string    `bson:"entity" json:"entity"`
	EntityID  string    `bson:"entity_id" json:"entityId"`
	Detail    string    `bson:"detail" json:"detail"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
}

package dto

// InteractionCreateRequest 互动。
type InteractionCreateRequest struct {
	TargetType string `json:"targetType" binding:"required,oneof=Artwork Exhibition"`
	TargetID   string `json:"targetId" binding:"required"`
	Type       string `json:"type" binding:"required,oneof=Like Comment Bookmark Share"`
	Comment    string `json:"comment"`
}

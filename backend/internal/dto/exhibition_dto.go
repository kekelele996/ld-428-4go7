package dto

// ExhibitionCreateRequest 创建展览。
type ExhibitionCreateRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Type        string   `json:"type" binding:"oneof=Solo Group Thematic Permanent"`
	CoverURL    string   `json:"coverUrl"`
	ArtworkIDs  []string `json:"artworkIds"`
}

// ExhibitionStatusRequest 发布展览。
type ExhibitionStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=Planning Active Ended Archived"`
}

// ExhibitionReadinessResponse 展览公开就绪评估结果。
type ExhibitionReadinessResponse struct {
	ExhibitionID       string   `json:"exhibitionId"`
	TotalArtworks      int      `json:"totalArtworks"`
	PublicReadyCount   int      `json:"publicReadyCount"`
	PendingReviewCount int      `json:"pendingReviewCount"`
	RejectedCount      int      `json:"rejectedCount"`
	DraftCount         int      `json:"draftCount"`
	SoldCount          int      `json:"soldCount"`
	ArchivedCount      int      `json:"archivedCount"`
	NoImageCount       int      `json:"noImageCount"`
	MissingCount       int      `json:"missingCount"`
	ReadyRatio         float64  `json:"readyRatio"`
	Ready              bool     `json:"ready"`
	BlockingReasons    []string `json:"blockingReasons"`
}

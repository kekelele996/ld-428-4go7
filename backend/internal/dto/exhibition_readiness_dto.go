package dto

// ExhibitionReadinessResponse GET /api/v1/exhibitions/:id/readiness 响应。
// 按关联作品实时统计，任何异常都作为阻断原因返回。
type ExhibitionReadinessResponse struct {
	ExhibitionID  string                          `json:"exhibitionId"`
	Exhibition    *ExhibitionBriefResponse        `json:"exhibition,omitempty"`
	Conclusion    string                          `json:"conclusion"`
	Ready         bool                            `json:"ready"`
	IsEmpty       bool                            `json:"isEmpty"`
	Total         int                             `json:"total"`
	ReadyCount    int                             `json:"readyCount"`
	PendingCount  int                             `json:"pendingCount"`
	SoldCount     int                             `json:"soldCount"`
	ArchivedCount int                             `json:"archivedCount"`
	NoImageCount  int                             `json:"noImageCount"`
	ReadyRatio    float64                         `json:"readyRatio"`
	Reasons       []ExhibitionReadinessReasonItem `json:"reasons"`
}

// ExhibitionReadinessReasonItem 阻断原因项。
type ExhibitionReadinessReasonItem struct {
	Code       string   `json:"code"`
	Message    string   `json:"message"`
	Count      int      `json:"count"`
	ArtworkIDs []string `json:"artworkIds"`
}

// ExhibitionBriefResponse 就绪接口内嵌的展览摘要。
type ExhibitionBriefResponse struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	ReviewStatus string `json:"reviewStatus"`
}

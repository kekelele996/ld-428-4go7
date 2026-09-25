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

package dto

// SizeInput 尺寸。
type SizeInput struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height,omitempty"`
}

// ArtworkCreateRequest 上传作品。
type ArtworkCreateRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Year        int       `json:"year"`
	Medium      string    `json:"medium"`
	Materials   string    `json:"materials"`
	Size        SizeInput `json:"size"`
	ImageURLs   []string  `json:"imageUrls"`
	VideoURL    string    `json:"videoUrl"`
	Tags        []string  `json:"tags"`
	Price       float64   `json:"price"`
}

// ArtworkStatusRequest 发布/下架。
type ArtworkStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=Draft Published Sold Archived"`
}

// ReviewRequest 内容审核。
type ReviewRequest struct {
	Result  string `json:"result" binding:"required,oneof=Approved Rejected Flagged"`
	Comment string `json:"comment"`
}

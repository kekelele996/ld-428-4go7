package model

// ReadinessReasonCode 展览就绪阻断原因码。
type ReadinessReasonCode string

// ExhibitionReadiness 展览公开展出就绪统计（按关联作品实时计算）。
type ExhibitionReadiness struct {
	ExhibitionID  string
	Exhibition    *Exhibition
	Total         int
	ReadyCount    int
	PendingCount  int
	SoldCount     int
	ArchivedCount int
	NoImageCount  int
	ReadyRatio    float64
	Ready         bool
	IsEmpty       bool
	Reasons       []ReadinessBlockReason
}

// ReadinessBlockReason 单个阻断原因及命中作品。
type ReadinessBlockReason struct {
	Code       ReadinessReasonCode
	Message    string
	Count      int
	ArtworkIDs []string
}

package constants

// ExhibitionStatus 展览状态枚举。
const (
	ExhibitionPlanning = "Planning"
	ExhibitionActive   = "Active"
	ExhibitionEnded    = "Ended"
	ExhibitionArchived = "Archived"
)

// ExhibitionType 展览类型。
const (
	ExhibitionTypeSolo      = "Solo"
	ExhibitionTypeGroup     = "Group"
	ExhibitionTypeThematic  = "Thematic"
	ExhibitionTypePermanent = "Permanent"
)

// ReviewStatus 内容审核状态。
const (
	ReviewPending = "Pending"
)

// ExhibitionReadinessReason 展览就绪阻断原因码。
const (
	ReadinessReasonEmptyExhibition = "EMPTY_EXHIBITION"
	ReadinessReasonMissingArtwork  = "MISSING_ARTWORK"
	ReadinessReasonDraft           = "ARTWORK_DRAFT"
	ReadinessReasonPendingReview   = "REVIEW_PENDING"
	ReadinessReasonRejected        = "REVIEW_REJECTED"
	ReadinessReasonFlagged         = "REVIEW_FLAGGED"
	ReadinessReasonReviewUnknown   = "REVIEW_UNKNOWN"
	ReadinessReasonSold            = "ARTWORK_SOLD"
	ReadinessReasonArchived        = "ARTWORK_ARCHIVED"
	ReadinessReasonStatusUnknown   = "ARTWORK_STATUS_UNKNOWN"
	ReadinessReasonNoImage         = "ARTWORK_NO_IMAGE"
)

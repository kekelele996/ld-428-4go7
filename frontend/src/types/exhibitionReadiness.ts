// ExhibitionReadinessReason 展览公开展出阻断原因。
export interface ExhibitionReadinessReason {
  code: string;
  message: string;
  count: number;
  artworkIds: string[];
}

// ExhibitionReadiness GET /api/v1/exhibitions/:id/readiness 响应。
// 按关联作品实时统计，任何异常都作为阻断原因。
export interface ExhibitionReadiness {
  exhibitionId: string;
  exhibition?: {
    id: string;
    title: string;
    status: string;
    reviewStatus: string;
  };
  conclusion: string;
  ready: boolean;
  isEmpty: boolean;
  total: number;
  readyCount: number;
  pendingCount: number;
  soldCount: number;
  archivedCount: number;
  noImageCount: number;
  readyRatio: number;
  reasons: ExhibitionReadinessReason[];
}

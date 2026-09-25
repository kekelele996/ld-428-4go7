import { ExhibitionStatus, ExhibitionType } from './enums';

export interface Exhibition {
  id: string;
  title: string;
  description: string;
  curatorId: string;
  startDate: string;
  endDate: string;
  type: ExhibitionType;
  coverUrl: string;
  artworkIds: string[];
  status: ExhibitionStatus;
  visitors: number;
}

export interface ExhibitionReadiness {
  exhibitionId: string;
  totalArtworks: number;
  publicReadyCount: number;
  pendingReviewCount: number;
  rejectedCount: number;
  draftCount: number;
  soldCount: number;
  archivedCount: number;
  noImageCount: number;
  missingCount: number;
  readyRatio: number;
  ready: boolean;
  blockingReasons: string[];
}

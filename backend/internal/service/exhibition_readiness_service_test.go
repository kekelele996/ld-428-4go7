package service

import (
	"testing"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/model"
)

func readyArtwork(id string) model.Artwork {
	return model.Artwork{
		ID:           id,
		Status:       constants.ArtworkPublished,
		ReviewStatus: constants.ReviewApproved,
		ImageURLs:    []string{"https://example.com/" + id + ".jpg"},
	}
}

func reasonCodes(r *model.ExhibitionReadiness) []string {
	codes := make([]string, 0, len(r.Reasons))
	for _, reason := range r.Reasons {
		codes = append(codes, string(reason.Code))
	}
	return codes
}

func TestEvaluateReadiness(t *testing.T) {
	tests := []struct {
		name            string
		artworkIDs      []string
		artworks        []model.Artwork
		wantReady       bool
		wantEmpty       bool
		wantTotal       int
		wantReadyCount  int
		wantPending     int
		wantSold        int
		wantArchived    int
		wantNoImage     int
		wantRatio       float64
		wantReasonCodes []string
	}{
		{
			name:           "空展览给说明",
			artworkIDs:     nil,
			wantReady:      false,
			wantEmpty:      true,
			wantTotal:      0,
			wantReadyCount: 0,
			wantRatio:      0,
		},
		{
			name:           "全部就绪",
			artworkIDs:     []string{"a1", "a2"},
			artworks:       []model.Artwork{readyArtwork("a1"), readyArtwork("a2")},
			wantReady:      true,
			wantTotal:      2,
			wantReadyCount: 2,
			wantRatio:      1,
		},
		{
			name:            "草稿且待审核",
			artworkIDs:      []string{"a1"},
			artworks:        []model.Artwork{{ID: "a1", Status: constants.ArtworkDraft, ReviewStatus: constants.ReviewPending, ImageURLs: []string{"x"}}},
			wantTotal:       1,
			wantPending:     1,
			wantReasonCodes: []string{constants.ReadinessReasonDraft, constants.ReadinessReasonPendingReview},
		},
		{
			name:            "已售且无图命中两个原因",
			artworkIDs:      []string{"a1"},
			artworks:        []model.Artwork{{ID: "a1", Status: constants.ArtworkSold, ReviewStatus: constants.ReviewApproved}},
			wantTotal:       1,
			wantSold:        1,
			wantNoImage:     1,
			wantReasonCodes: []string{constants.ReadinessReasonSold, constants.ReadinessReasonNoImage},
		},
		{
			name:            "下架",
			artworkIDs:      []string{"a1"},
			artworks:        []model.Artwork{{ID: "a1", Status: constants.ArtworkArchived, ReviewStatus: constants.ReviewApproved, ImageURLs: []string{"x"}}},
			wantTotal:       1,
			wantArchived:    1,
			wantReasonCodes: []string{constants.ReadinessReasonArchived},
		},
		{
			name:            "审核驳回与标记",
			artworkIDs:      []string{"a1", "a2"},
			artworks:        []model.Artwork{{ID: "a1", Status: constants.ArtworkPublished, ReviewStatus: constants.ReviewRejected, ImageURLs: []string{"x"}}, {ID: "a2", Status: constants.ArtworkPublished, ReviewStatus: constants.ReviewFlagged, ImageURLs: []string{"x"}}},
			wantTotal:       2,
			wantReasonCodes: []string{constants.ReadinessReasonRejected, constants.ReadinessReasonFlagged},
		},
		{
			name:            "关联引用缺失",
			artworkIDs:      []string{"gone"},
			artworks:        nil,
			wantTotal:       1,
			wantReasonCodes: []string{constants.ReadinessReasonMissingArtwork},
		},
		{
			name:            "未知审核状态",
			artworkIDs:      []string{"a1"},
			artworks:        []model.Artwork{{ID: "a1", Status: constants.ArtworkPublished, ReviewStatus: "Bizarre", ImageURLs: []string{"x"}}},
			wantTotal:       1,
			wantReasonCodes: []string{constants.ReadinessReasonReviewUnknown},
		},
		{
			name:            "未知作品状态",
			artworkIDs:      []string{"a1"},
			artworks:        []model.Artwork{{ID: "a1", Status: "Transit", ReviewStatus: constants.ReviewApproved, ImageURLs: []string{"x"}}},
			wantTotal:       1,
			wantReasonCodes: []string{constants.ReadinessReasonStatusUnknown},
		},
		{
			name:            "图片全为空白字符串视为无图",
			artworkIDs:      []string{"a1"},
			artworks:        []model.Artwork{{ID: "a1", Status: constants.ArtworkPublished, ReviewStatus: constants.ReviewApproved, ImageURLs: []string{"  "}}},
			wantTotal:       1,
			wantNoImage:     1,
			wantReasonCodes: []string{constants.ReadinessReasonNoImage},
		},
		{
			name:            "部分就绪比例 0.5",
			artworkIDs:      []string{"a1", "a2"},
			artworks:        []model.Artwork{readyArtwork("a1"), {ID: "a2", Status: constants.ArtworkSold, ReviewStatus: constants.ReviewApproved, ImageURLs: []string{"x"}}},
			wantTotal:       2,
			wantReadyCount:  1,
			wantSold:        1,
			wantRatio:       0.5,
			wantReady:       false,
			wantReasonCodes: []string{constants.ReadinessReasonSold},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &model.Exhibition{ID: "exh-1", ArtworkIDs: tt.artworkIDs}
			got := evaluateReadiness(e, tt.artworks)

			if got.Ready != tt.wantReady {
				t.Errorf("Ready = %v, want %v", got.Ready, tt.wantReady)
			}
			if got.IsEmpty != tt.wantEmpty {
				t.Errorf("IsEmpty = %v, want %v", got.IsEmpty, tt.wantEmpty)
			}
			if got.Total != tt.wantTotal {
				t.Errorf("Total = %d, want %d", got.Total, tt.wantTotal)
			}
			if got.ReadyCount != tt.wantReadyCount {
				t.Errorf("ReadyCount = %d, want %d", got.ReadyCount, tt.wantReadyCount)
			}
			if got.PendingCount != tt.wantPending {
				t.Errorf("PendingCount = %d, want %d", got.PendingCount, tt.wantPending)
			}
			if got.SoldCount != tt.wantSold {
				t.Errorf("SoldCount = %d, want %d", got.SoldCount, tt.wantSold)
			}
			if got.ArchivedCount != tt.wantArchived {
				t.Errorf("ArchivedCount = %d, want %d", got.ArchivedCount, tt.wantArchived)
			}
			if got.NoImageCount != tt.wantNoImage {
				t.Errorf("NoImageCount = %d, want %d", got.NoImageCount, tt.wantNoImage)
			}
			if got.ReadyRatio != tt.wantRatio {
				t.Errorf("ReadyRatio = %v, want %v", got.ReadyRatio, tt.wantRatio)
			}
			codes := reasonCodes(got)
			if len(codes) != len(tt.wantReasonCodes) {
				t.Fatalf("reason codes = %v, want %v", codes, tt.wantReasonCodes)
			}
			for i, code := range tt.wantReasonCodes {
				if codes[i] != code {
					t.Errorf("reason code[%d] = %s, want %s (full: %v)", i, codes[i], code, codes)
				}
			}
		})
	}
}

func TestBuildReadinessReasonsAggregatesArtworkIDs(t *testing.T) {
	aggs := map[model.ReadinessReasonCode]*readinessAgg{}
	a := &readinessAgg{}
	a.add("a1", "m1")
	a.add("a2", "m2")
	aggs[constants.ReadinessReasonSold] = a

	reasons := buildReadinessReasons(aggs)
	if len(reasons) != 1 {
		t.Fatalf("reasons len = %d, want 1", len(reasons))
	}
	if reasons[0].Count != 2 {
		t.Errorf("count = %d, want 2", reasons[0].Count)
	}
	if len(reasons[0].ArtworkIDs) != 2 || reasons[0].ArtworkIDs[0] != "a1" || reasons[0].ArtworkIDs[1] != "a2" {
		t.Errorf("artwork ids = %v, want [a1 a2]", reasons[0].ArtworkIDs)
	}
}

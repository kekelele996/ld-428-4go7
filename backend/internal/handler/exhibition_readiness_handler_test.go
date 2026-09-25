package handler

import (
	"testing"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/model"
)

func TestToExhibitionReadinessResponse(t *testing.T) {
	tests := []struct {
		name           string
		input          *model.ExhibitionReadiness
		wantConclusion string
		wantReady      bool
		wantEmpty      bool
		wantReasonLen  int
	}{
		{
			name: "全部就绪",
			input: &model.ExhibitionReadiness{
				Exhibition: &model.Exhibition{ID: "exh-1", Title: "展", Status: constants.ExhibitionActive, ReviewStatus: constants.ReviewApproved},
				Total:      2, ReadyCount: 2, ReadyRatio: 1, Ready: true,
				Reasons: []model.ReadinessBlockReason{},
			},
			wantConclusion: constants.MsgReadinessReady,
			wantReady:      true,
		},
		{
			name: "空展览",
			input: &model.ExhibitionReadiness{
				Exhibition: &model.Exhibition{ID: "exh-2", Title: "空"},
				IsEmpty:    true,
				Reasons:    []model.ReadinessBlockReason{},
			},
			wantConclusion: constants.MsgReadinessEmpty,
			wantEmpty:      true,
		},
		{
			name: "存在阻断",
			input: &model.ExhibitionReadiness{
				Exhibition: &model.Exhibition{ID: "exh-3", Title: "阻"},
				Total:      1, SoldCount: 1, ReadyRatio: 0,
				Reasons: []model.ReadinessBlockReason{{Code: constants.ReadinessReasonSold, Message: "m", Count: 1, ArtworkIDs: []string{"a1"}}},
			},
			wantConclusion: constants.MsgReadinessBlocked,
			wantReasonLen:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toExhibitionReadinessResponse(tt.input)
			if got.Conclusion != tt.wantConclusion {
				t.Errorf("Conclusion = %q, want %q", got.Conclusion, tt.wantConclusion)
			}
			if got.Ready != tt.wantReady {
				t.Errorf("Ready = %v, want %v", got.Ready, tt.wantReady)
			}
			if got.IsEmpty != tt.wantEmpty {
				t.Errorf("IsEmpty = %v, want %v", got.IsEmpty, tt.wantEmpty)
			}
			if len(got.Reasons) != tt.wantReasonLen {
				t.Errorf("Reasons len = %d, want %d", len(got.Reasons), tt.wantReasonLen)
			}
			if got.Exhibition == nil || got.Exhibition.ID != tt.input.Exhibition.ID {
				t.Errorf("Exhibition brief not mapped correctly: %+v", got.Exhibition)
			}
		})
	}
}

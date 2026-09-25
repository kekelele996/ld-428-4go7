package service

import (
	"testing"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/model"
)

func TestEvaluateReadiness(t *testing.T) {
	readyArtwork := func(id string) model.Artwork {
		return model.Artwork{ID: id, Status: constants.ArtworkPublished, ReviewStatus: constants.ReviewApproved, ImageURLs: []string{"https://example.com/a.jpg"}}
	}

	tests := []struct {
		name            string
		artworkIDs      []string
		artworks        []model.Artwork
		wantReady       bool
		wantPublicReady int
		wantPending     int
		wantSold        int
		wantArchived    int
		wantNoImage     int
		wantMissing     int
		wantRatio       float64
		wantReasons     int
	}{
		{
			name:        "empty exhibition is not ready with explanation",
			artworkIDs:  []string{},
			artworks:    []model.Artwork{},
			wantReady:   false,
			wantRatio:   0,
			wantReasons: 1,
		},
		{
			name:            "all artworks ready",
			artworkIDs:      []string{"a1", "a2"},
			artworks:        []model.Artwork{readyArtwork("a1"), readyArtwork("a2")},
			wantReady:       true,
			wantPublicReady: 2,
			wantRatio:       1,
			wantReasons:     0,
		},
		{
			name:       "every anomaly becomes a blocking reason",
			artworkIDs: []string{"a1", "a2", "a3", "a4", "a5", "a6", "a7"},
			artworks: []model.Artwork{
				readyArtwork("a1"),
				{ID: "a2", Status: constants.ArtworkPublished, ReviewStatus: "Pending", ImageURLs: []string{"x"}},
				{ID: "a3", Status: constants.ArtworkPublished, ReviewStatus: constants.ReviewRejected, ImageURLs: []string{"x"}},
				{ID: "a4", Status: constants.ArtworkSold, ReviewStatus: constants.ReviewApproved, ImageURLs: []string{"x"}},
				{ID: "a5", Status: constants.ArtworkArchived, ReviewStatus: constants.ReviewApproved, ImageURLs: []string{"x"}},
				{ID: "a6", Status: constants.ArtworkPublished, ReviewStatus: constants.ReviewApproved},
			},
			wantReady:       false,
			wantPublicReady: 1,
			wantPending:     1,
			wantSold:        1,
			wantArchived:    1,
			wantNoImage:     1,
			wantMissing:     1,
			wantRatio:       1.0 / 7.0,
			wantReasons:     6, // missing + pending + rejected + sold + archived + noImage
		},
		{
			name:        "draft artwork blocks readiness",
			artworkIDs:  []string{"a1"},
			artworks:    []model.Artwork{{ID: "a1", Status: constants.ArtworkDraft, ReviewStatus: "Pending", ImageURLs: []string{"x"}}},
			wantReady:   false,
			wantPending: 1,
			wantRatio:   0,
			wantReasons: 2, // draft + pending
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := evaluateReadiness("exh-test", tt.artworkIDs, tt.artworks)
			if got.Ready != tt.wantReady {
				t.Errorf("Ready = %v, want %v", got.Ready, tt.wantReady)
			}
			if got.PublicReadyCount != tt.wantPublicReady {
				t.Errorf("PublicReadyCount = %d, want %d", got.PublicReadyCount, tt.wantPublicReady)
			}
			if got.PendingReviewCount != tt.wantPending {
				t.Errorf("PendingReviewCount = %d, want %d", got.PendingReviewCount, tt.wantPending)
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
			if got.MissingCount != tt.wantMissing {
				t.Errorf("MissingCount = %d, want %d", got.MissingCount, tt.wantMissing)
			}
			if got.ReadyRatio != tt.wantRatio {
				t.Errorf("ReadyRatio = %v, want %v", got.ReadyRatio, tt.wantRatio)
			}
			if len(got.BlockingReasons) != tt.wantReasons {
				t.Errorf("BlockingReasons = %v, want %d reasons", got.BlockingReasons, tt.wantReasons)
			}
			if got.TotalArtworks != len(tt.artworkIDs) {
				t.Errorf("TotalArtworks = %d, want %d", got.TotalArtworks, len(tt.artworkIDs))
			}
		})
	}
}

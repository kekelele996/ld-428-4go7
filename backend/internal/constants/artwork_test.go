package constants

import "testing"

func TestArtworkStatus(t *testing.T) {
	if ArtworkPublished != "Published" || ArtworkDraft != "Draft" {
		t.Fatal("artwork status mismatch")
	}
	if MediumPhotography != "Photography" || MediumInstallation != "Installation" {
		t.Fatal("medium mismatch")
	}
}

func TestInteractionTypes(t *testing.T) {
	if InteractionLike != "Like" || InteractionBookmark != "Bookmark" {
		t.Fatal("interaction type mismatch")
	}
}

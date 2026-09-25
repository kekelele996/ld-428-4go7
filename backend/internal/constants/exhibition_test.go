package constants

import "testing"

func TestExhibitionStatuses(t *testing.T) {
	cases := []struct {
		got  string
		want string
	}{
		{ExhibitionPlanning, "Planning"},
		{ExhibitionActive, "Active"},
		{ExhibitionEnded, "Ended"},
		{ExhibitionArchived, "Archived"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Fatalf("exhibition status %q want %q", c.got, c.want)
		}
	}
}

func TestExhibitionTypes(t *testing.T) {
	if ExhibitionTypeSolo != "Solo" || ExhibitionTypeGroup != "Group" {
		t.Fatal("exhibition type mismatch")
	}
}

func TestRoles(t *testing.T) {
	if RoleAdmin != "Admin" || RoleCurator != "Curator" || RoleArtist != "Artist" || RoleViewer != "Viewer" {
		t.Fatal("role mismatch")
	}
}

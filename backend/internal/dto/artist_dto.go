package dto

// ArtistUpsertRequest 完善艺术家主页。
type ArtistUpsertRequest struct {
	ArtistName string   `json:"artistName"`
	Bio        string   `json:"bio"`
	AvatarURL  string   `json:"avatarUrl"`
	Mediums    []string `json:"mediums"`
}

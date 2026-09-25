package model

import "time"

// ArtworkSize 尺寸。
type ArtworkSize struct {
	Length float64 `bson:"length" json:"length"`
	Width  float64 `bson:"width" json:"width"`
	Height float64 `bson:"height" json:"height"`
}

// Artwork 作品。
type Artwork struct {
	ID            string      `bson:"_id" json:"id"`
	Title         string      `bson:"title" json:"title"`
	Description   string      `bson:"description" json:"description"`
	Year          int         `bson:"year" json:"year"`
	Medium        string      `bson:"medium" json:"medium"`
	Materials     string      `bson:"materials" json:"materials"`
	Size          ArtworkSize `bson:"size" json:"size"`
	ImageURLs     []string    `bson:"image_urls" json:"imageUrls"`
	VideoURL      string      `bson:"video_url" json:"videoUrl,omitempty"`
	Tags          []string    `bson:"tags" json:"tags"`
	ArtistID      string      `bson:"artist_id" json:"artistId"`
	ExhibitionIDs []string    `bson:"exhibition_ids" json:"exhibitionIds"`
	Status        string      `bson:"status" json:"status"`
	Price         float64     `bson:"price" json:"price,omitempty"`
	Views         int64       `bson:"views" json:"views"`
	Likes         int64       `bson:"likes" json:"likes"`
	Bookmarks     int64       `bson:"bookmarks" json:"bookmarks"`
	ReviewStatus  string      `bson:"review_status" json:"reviewStatus"`
	CreatedAt     time.Time   `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time   `bson:"updated_at" json:"updatedAt"`
}

package model

import "time"

// Exhibition 展览。
type Exhibition struct {
	ID           string    `bson:"_id" json:"id"`
	Title        string    `bson:"title" json:"title"`
	Description  string    `bson:"description" json:"description"`
	CuratorID    string    `bson:"curator_id" json:"curatorId"`
	StartDate    string    `bson:"start_date" json:"startDate"`
	EndDate      string    `bson:"end_date" json:"endDate"`
	Type         string    `bson:"type" json:"type"`
	CoverURL     string    `bson:"cover_url" json:"coverUrl"`
	ArtworkIDs   []string  `bson:"artwork_ids" json:"artworkIds"`
	Status       string    `bson:"status" json:"status"`
	Visitors     int64     `bson:"visitors" json:"visitors"`
	ReviewStatus string    `bson:"review_status" json:"reviewStatus"`
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt"`
}

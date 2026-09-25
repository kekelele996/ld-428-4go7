package model

import "time"

// Artist 艺术家。
type Artist struct {
	ID               string            `bson:"_id" json:"id"`
	UserID           string            `bson:"user_id" json:"userId"`
	ArtistName       string            `bson:"artist_name" json:"artistName"`
	Bio              string            `bson:"bio" json:"bio"`
	AvatarURL        string            `bson:"avatar_url" json:"avatarUrl"`
	Mediums          []string          `bson:"mediums" json:"mediums"`
	FeaturedWorkURLs []string          `bson:"featured_work_urls" json:"featuredWorkUrls"`
	SocialLinks      map[string]string `bson:"social_links" json:"socialLinks"`
	FollowingCount   int64             `bson:"following_count" json:"followingCount"`
	FollowerCount    int64             `bson:"follower_count" json:"followerCount"`
	Status           string            `bson:"status" json:"status"`
	CreatedAt        time.Time         `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time         `bson:"updated_at" json:"updatedAt"`
}

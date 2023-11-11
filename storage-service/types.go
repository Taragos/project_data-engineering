package main

import (
	"time"

	"github.com/gofrs/uuid"
)

type IGUser struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	FollowersCount int `json:"followers_count"`
	FollowsCount   int `json:"follows_count"`

	Media []IGMedia `json:"media"`
}

type IGMedia struct {
	Id uuid.UUID `json:"id"`

	MediaType string `json:"media_type"`
	MediaUrl  string `json:"media_url"`
	Caption   string `json:"caption"`

	IsCommentEnabled bool `json:"is_comment_enabled"`
	CommentsCount    int  `json:"comments_count"`
	LikeCount        int  `json:"like_count"`

	Timestamp time.Time `json:"timestamp"`

	Insights IGInsight `json:"insights"`
}

type IGInsight struct {
	Engagement  int `json:"engagement"`
	Impressions int `json:"impressions"`
	Reach       int `json:"reach"`
	Saved       int `json:"saved"`
}

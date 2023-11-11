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
	Caption   string `json:"caption"`

	IsCommentEnabled bool `json:"is_comment_enabled"`
}

type IGInsight struct {
	Engagement  int `json:"engagement"`
	Impressions int `json:"impressions"`
	Reach       int `json:"reach"`
	Saved       int `json:"saved"`
}

type Insight struct {
	Group     string    `json:"group"`
	Value     int       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}
type FullResult struct {
	Media    IGMedia   `json:"media"`
	Insights []Insight `json:"insights"`
}

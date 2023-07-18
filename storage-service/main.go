package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
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

func main() {
	// 1. Establish Connection to PostgreSQL + Clickhouse
	connStr := "postgresql://postgres:test@localhost:5432/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("could not open connecting to postgres: ", err)
	}

	go func(db *sql.DB) {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{"localhost:9092"},
			Topic:     "instagram-profiles",
			Partition: 0,
			MaxBytes:  10e6,
		})
		defer r.Close()

		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Fatal("could not read message: ", err)
			}

			user := IGUser{}
			json.Unmarshal(m.Value, &user)

			log.Print("test")
			// If User does not exist -> Create
			_, err = db.Exec(`INSERT INTO "ig_profiles"(id, name) values($1, $2) ON CONFLICT DO NOTHING`, user.Id, user.Name)
			if err != nil {
				log.Println(err)
			}

			for idx := range user.Media {
				media := user.Media[idx]
				_, err = db.Exec(`INSERT INTO "ig_media"(id, profile_id, media_type, media_url, caption, is_comment_enabled, ig_created_at) values($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING`,
					media.Id,
					user.Id,
					media.MediaType,
					media.MediaUrl,
					media.Caption,
					media.IsCommentEnabled,
					media.Timestamp,
				)

				if err != nil {
					log.Println(err)
				}
			}
		}
	}(db)

	for {
	}

	// 2. Establish Connection to Kafka
	// 3. Consume Message
	// 3.1 Check if Profile exists 	-> Update/Create
	// 3.2 Check if Medias exist 	-> Update/Create

}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
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

	Timestamp time.Time `json:"timestamp"`

	Insights Insights `json:"insights"`
}

type Insights struct {
	Id          uuid.UUID `json:"id"`
	Comments    int       `json:"comments"`
	Engagement  int       `json:"engagement"`
	Impressions int       `json:"impressions"`
	Likes       int       `json:"likes"`
	Reach       int       `json:"reach"`
	Saved       int       `json:"saved"`
}

func main() {

	kafkaBootstrapServers := flag.String("kafkaBootstrapServers", "", "kafka bootstrap servers address")
	numProfiles := flag.Int("profiles", 2, "number of profiles to simulate")
	profileUpdateFreq := flag.Int("profileUpdateFreq", 5000, "milliseconds to wait in between profile updates")

	flag.Parse()

	log.Println("connection to kafka: ", *kafkaBootstrapServers)

	users := generateProfiles(*numProfiles)

	var wg sync.WaitGroup

	for _, user := range users {
		wg.Add(1)
		go func(user *IGUser, profileUpdateFreq int) {
			conn, err := kafka.DialLeader(context.Background(), "tcp", *kafkaBootstrapServers, "instagram-profiles", 0)
			if err != nil {
				log.Fatal("failed to dial leader:", err)
			}
			defer conn.Close()

			for {
				publishProfile(conn, *user)
				updateProfile(user)
				time.Sleep(time.Duration(profileUpdateFreq) * time.Millisecond)

			}
		}(&user, *profileUpdateFreq)

		for idx := range user.Media {
			wg.Add(1)
			go func(media *IGMedia) {
				conn, err := kafka.DialLeader(context.Background(), "tcp", *kafkaBootstrapServers, "instagram-insights", 0)
				if err != nil {
					log.Fatal("failed to dial leader:", err)
				}
				defer conn.Close()

				for {
					publishMedia(conn, *media)
					updateMedia(media)
					time.Sleep(300 * time.Millisecond)
				}
			}(&user.Media[idx])

		}
	}

	wg.Wait()
}

func generateProfiles(numUsers int) (users []IGUser) {
	for i := 0; i < numUsers; i++ {
		medias := []IGMedia{}
		for j := 0; j < 5; j++ {
			mediaId := uuid.New()
			medias = append(medias, IGMedia{
				Id:               mediaId,
				MediaType:        "IMAGE",
				MediaUrl:         "HOST",
				Caption:          "Hello World",
				IsCommentEnabled: true,
				Timestamp:        time.Now(),
				Insights: Insights{
					Id:          mediaId,
					Comments:    rand.Intn(150),
					Likes:       rand.Intn(10000),
					Engagement:  rand.Intn(10000),
					Impressions: rand.Intn(10000),
					Reach:       rand.Intn(10000),
					Saved:       rand.Intn(10000),
				},
			})
		}

		users = append(users, IGUser{
			Id:             uuid.New(),
			Name:           "Placeholder",
			FollowersCount: rand.Intn(100000),
			FollowsCount:   rand.Intn(2500),
			Media:          medias,
		})
	}
	return
}

func publishProfile(conn *kafka.Conn, user IGUser) {
	log.Println("publishing user:", user.Id)

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(user)

	_, err := conn.WriteMessages(
		kafka.Message{Value: reqBodyBytes.Bytes()},
	)

	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}

func publishMedia(conn *kafka.Conn, media IGMedia) {
	log.Println("publishing media:", media.Id)

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(media.Insights)

	_, err := conn.WriteMessages(
		kafka.Message{Value: reqBodyBytes.Bytes()},
	)

	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}

func updateMedia(media *IGMedia) {
	media.Insights.Comments = media.Insights.Comments + rand.Intn(2)
	media.Insights.Likes = media.Insights.Likes + rand.Intn(10)
	media.Insights.Engagement = media.Insights.Engagement + rand.Intn(10)
	media.Insights.Impressions = media.Insights.Impressions + rand.Intn(20)
	media.Insights.Reach = media.Insights.Reach + rand.Intn(20)
	media.Insights.Saved = media.Insights.Saved + rand.Intn(5)
}

func updateProfile(user *IGUser) {
	log.Println("updating user:", user.Id)
	user.FollowersCount = user.FollowersCount + rand.Intn(100)
	user.FollowsCount = user.FollowersCount + rand.Intn(2)
}

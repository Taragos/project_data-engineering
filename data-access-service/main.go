package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"

	ch "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func main() {
	db := setupDatabase()
	clickhouseClient := setupClickhouse()

	app := fiber.New()

	mediaStatement, err := db.Prepare(`SELECT ig_media.id,media_type,caption,is_comment_enabled, ig_user.name FROM ig_media join ig_user on profile_id = ig_user.id;`)
	if err != nil {
		log.Fatal("failed to prepare statment: ", err)
	}

	app.Get("/api/all", func(c *fiber.Ctx) error {
		rows, err := mediaStatement.Query()
		if err != nil {
			log.Fatal("could not retrieve media for user: ", err)
		}
		defer rows.Close()

		userMediaMap := make(map[string][]IGMedia)

		for rows.Next() {
			var media IGMedia
			var username string
			err := rows.Scan(&media.Id, &media.MediaType, &media.Caption, &media.IsCommentEnabled, &username)
			if err != nil {
				log.Fatal(err)
			}

			userMedia, ok := userMediaMap[username]
			if ok {
				userMedia = append(userMedia, media)
			} else {
				userMedia = []IGMedia{media}
			}
			userMediaMap[username] = userMedia
		}

		response := make(map[string][]FullResult)
		for username, medias := range userMediaMap {
			userResults := []FullResult{}
			for _, media := range medias {
				currentResult := FullResult{}
				currentResult.Media = media

				insights := []Insight{}
				rows, err := clickhouseClient.Query(context.Background(),
					fmt.Sprintf("SELECT round(avg(comments)),round(avg(engagement)),round(avg(impressions)),round(avg(likes)),round(avg(reach)),round(avg(saved)),toStartOfMinute((fromUnixTimestamp(_timestamp))) as minute FROM ig_insights WHERE id='%s' GROUP BY minute ORDER BY minute", media.Id))
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()

				for rows.Next() {
					var comments float64
					var engagement float64
					var impressions float64
					var likes float64
					var reach float64
					var saved float64
					var timestamp time.Time

					err := rows.Scan(&comments, &engagement, &impressions, &likes, &reach, &saved, &timestamp)
					if err != nil {
						log.Fatal(err)
					}

					insights = append(insights, Insight{Group: "Comments", Timestamp: timestamp, Value: int(comments)})
					insights = append(insights, Insight{Group: "Engagement", Timestamp: timestamp, Value: int(engagement)})
					insights = append(insights, Insight{Group: "Impressions", Timestamp: timestamp, Value: int(impressions)})
					insights = append(insights, Insight{Group: "Likes", Timestamp: timestamp, Value: int(likes)})
					insights = append(insights, Insight{Group: "Reach", Timestamp: timestamp, Value: int(reach)})
					insights = append(insights, Insight{Group: "Saved", Timestamp: timestamp, Value: int(saved)})
				}
				currentResult.Insights = insights

				userResults = append(userResults, currentResult)
			}

			response[username] = userResults

		}
		return c.JSON(response)
	})

	log.Fatal(app.Listen(":" + loadEnvOrCrash("DATA_ACCESS_SERVICE_PORT")))
}

func loadEnvOrCrash(env string) string {
	result, exists := os.LookupEnv(env)

	if !exists {
		log.Fatal("env variable not set: ", env)
	}

	return result
}

func setupDatabase() *sql.DB {
	pgHost := loadEnvOrCrash("POSTGRES_HOST")
	pgPort := loadEnvOrCrash("POSTGRES_PORT")
	pgPassword := loadEnvOrCrash("POSTGRES_PASSWORD")
	pgUser := loadEnvOrCrash("POSTGRES_USER")
	pgDb := loadEnvOrCrash("POSTGRES_DB")

	pgUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDb)

	db, err := sql.Open("postgres", pgUrl)

	if err != nil {
		log.Fatal(err)
	}

	return db

}

/*
Creates minio client, waits for minio to be responsive and creates required bucket if it doesn't exist already
*/
func setupClickhouse() ch.Conn {
	clickhouseEndpoint := loadEnvOrCrash("CLICKHOUSE_ENDPOINT")
	clickhouseClient, err := clickhouse.Open(&clickhouse.Options{Addr: []string{clickhouseEndpoint}})
	if err != nil {
		log.Fatal(err)
	}
	v, _ := clickhouseClient.ServerVersion()
	fmt.Println(v.String())
	return clickhouseClient
}

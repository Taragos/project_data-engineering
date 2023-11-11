package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/segmentio/kafka-go"
)

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
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err != nil {
		log.Fatal(err)
	}

	return db

}

/*
Creates minio client, waits for minio to be responsive and creates required bucket if it doesn't exist already
*/
func setupMinio() *minio.Client {

	s3Endpoint := loadEnvOrCrash("S3_ENDPOINT")
	s3AccessKeyID := loadEnvOrCrash("S3_ACCESS_KEY_ID")
	s3SecretAccessKey := loadEnvOrCrash("S3_SECRET_ACCESS_KEY")

	minioClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKeyID, s3SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("could not open connection to minio: ", err)
	}
	cancel, err := minioClient.HealthCheck(time.Second * 5)
	if err != nil {
		log.Fatal("healthcheck already running", err)
	}

	for minioClient.IsOffline() {
		time.Sleep(time.Second * 5)
		log.Println("Minio still offline")
	}
	cancel()

	return minioClient
}

func setupBucket(minioClient *minio.Client, bucket string) {
	exists, err := minioClient.BucketExists(context.Background(), bucket)

	if err != nil {
		log.Fatal("could not check whether bucket exists or not: ", err)
	}

	if !exists {
		minioClient.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
		policy := `{ "Version": "2012-10-17", "Statement": [ { "Effect": "Allow", "Principal": { "AWS": [ "*" ] }, "Action": [ "s3:GetBucketLocation", "s3:ListBucket" ], "Resource": [ "arn:aws:s3:::test" ] }, { "Effect": "Allow", "Principal": { "AWS": [ "*" ] }, "Action": [ "s3:GetObject" ], "Resource": [ "arn:aws:s3:::test/*" ] } ] }`
		minioClient.SetBucketPolicy(context.Background(), bucket, policy)
	}
}

func setupKafka(id int) *kafka.Reader {
	kafkaBootstrapServers := loadEnvOrCrash("KAFKA_BOOTSTRAP_SERVERS")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaBootstrapServers},
		Topic:     "instagram-profiles",
		Partition: 0,
		MaxBytes:  10e6,
	})
}

func main() {
	// 1. Establish Connection to PostgreSQL + Clickhouse

	s3Bucket := loadEnvOrCrash("S3_BUCKET")

	db := setupDatabase()
	minioClient := setupMinio()
	setupBucket(minioClient, s3Bucket)

	var wg sync.WaitGroup

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go worker(db, minioClient, s3Bucket, i)
	}

	wg.Wait()
}

/*
Main worker process that subscribes to Kafka and handles new IG objects
*/
func worker(db *sql.DB, minioClient *minio.Client, bucket string, id int) {
	// 2	. Establish Connection to Kafka
	kafka := setupKafka(id)
	defer kafka.Close()

	userCache := make(map[uuid.UUID][]uuid.UUID)

	for {
		// 3. Consume Message
		m, err := kafka.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("could not read message: ", err)
		}

		user := IGUser{}
		json.Unmarshal(m.Value, &user)
		userMedias, ok := userCache[user.Id]
		if !ok {
			insertUserIntoDB(user, db)
			userMedias = []uuid.UUID{}
		}

		for idx := range user.Media {
			media := user.Media[idx]
			if slices.Contains(userMedias, media.Id) {
				continue
			}
			insertMediaIntoDB(media, user.Id, db)
			uploadImageToS3(media, minioClient, bucket)
			userMedias = append(userMedias, media.Id)
		}

		userCache[user.Id] = userMedias
	}
}

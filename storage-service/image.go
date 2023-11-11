package main

import (
	"context"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
)

/*
Downloads the image referenced in media.MediaUrl and re-uploads it to a S3 bucket
*/
func uploadImageToS3(media IGMedia, minioClient *minio.Client, bucket string) {
	resp, err := http.Get("http://" + media.MediaUrl)
	if err != nil {
		log.Fatalf("failed to download image %s because of: %v", media.MediaUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("bad status: %s", resp.Status)
	}

	log.Println("downloaded media using url: ", media.MediaUrl)
	log.Println(resp)
	minioClient.PutObject(
		context.Background(),
		bucket,
		media.Id.String(),
		resp.Body,
		resp.ContentLength,
		minio.PutObjectOptions{
			ContentType: resp.Header.Get("Content-type"),
		},
	)
}

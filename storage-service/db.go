package main

import (
	"database/sql"
	"log"
	"regexp"
	"strings"

	"github.com/gofrs/uuid"
)

/*
Inserts an IGUser object into the database table ig_user
Does nothing if the user alredy exists
*/
func insertUserIntoDB(user IGUser, db *sql.DB) {
	// If User does not exist -> Create
	log.Printf("inserting user: %s", user.Name)
	_, err := db.Exec(`INSERT INTO "ig_user"(id, name) values($1, $2) ON CONFLICT DO NOTHING`, user.Id, user.Name)
	if err != nil {
		log.Fatal("could not insert/update user: ", err)
	}
}

/*
Inserts an IGMedia object into the database table ig_media
Does nothing if the user already exists
*/
func insertMediaIntoDB(media IGMedia, userId uuid.UUID, db *sql.DB) {
	// Check if Media exists
	var count int
	_ = db.QueryRow(`SELECT COUNT(*) FROM ig_media WHERE id=$1`, media.Id).Scan(&count)

	if count == 0 {
		// Create Media
		_, err := db.Exec(`INSERT INTO "ig_media"(id, profile_id, media_type, media_url, caption, is_comment_enabled, ig_created_at) values($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING`,
			media.Id,
			userId,
			media.MediaType,
			media.MediaUrl,
			media.Caption,
			media.IsCommentEnabled,
			media.Timestamp,
		)
		if err != nil {
			log.Println("insert error: ", err)
		}

		tags := extractTags(media.Caption)

		log.Println("found tags: ", tags)
		for _, tag := range tags {
			var id uuid.UUID
			err = db.QueryRow("SELECT id FROM ig_tag WHERE tag='$1'", tag).Scan(&id)

			if err == sql.ErrNoRows {
				db.QueryRow(`INSERT INTO ig_tag(tag) VALUES($1) RETURNING id`, tag).Scan(&id)
			}

			db.Exec(`INSERT INTO ig_media_tag(ig_media_id, ig_tag_id) VALUES($1, $2)`, media.Id, id)
		}

	}
}

func extractTags(caption string) (tags []string) {
	re := regexp.MustCompile(`(#[a-zA-Z]+)`)
	for _, match := range re.FindAllString(caption, -1) {
		tags = append(tags, strings.ReplaceAll(match, "#", ""))
	}
	return tags
}

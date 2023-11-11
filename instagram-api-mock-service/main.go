package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

var captions = []string{
	"Chasing dreams and catching sunsets. #DreamChaser #SunsetMagic 🌅✨",
	"Life is short, smile while you still have teeth. #SmileMore #EnjoyLife 😁💖",
	"Coffee and kindness. #CoffeeLover #KindnessMatters ☕️💕",
	"Creating my own sunshine. #PositiveVibes #SunshineState 🌞💛",
	"Adventure awaits. Go find it. #AdventureTime #ExploreMore 🚀🌍",
	"Slaying all day. #SlayQueen #Confidence 👑💄",
	"Good vibes only. #GoodVibes #Positivity ✌️😊",
	"Do more things that make you forget to check your phone. #Disconnect #LiveInTheMoment 📵🌟",
	"Be a voice, not an echo. #BeUnique #SpeakUp 🗣️🔊",
	"Making memories in my favorite place. #MakingMemories #FavoriteSpot 🗺️📸",
	"Just because you're awake doesn't mean you should stop dreaming. #DreamBig #StayInspired 🌌💭",
	"Sunshine mixed with a little hurricane. #SunshineAndStorm #Balance 🌦️⚖️",
	"Do it with passion or not at all. #PassionFirst #LoveWhatYouDo ❤️‍🔥🚀",
	"Living my story and loving it. #LivingMyStory #LoveLife 📖💖",
	"Radiate positive vibes. #PositiveEnergy #GoodVibesOnly ✨😇",
	"Embracing the journey. #JourneyOfLife #AdventureAwaits 🌄🛤️",
	"Capturing moments that turn into memories. #CaptureTheMoment #MemoriesMade 📷🎉",
	"In the pursuit of happiness. #PursuitOfHappiness #ChooseJoy 😊💪",
	"Smiling my way through life's adventures. #SmileAlways #AdventureTime 😄🌟",
	"Elevating my vibes, one post at a time. #ElevateYourVibe #PositivityOnPoint 🚀🔝",
	"Wander often, wonder always. #Wanderlust #StayCurious 🌍❓",
	"Blessed with a resting beach face. #BeachVibes #RestingBeachFace 🏖️😎",
	"Chasing sunsets and dreams. #SunsetChaser #Dreamer 🌇💫",
	"Living my fairytale. #FairytaleLife #HappilyEverAfter 👑🏰",
	"Sunshine on my mind. #SunshineThoughts #PositiveMindset ☀️🤔",
	"Good times + Crazy friends = Amazing memories. #FriendshipGoals #MemoriesForever 🤪👫",
	"Born to stand out. #StandOut #BeBold 🌟🎨",
	"Life is better in flip-flops. #BeachLife #FlipFlopStateOfMind 🏝️👣",
	"Sassy, classy with a touch of bad-assy. #SassyStyle #ClassyAndSassy 💁‍♀️💋",
	"Finding paradise wherever I go. #ParadiseFound #Wanderlust 🌺🌴",
	"Dream big, sparkle more, shine bright. #DreamBig #SparkleAndShine ✨💖",
	"Adventure is calling, and I must go. #AdventureAwaits #AnswerTheCall 🌲🏞️",
	"Not all who wander are lost. #Wanderlust #FindYourPath 🚶‍♂️🌌",
	"Living my life in a flip-flop state of mind. #FlipFlopLife #BeachVibes 🏖️🌊",
	"Chin up, princess. Or the crown slips. #ChinUpPrincess #CrownOnPoint 👑😊",
	"Happiness looks gorgeous on me. #HappinessIsKey #Glowing 💫😁",
	"Life's a journey, not a destination. #JourneyOfLife #EnjoyTheRide 🚗🌟",
	"Sweeter than honey. #SweetLife #HoneyVibes 🍯🌼",
	"Living my happily ever after. #HappilyEverAfter #FairyTaleEnding 👑💖",
	"Be a voice, not an echo. #BeUnique #SpeakYourMind 🗣️💬",
	"Adventure is out there. #AdventureAwaits #ExploreMore 🌍🚀",
	"Grateful for this beautiful life. #Gratitude #BeautifulLife 🙏🌈",
	"Confidence level: Selfie with no filter. #ConfidenceGoals #NoFilterNeeded 😎📸",
	"Sunshine mixed with a little hurricane. #SunshineAndStorm #Balance 🌦️⚖️",
	"Let the sea set you free. #SeaLove #OceanEscape 🌊💙",
	"Living my own fairy tale. #FairyTaleLife #LivingTheDream 👸✨",
	"Elegance is an attitude. #Elegance #AttitudeOnPoint 💃🔥",
	"Positive mind, positive vibes, positive life. #PositiveMindset #GoodVibesOnly 🌈😊",
	"Wander often, wonder always. #Wanderlust #StayCurious 🌍❓",
	"Blessed with a resting beach face. #BeachVibes #RestingBeachFace 🏖️😎",
	"Chasing sunsets and dreams. #SunsetChaser #Dreamer 🌇💫",
	"Living my fairytale. #FairytaleLife #HappilyEverAfter 👑🏰",
	"Sunshine on my mind. #SunshineThoughts #PositiveMindset ☀️🤔",
	"Good times + Crazy friends = Amazing memories. #FriendshipGoals #MemoriesForever 🤪👫",
	"Born to stand out. #StandOut #BeBold 🌟🎨",
	"Life is better in flip-flops. #BeachLife #FlipFlopStateOfMind 🏝️👣",
	"Sassy, classy with a touch of bad-assy. #SassyStyle #ClassyAndSassy 💁‍♀️💋",
	"Finding paradise wherever I go. #ParadiseFound #Wanderlust 🌺🌴",
	"Dream big, sparkle more, shine bright. #DreamBig #SparkleAndShine ✨💖",
	"Adventure is calling, and I must go. #AdventureAwaits #AnswerTheCall 🌲🏞️",
}

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
	kafkaBootstrapServers := loadEnvOrCrash("KAFKA_BOOTSTRAP_SERVERS")
	numProfiles := loadEnvOrCrashInt("NUM_PROFILES")
	numPicturesPerUser := loadEnvOrCrashInt("NUM_PICTURES_PER_PROFILE")
	insightUpdateFreq := loadEnvOrCrashInt("INSIGHT_UPDATE_FREQ_MS")
	profileUpdateFreq := loadEnvOrCrashInt("PROFILE_UPDATE_FREQ_MS")

	flag.Parse()

	log.Println("connection to kafka: ", kafkaBootstrapServers)

	getPictures(numProfiles, numPicturesPerUser)
	users := generateProfiles(numProfiles, numPicturesPerUser)
	var wg sync.WaitGroup

	for _, user := range users {
		wg.Add(1)
		go func(workUser IGUser, profileUpdateFreq int) {
			conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaBootstrapServers, "instagram-profiles", 0)
			if err != nil {
				log.Fatal("failed to dial leader:", err)
			}
			defer conn.Close()

			for {
				publishProfile(conn, workUser)
				updateProfile(&workUser)
				time.Sleep(time.Duration(profileUpdateFreq) * time.Millisecond)

			}
		}(user, profileUpdateFreq)

		for idx := range user.Media {
			wg.Add(1)
			go func(workMedia IGMedia) {
				conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaBootstrapServers, "instagram-insights", 0)
				if err != nil {
					log.Fatal("failed to dial leader:", err)
				}
				defer conn.Close()

				for {
					publishMedia(conn, workMedia)
					updateMedia(&workMedia)
					time.Sleep(time.Duration(insightUpdateFreq) * time.Millisecond)
				}
			}(user.Media[idx])
		}
	}

	app := fiber.New()

	app.Static("/images", "./images")

	app.Listen(":3000")
}

func getPictures(numUsers, numPicturesPerUser int) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := "https://thispersondoesnotexist.com/"

	for i := 0; i < numUsers*numPicturesPerUser; i++ {
		response, err := http.Get(url)
		if err != nil {
			log.Fatal("failed to download image: ", err)
		}

		defer response.Body.Close()

		file, err := os.Create(fmt.Sprintf("images/%d.jpeg", i))
		if err != nil {
			log.Fatal("failed to create file for img: ", err)
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal("failed to write to image: ", err)
		}
		log.Printf("Successfully downloaded image for jpeg images/%d.jpeg\n", i)
		// Required so new pictures are generated
		time.Sleep(time.Millisecond * 500)
	}

}

func generateProfiles(numUsers, numPicturesPerUser int) (users []IGUser) {
	hostname, _ := os.Hostname()
	for i := 0; i < numUsers; i++ {
		medias := []IGMedia{}
		for j := i * numPicturesPerUser; j < (i*numPicturesPerUser)+numPicturesPerUser; j++ {
			mediaId := uuid.New()
			medias = append(medias, IGMedia{
				Id:               mediaId,
				MediaType:        "IMAGE",
				MediaUrl:         fmt.Sprintf("%s:%d/images/%d.jpeg", hostname, 3000, j),
				Caption:          captions[rand.Intn(len(captions))],
				IsCommentEnabled: true,
				Timestamp:        time.Now(),
				Insights: Insights{
					Id:          mediaId,
					Comments:    0,
					Likes:       0,
					Engagement:  0,
					Impressions: 0,
					Reach:       0,
					Saved:       0,
				},
			})
		}

		users = append(users, IGUser{
			Id:             uuid.New(),
			Name:           fmt.Sprintf("User %d", i),
			FollowersCount: rand.Intn(100000),
			FollowsCount:   rand.Intn(2500),
			Media:          medias,
		})
	}
	return
}

func loadEnvOrCrashInt(env string) int {
	result, err := strconv.Atoi(loadEnvOrCrash(env))

	if err != nil {
		log.Fatalf("failed to convert env %s to int: %s", env, err)
	}

	return result
}

func loadEnvOrCrash(env string) string {
	result, exists := os.LookupEnv(env)

	if !exists {
		log.Fatal("env variable not set: ", env)
	}

	return result
}

func publishProfile(conn *kafka.Conn, user IGUser) {
	log.Println("publishing user:", user.Name)

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(user)

	_, err := conn.WriteMessages(
		kafka.Message{Value: reqBodyBytes.Bytes()},
	)

	log.Println("published user:", user.Name)

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
	media.Insights.Comments = max(media.Insights.Comments+rand.Intn(5)-rand.Intn(3), 0)
	media.Insights.Likes = max(media.Insights.Likes+rand.Intn(10)-rand.Intn(6), 0)
	media.Insights.Engagement = max(media.Insights.Engagement+rand.Intn(10)-rand.Intn(6), 0)
	media.Insights.Impressions = max(media.Insights.Impressions+rand.Intn(20)-rand.Intn(12), 0)
	media.Insights.Reach = max(media.Insights.Reach+rand.Intn(20)-rand.Intn(12), 0)
	media.Insights.Saved = max(media.Insights.Saved+rand.Intn(5)-rand.Intn(3), 0)
}

func updateProfile(user *IGUser) {
	log.Println("updating user:", user.Id)
	user.FollowersCount = max(user.FollowersCount+rand.Intn(100)-rand.Intn(75), 0)
	user.FollowsCount = max(user.FollowersCount+rand.Intn(5)-rand.Intn(3), 0)
}

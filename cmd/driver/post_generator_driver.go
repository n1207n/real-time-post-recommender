package main

import (
	"github.com/n1207n/real-time-post-recommender/post"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"log"
	"math/rand"
	"sync"
	"time"
)

var _ = gentleman.Client{}
var _ = post.Faker{}

var (
	apiClient = gentleman.New()
	faker     = &post.Faker{}
)

const (
	MinParagraph = 5
	MaxParagraph = 10
	MinSentence  = 5
	MaxSentence  = 10
	MinWord      = 10
	MaxWord      = 20

	DataN   = 100_000
	WorkerN = 4
)

func main() {
	wg := &sync.WaitGroup{}

	// Define base URL
	apiClient.URL("http://localhost:8080")

	log.Printf("Post fake data generator started\n")

	for i := 0; i < WorkerN; i++ {
		wg.Add(1)
		go createFakePostsAndVote(wg)
	}

	wg.Wait()
}

func createFakePostsAndVote(wg *sync.WaitGroup) {
	for i := 0; i < DataN; i++ {
		data := map[string]string{
			"title": post.GeneratePostTitle(MinWord, MaxWord, faker),
			"body":  post.GeneratePostBody(MinParagraph, MaxParagraph, MinSentence, MaxSentence, MinWord, MaxWord, faker),
		}

		// Create post
		req := apiClient.Request()
		req.Path("/posts")
		req.Method("POST")
		req.Use(body.JSON(data))

		res, err := req.Send()
		if err != nil {
			log.Printf("Request error: %s\n", err)
			continue
		}
		if !res.Ok {
			log.Printf("Invalid server response: %d\n", res.StatusCode)
			continue
		}

		// Parse createPost response
		var postResponse struct {
			Id        string    `json:"id"`
			Title     string    `json:"title"`
			Body      string    `json:"body"`
			Votes     int       `json:"votes"`
			Timestamp time.Time `json:"timestamp"`
		}
		err = res.JSON(&postResponse)
		if err != nil {
			log.Printf("Failed to parse JSON response: %d\n", res.StatusCode)
			continue
		}

		// Generate random vote actions up to 100 times
		voteCount := rand.Intn(100)
		postId := postResponse.Id

		for i := 0; i < voteCount; i++ {
			req = apiClient.Request()
			req.Path("/posts/vote")
			req.Method("POST")

			var request struct {
				Id       string `json:"id"`
				IsUpvote bool   `json:"is_upvote"`
			}
			request.Id = postId

			if prob := rand.Float64(); prob <= 0.5 {
				request.IsUpvote = false
			} else {
				request.IsUpvote = true
			}

			req.Use(body.JSON(request))

			res, err = req.Send()
			if err != nil {
				log.Printf("Request error: %s\n", err)
				continue
			}
			if !res.Ok {
				log.Printf("Invalid server response: %d\n", res.StatusCode)
				continue
			}
		}

		// Query Post by ID
		req = apiClient.Request()
		req.Path("/posts/:id")
		req.Param("id", postId)
		req.Method("GET")

		res, err = req.Send()
		if err != nil {
			log.Printf("Request error: %s\n", err)
			continue
		}
		if !res.Ok {
			log.Printf("Invalid server response: %d\n", res.StatusCode)
			continue
		}

		err = res.JSON(&postResponse)
		if err != nil {
			log.Printf("Failed to parse JSON response: %d\n", res.StatusCode)
			continue
		}

		log.Printf("ID: %s - VoteCount: %d - Final Votes: %d\n", postId, voteCount, postResponse.Votes)
	}

	wg.Done()
}

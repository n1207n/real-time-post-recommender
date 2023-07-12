package main

import (
	"github.com/imroc/req/v3"
	"github.com/n1207n/real-time-post-recommender/post"
	"log"
	"math/rand"
	"sync"
	"time"
)

var _ = req.Client{}
var _ = post.Faker{}

var (
	client = req.C()
	faker  = &post.Faker{}
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
	log.Printf("Post fake data generator started\n")

	for i := 0; i < WorkerN; i++ {
		wg.Add(1)
		go createFakePostsAndVote(wg)
	}

	wg.Wait()
}

func createFakePostsAndVote(wg *sync.WaitGroup) {
	for i := 0; i < DataN; i++ {
		var createPostPayload struct {
			Title string
			Body  string
		}
		createPostPayload.Title = post.GeneratePostTitle(MinWord, MaxWord, faker)
		createPostPayload.Body = post.GeneratePostBody(MinParagraph, MaxParagraph, MinSentence, MaxSentence, MinWord, MaxWord, faker)

		// Create post
		request := client.R()
		request.SetBodyJsonMarshal(&createPostPayload)
		response, err := request.Post("http://localhost:8080/posts")

		if err != nil {
			log.Printf("Request error: %s\n", err)
			continue
		}
		if !response.IsSuccessState() {
			log.Printf("Invalid server response: %d\n", response.StatusCode)
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
		err = response.Unmarshal(&postResponse)
		if err != nil {
			log.Printf("Failed to parse JSON response: %v\n", err.Error())
			continue
		}

		// Generate random vote actions up to 100 times
		voteCount := rand.Intn(100)
		postId := postResponse.Id

		for i := 0; i < voteCount; i++ {
			// Create post
			voteRequest := client.R()

			var votePayload struct {
				Id       string `json:"id"`
				IsUpvote bool   `json:"is_upvote"`
			}
			votePayload.Id = postId

			if prob := rand.Float64(); prob <= 0.5 {
				votePayload.IsUpvote = false
			} else {
				votePayload.IsUpvote = true
			}

			voteRequest.SetBodyJsonMarshal(&votePayload)
			voteResponse, voteErr := voteRequest.Post("http://localhost:8080/posts/vote")

			if voteErr != nil {
				log.Printf("Request error: %s\n", voteErr)
				continue
			}
			if !voteResponse.IsSuccessState() {
				log.Printf("Invalid server response: %d\n", voteResponse.StatusCode)
				continue
			}
		}

		// Query Post by ID
		queryRequest := client.R()
		queryRequest.SetPathParam("postId", postId)
		queryResponse, queryErr := queryRequest.Get("http://localhost:8080/posts/{postId}")

		if queryErr != nil {
			log.Printf("Request error: %s\n", queryErr)
			continue
		}
		if !queryResponse.IsSuccessState() {
			log.Printf("Invalid server response: %d\n", queryResponse.StatusCode)
			continue
		}

		err = queryResponse.Unmarshal(&postResponse)
		if err != nil {
			log.Printf("Failed to parse JSON response: %v\n", err.Error())
			continue
		}

		log.Printf("ID: %s - VoteCount: %d - Final Votes: %d\n", postId, voteCount, postResponse.Votes)
	}

	wg.Done()
}

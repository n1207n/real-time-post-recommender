package main

import (
	"fmt"
	"github.com/n1207n/real-time-post-recommender/post"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"sync"
)

var _ gentleman.Client = gentleman.Client{}
var _ post.Faker = post.Faker{}

var (
	apiClient *gentleman.Client = gentleman.New()
	faker                       = &post.Faker{}
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

	fmt.Printf("Post fake data generator started\n")

	for _ = range [WorkerN]int{} {
		wg.Add(1)
		go createFakePosts(wg)
	}

	wg.Wait()
}

func createFakePosts(wg *sync.WaitGroup) {
	for i := 0; i < DataN; i++ {
		data := map[string]string{
			"title": post.GeneratePostTitle(MinWord, MaxWord, faker),
			"body":  post.GeneratePostBody(MinParagraph, MaxParagraph, MinSentence, MaxSentence, MinWord, MaxWord, faker),
		}

		req := apiClient.Request()
		req.Path("/posts")
		req.Method("POST")

		// Serialize
		req.Use(body.JSON(data))

		res, err := req.Send()
		if err != nil {
			fmt.Printf("Request error: %s\n", err)
			continue
		}
		if !res.Ok {
			fmt.Printf("Invalid server response: %d\n", res.StatusCode)
			continue
		}

		fmt.Printf("Status: %d - Body: %s\n\n", res.StatusCode, res.String())
	}

	wg.Done()
}

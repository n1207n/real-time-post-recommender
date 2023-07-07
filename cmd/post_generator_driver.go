package main

import (
	"fmt"
	"github.com/n1207n/real-time-post-recommender/post"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"sync"
)

var _ gentleman.Client = gentleman.Client{}

var (
	apiClient *gentleman.Client = gentleman.New()
)

const (
	DATA_N   = 100_000
	WORKER_N = 4
)

func main() {
	wg := &sync.WaitGroup{}

	// Define base URL
	apiClient.URL("http://localhost:8080")

	fmt.Printf("Post fake data generator started\n")

	for _ = range [WORKER_N]int{} {
		wg.Add(1)
		go createFakePosts(wg)
	}

	wg.Wait()
}

func createFakePosts(wg *sync.WaitGroup) {
	for i := 0; i < DATA_N; i++ {
		data := map[string]string{"title": post.GeneratePostTitle(), "body": post.GeneratePostBody()}

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

		fmt.Printf("Status: %d - Body: %s\n\n", res.String())
	}

	wg.Done()
}

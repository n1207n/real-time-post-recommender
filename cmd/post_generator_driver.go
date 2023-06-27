package main

import (
	"fmt"
	post "github.com/n1207n/real-time-post-recommender/post"
)

func main() {
	post := post.NewPost(post.GeneratePostTitle(), post.GeneratePostBody())

	fmt.Println("Post ID:", post.ID)
	fmt.Println("Title:", post.Title)
	fmt.Println("Body:", post.Body)
	fmt.Println("Votes:", post.Votes)
	fmt.Println("Timestamp:", post.Timestamp)
}

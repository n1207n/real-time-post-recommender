package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/n1207n/real-time-post-recommender/post"
	"github.com/n1207n/real-time-post-recommender/ranking"
	"log"
	"net/http"
)

type postCreateRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

// CreatePost creates a new Post instsance
func CreatePost(context *gin.Context) {
	var payload postCreateRequest

	if err := context.ShouldBindJSON(&payload); err != nil {
		log.Printf("CreatePost - context.ShouldBindJSON(&payload): %v", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPost := post.NewPost(payload.Title, payload.Body)
	post.PostServiceInstance.Create(*newPost)

	if err := ranking.Ranker.PushPostScore(newPost); err != nil {
		log.Printf("CreatePost - ranking.Ranker.PushPostScore(newPost): %v", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, newPost)
}

package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/n1207n/real-time-post-recommender/post"
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
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	newPost := post.NewPost(payload.Title, payload.Body)
	post.PostServiceInstance.CreatePost(newPost)

	context.JSON(http.StatusOK, newPost)
}

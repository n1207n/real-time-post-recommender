package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/n1207n/real-time-post-recommender/post"
	"net/http"
	"strconv"
)

const (
	DEFAULT_LIMIT  = "20"
	DEFAULT_OFFSET = "0"
)

// ListPosts returns a paginated list of Post instances
func ListPosts(context *gin.Context) {
	limit, err := strconv.Atoi(context.DefaultQuery("limit", DEFAULT_LIMIT))
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(context.DefaultQuery("offset", DEFAULT_OFFSET))
	if err != nil {
		offset = 0
	}

	posts := post.PostServiceInstance.List(limit, offset)
	context.JSON(http.StatusOK, gin.H{
		"data":   posts,
		"count":  len(posts),
		"limit":  limit,
		"offset": offset,
	})
}

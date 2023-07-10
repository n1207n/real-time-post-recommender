package routing

import (
	"github.com/gin-gonic/gin"
)

// BuildRouters registers the API endpoints
func BuildRouters(r *gin.Engine) {
	r.GET("/", Index)
	r.GET("/posts", ListPosts)
	r.GET("/posts/:id", GetPost)
	r.GET("/posts/top", ListTopRankedPosts)
	r.POST("/posts", CreatePost)
	r.POST("/posts/vote", VotePost)
}

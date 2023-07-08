package routing

import (
	"github.com/gin-gonic/gin"
)

// BuildRouters registers the API endpoints
func BuildRouters(r *gin.Engine) {
	r.GET("/", Index)
	r.GET("/posts", ListPosts)
	r.POST("/posts", CreatePost)
	r.POST("/posts/vote", VotePost)
}

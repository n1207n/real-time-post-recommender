package routing

import (
	"github.com/gin-gonic/gin"
	routing "github.com/n1207n/real-time-post-recommender/routing/posts"
)

// BuildRouters registers the API endpoints
func BuildRouters(r *gin.Engine) {
	r.GET("/", Index)
	r.POST("/posts", routing.CreatePost)
}

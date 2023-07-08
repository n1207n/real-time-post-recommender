package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/n1207n/real-time-post-recommender/post"
	"net/http"
)

// GetPost returns a single Post instance
func GetPost(context *gin.Context) {
	idStr := context.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{})
		return
	}

	p, err2 := post.PostServiceInstance.GetByID(id)
	if err2 != nil {
		context.JSON(http.StatusNotFound, gin.H{})
		return
	}

	context.JSON(http.StatusOK, p)
}

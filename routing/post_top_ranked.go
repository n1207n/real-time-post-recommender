package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/n1207n/real-time-post-recommender/ranking"
	"net/http"
	"time"
)

const (
	DEFAULT_TOP_N = 50
)

// ListTopRankedPosts returns a list of Post instances ordered by ranking score up to DEFAULT_TOP_N
func ListTopRankedPosts(context *gin.Context) {
	dateForRanking := context.DefaultQuery("date", time.Now().Format("2006-01-02"))
	datetime, err := time.Parse("2006-01-02", dateForRanking)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	posts := ranking.Ranker.GetTopRankedPosts(datetime, DEFAULT_TOP_N)
	context.JSON(http.StatusOK, gin.H{
		"data":  posts,
		"count": len(posts),
	})
}

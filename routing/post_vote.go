package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/n1207n/real-time-post-recommender/post"
	"github.com/n1207n/real-time-post-recommender/ranking"
	"net/http"
)

type postVoteRequest struct {
	ID       uuid.UUID `json:"id" binding:"required"`
	IsUpvote *bool     `json:"is_upvote" binding:"required" `
}

// VotePost increments or decrements the Post instance's votes
func VotePost(context *gin.Context) {
	var payload postVoteRequest

	if err := context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := post.PostServiceInstance.Vote(payload.ID, *payload.IsUpvote)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ranking.Ranker.PushPostScore(p); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var voteAction string

	if *payload.IsUpvote == true {
		voteAction = "Upvoted"
	} else {
		voteAction = "Downvoted"
	}

	context.JSON(http.StatusOK, gin.H{
		"status": voteAction,
	})
}

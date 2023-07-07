package post

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Votes     int       `json:"votes"`
	Timestamp time.Time `json:"timestamp"`
}

func NewPost(title string, body string) *Post {
	return &Post{
		ID:        uuid.New(),
		Title:     title,
		Body:      body,
		Votes:     0,
		Timestamp: time.Now(),
	}
}

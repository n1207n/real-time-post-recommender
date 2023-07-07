package post

import (
	"github.com/google/uuid"
	"github.com/n1207n/real-time-post-recommender/persistance"
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

func (p *Post) Save() {
	stmt := persistance.SqlServiceInstance.DbClient.Rebind(
		"INSERT INTO posts (id, title, body, votes, timestamp) VALUES (:id, :title, :body, :votes, :timestamp)",
	)
	_, err := persistance.SqlServiceInstance.DbClient.NamedExec(stmt, p)
	if err != nil {
		panic(err)
	}
}

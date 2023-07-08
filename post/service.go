package post

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/n1207n/real-time-post-recommender/sql"
)

type PostService struct {
	db *sql.SqlService
}

// Compliation check
var _ PostService = PostService{}

var (
	PostServiceInstance *PostService
)

func NewPostService() *PostService {
	PostServiceInstance = &PostService{db: sql.DB}
	return PostServiceInstance
}

func (s *PostService) Create(newPost *Post) uuid.UUID {
	var lastInsertedUUID uuid.UUID

	err := s.db.Client.QueryRow(
		"INSERT INTO posts (id, title, body, votes, timestamp) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		newPost.ID, newPost.Title, newPost.Body, newPost.Votes, newPost.Timestamp,
	).Scan(&lastInsertedUUID)
	if err != nil {
		panic(err)
	}

	return lastInsertedUUID
}

func (s *PostService) List(limit int, offset int) []Post {
	var results []Post

	stmt := fmt.Sprintf("SELECT * FROM posts LIMIT %d OFFSET %d", limit, offset)
	err := sql.DB.Client.Select(&results, stmt)
	if err != nil {
		panic(err)
	}

	return results
}

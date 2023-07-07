package post

import (
	"fmt"
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

func (s *PostService) Create(newPost *Post) {
	stmt := sql.DB.Client.Rebind(
		"INSERT INTO posts (id, title, body, votes, timestamp) VALUES (:id, :title, :body, :votes, :timestamp)",
	)
	_, err := s.db.Client.NamedExec(stmt, newPost)
	if err != nil {
		panic(err)
	}
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

package post

import (
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
	results := make([]Post, limit)

	err := sql.DB.Client.Select(&results, "SELECT * FROM posts LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		panic(err)
	}

	return results
}

func (s *PostService) Vote(postId uuid.UUID, isUpvote bool) error {
	var post Post

	err := sql.DB.Client.Get(&post, "SELECT * FROM posts WHERE id = $1", postId.String())
	if err != nil {
		return err
	}

	if isUpvote {
		post.Votes += 1
	} else {
		post.Votes -= 1
	}

	_, err = sql.DB.Client.Exec(`UPDATE posts SET votes = $1 WHERE id = $2`, post.Votes, postId.String())
	return err
}

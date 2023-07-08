package post

import (
	"github.com/google/uuid"
	"github.com/n1207n/real-time-post-recommender/sql"
	"github.com/n1207n/real-time-post-recommender/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	dbHost      string
	dbPort      int
	dbUsername  string
	dbPassword  string
	dbName      string
	postService *PostService = NewPostService()
)

func setUp() {
	dbHost, dbPort, dbUsername, dbPassword, dbName = utils.LoadDBEnvVariables()

	sql.NewSqlService(dbUsername, dbPassword, dbHost, dbPort, dbName)
}

func TestPostService_Create(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")
	dbId := postService.Create(newPost)

	// Verify that the post is created successfully
	var post Post
	err := sql.DB.Client.Select(&post, "SELECT * FROM posts WHERE id = ?", dbId)
	assert.NoError(t, err)

	// Perform assertions
	assert.NotEqual(t, uuid.Nil, newPost.ID)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestPostService_List(t *testing.T) {
	setUp()

	const dataN = 10
	for range [dataN]int{} {
		newPost := NewPost("Test Title", "Test Body")
		postService.Create(newPost)
	}

	posts := postService.List(dataN, 0)

	assert.NotNil(t, posts)
	assert.Len(t, posts, dataN)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestPostService_UpvotePost(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")
	dbId := postService.Create(newPost)

	var post Post
	err := sql.DB.Client.Select(&post, "SELECT * FROM posts WHERE id = ?", dbId)

	// Upvote the post
	err = postService.Vote(post.ID, true)
	assert.NoError(t, err)

	// Retrieve the post from the database
	var post Post
	err = sql.DB.Client.Select(&post, "SELECT * FROM posts WHERE id = ?", dbId)

	// Assert that the vote count has increased by 1
	assert.Equal(t, 1, post.Votes)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestPostService_DownvotePost(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")
	dbId := postService.Create(newPost)

	var post Post
	err := sql.DB.Client.Select(&post, "SELECT * FROM posts WHERE id = ?", dbId)

	// Upvote the post
	err = postService.Vote(post.ID, false)
	assert.NoError(t, err)

	// Retrieve the post from the database
	var post Post
	err = sql.DB.Client.Select(&post, "SELECT * FROM posts WHERE id = ?", dbId)

	// Assert that the vote count has increased by 1
	assert.Equal(t, -1, post.Votes)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestPostService_Vote_NonExistingPost(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")
	newPost.ID = uuid.New()

	// Try to upvote the non-existing post
	err := postService.Vote(newPost.ID, true)
	assert.Error(t, err)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestPostService_Vote_MultipleVotes(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")
	dbId := postService.Create(newPost)

	var post Post
	err := sql.DB.Client.Select(&post, "SELECT * FROM posts WHERE id = ?", dbId)

	// Upvote the post multiple times
	err = postService.Vote(post.ID, true)
	assert.NoError(t, err)
	err = postService.Vote(post.ID, true)
	assert.NoError(t, err)

	// Downvote the post multiple times
	err = postService.Vote(post.ID, false)
	assert.NoError(t, err)
	err = postService.Vote(post.ID, false)
	assert.NoError(t, err)

	// Retrieve the post from the database
	var post Post
	err = sql.DB.Client.Select(&post, "SELECT * FROM posts WHERE id = ?", dbId)

	// Assert that the vote count reflects the cumulative effect of the multiple votes
	assert.Equal(t, 0, post.Votes)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

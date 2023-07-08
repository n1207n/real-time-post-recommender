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
	postService *PostService
)

func setUp() {
	dbHost, dbPort, dbUsername, dbPassword, dbName = utils.LoadDBEnvVariables()

	sql.NewSqlService(dbUsername, dbPassword, dbHost, dbPort, dbName)
	postService = NewPostService()
}

func TestPostService_Create(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")
	dbId := postService.Create(*newPost)

	// Verify that the post is created successfully
	var post Post
	err := sql.DB.Client.Get(&post, "SELECT * FROM posts WHERE id = $1", dbId.String())
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
		postService.Create(*newPost)
	}

	posts := postService.List(dataN, 0)

	assert.NotNil(t, posts)
	assert.Len(t, posts, dataN)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestPostService_Get(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")
	postService.Create(*newPost)

	post, err := postService.GetByID(newPost.ID)
	assert.NoError(t, err)

	assert.NotNil(t, post)
	assert.Equal(t, post.Title, newPost.Title)
	assert.Equal(t, post.Body, newPost.Body)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestPostService_UpvotePost(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")
	dbId := postService.Create(*newPost)

	post, err := postService.GetByID(dbId)
	assert.NoError(t, err)

	// Upvote the post
	post, err = postService.Vote(post.ID, true)
	assert.NoError(t, err)

	// Assert that the vote count has increased by 1
	assert.Equal(t, 1, post.Votes)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestPostService_DownvotePost(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")
	dbId := postService.Create(*newPost)

	post, err := postService.GetByID(dbId)
	assert.NoError(t, err)

	// Upvote the post
	post, err = postService.Vote(post.ID, false)
	assert.NoError(t, err)

	// Assert that the vote count has decreased by 1
	assert.Equal(t, -1, post.Votes)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestPostService_Vote_NonExistingPost(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")
	newPost.ID = uuid.New()

	// Try to upvote the non-existing post
	_, err := postService.GetByID(newPost.ID)
	assert.Error(t, err)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestPostService_Vote_MultipleVotes(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")
	dbId := postService.Create(*newPost)

	post, err := postService.GetByID(dbId)
	assert.NoError(t, err)

	// Upvote the post multiple times
	post, err = postService.Vote(post.ID, true)
	assert.NoError(t, err)
	post, err = postService.Vote(post.ID, true)
	assert.NoError(t, err)

	// Downvote the post multiple times
	post, err = postService.Vote(post.ID, false)
	assert.NoError(t, err)
	post, err = postService.Vote(post.ID, false)
	assert.NoError(t, err)

	// Assert that the vote count reflects the cumulative effect of the multiple votes
	assert.Equal(t, 0, post.Votes)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

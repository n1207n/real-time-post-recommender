package post

import (
	"github.com/google/uuid"
	"github.com/n1207n/real-time-post-recommender/sql"
	"github.com/n1207n/real-time-post-recommender/utils"
	"github.com/stretchr/testify/assert"
	"sync"
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

func setUp() func() {
	once := sync.Once{}
	once.Do(func() {
		dbHost, dbPort, dbUsername, dbPassword, dbName = utils.LoadDBEnvVariables()

		sql.NewSqlService(dbUsername, dbPassword, dbHost, dbPort, dbName)
		postService = NewPostService()
	})

	// Teardown function as return value
	return func() {
		// sql.DB.Client.MustExec("TRUNCATE TABLE posts")
	}
}

func TestPostService_Create(t *testing.T) {
	teardown := setUp()
	defer teardown()

	newPost := NewPost("Test Title", "Test Body")
	dbId := postService.Create(*newPost)

	// Verify that the post is created successfully
	p, err := postService.GetByID(dbId)
	assert.NoError(t, err)

	// Perform assertions
	assert.NotEqual(t, uuid.Nil, p.ID)
	assert.Equal(t, newPost.ID, p.ID)
}

func TestPostService_List(t *testing.T) {
	teardown := setUp()
	defer teardown()

	const dataN = 10
	for i := 0; i < dataN; i++ {
		newPost := NewPost("Test Title", "Test Body")
		postService.Create(*newPost)
	}

	posts := postService.List(dataN, 0)

	assert.NotNil(t, posts)
	assert.Len(t, posts, dataN)
}

func TestPostService_FilterByIDs(t *testing.T) {
	teardown := setUp()
	defer teardown()

	const dataN = 20
	var targetIds []string

	for i := 0; i < dataN; i++ {
		newPost := NewPost("Test Title", "Test Body")
		postService.Create(*newPost)

		if i < dataN/2 {
			targetIds = append(targetIds, newPost.ID.String())
		}
	}

	posts := postService.FilterByIDs(targetIds)

	assert.NotNil(t, posts)
	assert.Len(t, posts, dataN/2)
}

func TestPostService_Get(t *testing.T) {
	teardown := setUp()
	defer teardown()

	newPost := NewPost("Test Title", "Test Body")
	postService.Create(*newPost)

	post, err := postService.GetByID(newPost.ID)
	assert.NoError(t, err)

	assert.NotNil(t, post)
	assert.Equal(t, post.Title, newPost.Title)
	assert.Equal(t, post.Body, newPost.Body)
}

func TestPostService_UpvotePost(t *testing.T) {
	teardown := setUp()
	defer teardown()

	newPost := NewPost("Test Title", "Test Body")
	dbId := postService.Create(*newPost)

	post, err := postService.GetByID(dbId)
	assert.NoError(t, err)

	// Upvote the post
	post, err = postService.Vote(post.ID, true)
	assert.NoError(t, err)

	// Assert that the vote count has increased by 1
	assert.Equal(t, 1, post.Votes)
}

func TestPostService_DownvotePost(t *testing.T) {
	teardown := setUp()
	defer teardown()

	newPost := NewPost("Test Title", "Test Body")
	dbId := postService.Create(*newPost)

	post, err := postService.GetByID(dbId)
	assert.NoError(t, err)

	// Upvote the post
	post, err = postService.Vote(post.ID, false)
	assert.NoError(t, err)

	// Assert that the vote count has decreased by 1
	assert.Equal(t, -1, post.Votes)
}

func TestPostService_Vote_NonExistingPost(t *testing.T) {
	teardown := setUp()
	defer teardown()

	newPost := NewPost("Test Title", "Test Body")
	newPost.ID = uuid.New()

	// Try to upvote the non-existing post
	_, err := postService.GetByID(newPost.ID)
	assert.Error(t, err)
}

func TestPostService_Vote_MultipleVotes(t *testing.T) {
	teardown := setUp()
	defer teardown()

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
}

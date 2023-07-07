package post

import (
	"github.com/google/uuid"
	"github.com/n1207n/real-time-post-recommender/sql"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strconv"
	"testing"
)

var (
	ss         *sql.SqlService
	dbHost     string
	dbPort     int
	dbUsername string
	dbPassword string
	dbName     string
)

func loadDBEnvVariables() (string, int, string, string, string) {
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	portStr := os.Getenv("DB_PORT")
	name := "test_db"

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Failed to parse DB_PORT: %v", err)
	}

	return host, port, username, password, name
}

func setUp() {
	dbHost, dbPort, dbUsername, dbPassword, dbName = loadDBEnvVariables()
	ss = sql.NewSqlService(dbUsername, dbPassword, dbHost, dbPort, dbName)
}

func TestPostService_Create(t *testing.T) {
	setUp()

	newPost := NewPost("Test Title", "Test Body")

	postService := NewPostService()
	postService.Create(newPost)

	// Verify that the post is created successfully
	// ... (e.g., query the database to check if the post exists)

	// Perform assertions
	assert.NotEqual(t, uuid.Nil, newPost.ID)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestPostService_List(t *testing.T) {
	setUp()

	postService := NewPostService()

	const dataN = 10
	for _ = range [dataN]int{} {
		newPost := NewPost("Test Title", "Test Body")
		postService.Create(newPost)
	}

	posts := postService.List(dataN, 0)

	assert.NotNil(t, posts)
	assert.Len(t, posts, dataN)

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

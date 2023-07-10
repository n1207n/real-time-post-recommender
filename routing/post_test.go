package routing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/n1207n/real-time-post-recommender/cache"
	"github.com/n1207n/real-time-post-recommender/post"
	"github.com/n1207n/real-time-post-recommender/ranking"
	"github.com/n1207n/real-time-post-recommender/sql"
	"github.com/n1207n/real-time-post-recommender/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	r         *gin.Engine
	postCount int
)

func setUp() func() {
	once := sync.Once{}
	once.Do(func() {
		dbHost, dbPort, dbUsername, dbPassword, dbName := utils.LoadDBEnvVariables()

		sql.NewSqlService(dbUsername, dbPassword, dbHost, dbPort, dbName)
		post.NewPostService()

		redisHost, redisPort, redisName := utils.LoadRedisEnvVariables()
		cache.NewCacheService(redisHost, redisPort, redisName)
		ranking.NewRanker()

		r = gin.Default()
		BuildRouters(r)
	})

	// Teardown function as return value
	return func() {
		// sql.DB.Client.MustExec("TRUNCATE TABLE posts")
	}
}

func TestCreatePost(t *testing.T) {
	teardown := setUp()
	defer teardown()

	data := map[string]string{
		"title": "Test title",
		"body":  "Test body",
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBytes))

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, data["title"], response["title"])
	assert.Equal(t, data["body"], response["body"])
}

func TestListPosts(t *testing.T) {
	teardown := setUp()
	defer teardown()

	err := sql.DB.Client.QueryRow("SELECT COUNT(*) FROM posts").Scan(&postCount)
	if err != nil {
		panic(err)
	}

	const dataN = 10
	for i := 0; i < dataN; i++ {
		data := map[string]string{
			"title": fmt.Sprintf("Test title %d", i),
			"body":  fmt.Sprintf("Test body %d", i),
		}

		jsonBytes, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBytes))

		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)

	// Set correct offset
	q := req.URL.Query()
	q.Add("offset", strconv.Itoa(postCount))
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var response struct {
		Count int `json:"count"`
		Data  []struct {
			Id        string    `json:"id"`
			Title     string    `json:"title"`
			Body      string    `json:"body"`
			Votes     int       `json:"votes"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"data"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, response)
	assert.Equal(t, dataN, len(response.Data))
}

func TestGetPost(t *testing.T) {
	teardown := setUp()
	defer teardown()

	data := map[string]string{
		"title": fmt.Sprintf("Test title"),
		"body":  fmt.Sprintf("Test body"),
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBytes))

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var response struct {
		Id        string    `json:"id"`
		Title     string    `json:"title"`
		Body      string    `json:"body"`
		Votes     int       `json:"votes"`
		Timestamp time.Time `json:"timestamp"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		panic(err)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/posts/%s", response.Id), nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	fmt.Printf(w.Body.String())

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, response)
	assert.Equal(t, data["title"], response.Title)
	assert.Equal(t, data["body"], response.Body)
}

func TestVotePost(t *testing.T) {
	teardown := setUp()
	defer teardown()

	// Create a new post
	data := map[string]string{
		"title": "Test title",
		"body":  "Test body",
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBytes))

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	postCount += 1

	var response struct {
		Id        string    `json:"id"`
		Title     string    `json:"title"`
		Body      string    `json:"body"`
		Votes     int       `json:"votes"`
		Timestamp time.Time `json:"timestamp"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		panic(err)
	}

	// Upvote
	var request struct {
		Id       string `json:"id"`
		IsUpvote bool   `json:"is_upvote"`
	}

	request.Id = response.Id
	request.IsUpvote = true

	jsonBytes, err = json.Marshal(request)
	if err != nil {
		panic(err)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/posts/vote", bytes.NewBuffer(jsonBytes))

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var voteResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &voteResponse)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "Upvoted", voteResponse["status"])

	// Ranking check
	ctx := cache.Cache.Ctx
	key := fmt.Sprintf("post-scores-%s", time.Now().Format("2006-01-02"))
	result, keyErr := cache.Cache.RedisClient.Exists(ctx, key).Result()
	assert.NoError(t, keyErr)
	assert.Equal(t, result, int64(1))

	// Downvote
	request.IsUpvote = false

	jsonBytes, err = json.Marshal(request)
	if err != nil {
		panic(err)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/posts/vote", bytes.NewBuffer(jsonBytes))

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &voteResponse)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "Downvoted", voteResponse["status"])

	// Ranking check
	result, keyErr = cache.Cache.RedisClient.Exists(ctx, key).Result()
	assert.NoError(t, keyErr)
	assert.Equal(t, result, int64(1))
}

func TestGetRankedPosts(t *testing.T) {
	teardown := setUp()
	defer teardown()

	const dataN = 10
	for i := 0; i < dataN; i++ {
		// Create Post
		data := map[string]string{
			"title": fmt.Sprintf("Test title %d", i),
			"body":  fmt.Sprintf("Test body %d", i),
		}

		jsonBytes, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBytes))

		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var response struct {
			Id        string    `json:"id"`
			Title     string    `json:"title"`
			Body      string    `json:"body"`
			Votes     int       `json:"votes"`
			Timestamp time.Time `json:"timestamp"`
		}

		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			panic(err)
		}

		// Upvote
		var request struct {
			Id       string `json:"id"`
			IsUpvote bool   `json:"is_upvote"`
		}

		request.Id = response.Id
		request.IsUpvote = true

		jsonBytes, err = json.Marshal(request)
		if err != nil {
			panic(err)
		}

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/posts/vote", bytes.NewBuffer(jsonBytes))

		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var voteResponse map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &voteResponse)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, "Upvoted", voteResponse["status"])
	}

	// Top ranked
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/top", nil)
	q := req.URL.Query()

	// date GET params
	q.Add("date", time.Now().Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var response struct {
		Count int `json:"count"`
		Data  []struct {
			Id        string    `json:"id"`
			Title     string    `json:"title"`
			Body      string    `json:"body"`
			Votes     int       `json:"votes"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 50, response.Count)
}

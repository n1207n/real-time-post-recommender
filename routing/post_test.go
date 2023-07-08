package routing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/n1207n/real-time-post-recommender/post"
	"github.com/n1207n/real-time-post-recommender/sql"
	"github.com/n1207n/real-time-post-recommender/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreatePost(t *testing.T) {
	dbHost, dbPort, dbUsername, dbPassword, dbName := utils.LoadDBEnvVariables()

	sql.NewSqlService(dbUsername, dbPassword, dbHost, dbPort, dbName)
	post.NewPostService()

	r := gin.Default()
	BuildRouters(r)

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

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

func TestListPosts(t *testing.T) {
	dbHost, dbPort, dbUsername, dbPassword, dbName := utils.LoadDBEnvVariables()

	sql.NewSqlService(dbUsername, dbPassword, dbHost, dbPort, dbName)
	post.NewPostService()

	r := gin.Default()
	BuildRouters(r)

	const dataN = 10
	for i := range [dataN]int{} {
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

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	fmt.Printf(w.Body.String())

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

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, response)
	assert.Equal(t, dataN, len(response.Data))

	// Cleanup
	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
}

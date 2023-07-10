package ranking

import (
	"context"
	"fmt"
	"github.com/n1207n/real-time-post-recommender/cache"
	"github.com/n1207n/real-time-post-recommender/post"
	"github.com/n1207n/real-time-post-recommender/sql"
	"github.com/n1207n/real-time-post-recommender/utils"
	"github.com/stretchr/testify/assert"
	"math"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	redisHost  string
	redisPort  int
	redisDb    int
	dbHost     string
	dbPort     int
	dbUsername string
	dbPassword string
	dbName     string
	ranker     *HackerNewsRanker
	date       = time.Now()
	redisKey   = fmt.Sprintf("post-scores-%s", date.Format("2006-01-02"))
	ctx        context.Context
)

func setUp() func() {
	once := sync.Once{}
	once.Do(func() {
		redisHost, redisPort, redisDb = utils.LoadRedisEnvVariables()

		cache.NewCacheService(redisHost, redisPort, redisDb)
		ranker = NewRanker()
		ctx = ranker.CacheInstance.Ctx

		dbHost, dbPort, dbUsername, dbPassword, dbName = utils.LoadDBEnvVariables()

		sql.NewSqlService(dbUsername, dbPassword, dbHost, dbPort, dbName)
		post.NewPostService()
	})

	// Teardown function as return value
	return func() {
		sql.DB.Client.MustExec("TRUNCATE TABLE posts")
		ranker.CacheInstance.RedisClient.Del(ctx, redisKey)
	}
}

func TestHackerNewsRanker_PushPostScore(t *testing.T) {
	teardown := setUp()
	defer teardown()

	newPost := post.NewPost("Test Title", "Test Body")

	err := ranker.PushPostScore(newPost)
	assert.NoError(t, err)

	result, keyErr := ranker.CacheInstance.RedisClient.Exists(ctx, redisKey).Result()
	assert.NoError(t, keyErr)
	assert.Equal(t, result, int64(1))
}

func TestHackerNewsRanker_calculateScore(t *testing.T) {
	teardown := setUp()
	defer teardown()

	points := 10
	creationTime := time.Now().Add(-24 * time.Hour) // Set creation time to 24 hours ago
	gravity := 1.0

	expectedScore := float64(points-1) / math.Pow(24*60*60+2, gravity) // Assuming seconds as the unit of age

	score := ranker.calculateScore(points, creationTime, gravity)
	assert.Equal(
		t,
		strconv.FormatFloat(expectedScore, 'f', 12, 64),
		strconv.FormatFloat(score, 'f', 12, 64),
	)
}

func TestHackerNewsRanker_GetTopRankedPosts(t *testing.T) {
	teardown := setUp()
	defer teardown()

	for i := 0; i < 50; i++ {
		newPost := post.NewPost(fmt.Sprintf("Test Title %d", i), fmt.Sprintf("Test Body %d", i))
		newPost.Timestamp = date
		post.PostServiceInstance.Create(*newPost)

		// Upvote the post
		p, err := post.PostServiceInstance.Vote(newPost.ID, true)
		assert.NoError(t, err)

		// Store into leaderboard
		err = ranker.PushPostScore(p)
		assert.NoError(t, err)
	}

	// Retrieve the top ranked posts
	topPosts := ranker.GetTopRankedPosts(date, 50)

	// Assertions
	assert.NotNil(t, topPosts)
	assert.Equal(t, 50, len(topPosts))
	// Assert the order of the top posts based on their scores
}

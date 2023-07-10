package ranking

import (
	"fmt"
	"github.com/n1207n/real-time-post-recommender/cache"
	"github.com/n1207n/real-time-post-recommender/post"
	"github.com/n1207n/real-time-post-recommender/sql"
	"github.com/n1207n/real-time-post-recommender/utils"
	"github.com/stretchr/testify/assert"
	"math"
	"strconv"
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
)

func setUp() {
	redisHost, redisPort, redisDb = utils.LoadRedisEnvVariables()

	cache.NewCacheService(redisHost, redisPort, redisDb)
	ranker = &HackerNewsRanker{
		CacheInstance: cache.Cache,
	}

	dbHost, dbPort, dbUsername, dbPassword, dbName = utils.LoadDBEnvVariables()

	sql.NewSqlService(dbUsername, dbPassword, dbHost, dbPort, dbName)
	post.NewPostService()
}

func TestHackerNewsRanker_PushPostScore(t *testing.T) {
	setUp()

	newPost := post.NewPost("Test Title", "Test Body")

	err := ranker.PushPostScore(newPost)
	assert.NoError(t, err)

	key := fmt.Sprintf("post-scores-%s", newPost.Timestamp.Format("2006-01-02"))
	ctx := ranker.CacheInstance.Ctx

	result, keyErr := ranker.CacheInstance.RedisClient.Exists(ctx, key).Result()
	assert.NoError(t, keyErr)
	assert.Equal(t, result, int64(1))

	// Cleanup
	ranker.CacheInstance.RedisClient.ZRem(ctx, key, newPost.ID.String())
}

func TestHackerNewsRanker_calculateScore(t *testing.T) {
	setUp()

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
	setUp()

	const dataN = 50
	for range [dataN]int{} {
		newPost := post.NewPost("Test Title", "Test Body")
		post.PostServiceInstance.Create(*newPost)

		// Upvote the post
		_, err := post.PostServiceInstance.Vote(newPost.ID, true)
		assert.NoError(t, err)
	}

	// Retrieve the top ranked posts
	topPosts, err := ranker.GetTopRankedPosts(dataN)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, topPosts)
	assert.Equal(t, dataN, len(topPosts))
	// Assert the order of the top posts based on their scores

	// Cleanup
	key := fmt.Sprintf("post-scores-%s", time.Now().Format("2006-01-02"))
	ctx := ranker.CacheInstance.Ctx

	sql.DB.Client.MustExec("TRUNCATE TABLE posts")
	ranker.CacheInstance.RedisClient.Del(ctx, key)
}

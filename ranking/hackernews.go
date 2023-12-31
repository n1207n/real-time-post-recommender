package ranking

import (
	"context"
	"fmt"
	"github.com/n1207n/real-time-post-recommender/cache"
	"github.com/n1207n/real-time-post-recommender/post"
	"github.com/redis/go-redis/v9"
	"math"
	"time"
)

type HackerNewsRanker struct {
	CacheInstance *cache.CacheService
}

var _ HackerNewsRanker = HackerNewsRanker{}

var (
	Ranker *HackerNewsRanker
)

func NewRanker() *HackerNewsRanker {
	Ranker = &HackerNewsRanker{CacheInstance: cache.Cache}
	return Ranker
}

// HackerNewsScore calculates the ranking score based on public formula of HackerNews ranking
// `score = (points - 1) / (age + 2)^gravity`
func (r *HackerNewsRanker) calculateScore(points int, creationTime time.Time, gravity float64) float64 {
	if gravity == 0 {
		gravity = 1.0
	}

	age := time.Now().Sub(creationTime).Seconds()
	return float64(points-1) / math.Pow(age+2, gravity)
}

func (r *HackerNewsRanker) PushPostScore(p *post.Post) error {
	ctx := context.Background()
	key := fmt.Sprintf("post-scores:%s", p.Timestamp.Format("2006-01-02"))

	// TODO: Dynamically adjust gravity based on the # of posts and average gravity or something else
	score := r.calculateScore(p.Votes, p.Timestamp, r.computeGravity())

	// Upsert sorted set entry
	_, err := cache.Cache.RedisClient.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: p.ID.String(),
	}).Result()

	if err != nil {
		return err
	}

	return nil
}

func (r *HackerNewsRanker) computeGravity() float64 {
	return 1.0
}

func (r *HackerNewsRanker) GetTopRankedPosts(date time.Time, n int) []post.Post {
	ctx := context.Background()
	key := fmt.Sprintf("post-scores:%s", date.Format("2006-01-02"))

	var topPosts []post.Post

	postIds, err := cache.Cache.RedisClient.ZRevRange(ctx, key, 0, int64(n-1)).Result()

	if err != nil {
		return topPosts
	}

	topPosts = post.PostServiceInstance.FilterByIDs(postIds)
	return topPosts
}

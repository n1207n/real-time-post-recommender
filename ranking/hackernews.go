package ranking

import (
	"fmt"
	"github.com/n1207n/real-time-post-recommender/cache"
	"github.com/n1207n/real-time-post-recommender/post"
	"github.com/redis/go-redis/v9"
	"math"
	"time"
)

type HackerNewsRanker struct {
	cacheInstance *cache.CacheService
}

var _ HackerNewsRanker = HackerNewsRanker{}

var (
	Ranker *HackerNewsRanker
)

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
	ctx := r.cacheInstance.Ctx

	key := fmt.Sprintf("post-scores-%s", p.Timestamp.Format("2006-01-02"))

	// TODO: Dynamically adjust gravity based on the # of posts and average gravity or something else
	score := r.calculateScore(p.Votes, p.Timestamp, r.computeGravity())

	err := cache.Cache.RedisClient.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: p.ID.String(),
	}).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *HackerNewsRanker) computeGravity() float64 {
	return 1.0
}

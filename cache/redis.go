package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	redisClient *redis.Client
}

// Compilation check
var _ CacheService = CacheService{}

var (
	Cache *CacheService
	ctx   = context.Background()
)

func NewCacheService(redisHost string, redisPort int, redisDB int) *CacheService {
	redisAddr := fmt.Sprintf("%s:%d", redisHost, redisPort)

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   redisDB,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error initializing redis = {%s}", pong))
	}

	fmt.Printf("\nRedis started successfully: pong message = {%s}", pong)

	Cache = &CacheService{
		redisClient: redisClient,
	}

	return Cache
}

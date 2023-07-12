package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	RedisClient *redis.Client
}

// Compilation check
var _ CacheService = CacheService{}

var (
	Cache *CacheService
)

func NewCacheService(redisHost string, redisPort int, redisDB int) *CacheService {
	redisAddr := fmt.Sprintf("%s:%d", redisHost, redisPort)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       redisDB,
	})

	Cache = &CacheService{
		RedisClient: redisClient,
	}

	ctx := context.Background()
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error initializing redis = {%s}", pong))
	}

	fmt.Printf("\nRedis address: %s - Redis DB: %d\n", redisAddr, redisDB)

	fmt.Printf("\nRedis started successfully: pong message = {%s}\n", pong)

	return Cache
}

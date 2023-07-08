package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	RedisClient *redis.Client
	Ctx         context.Context
}

// Compilation check
var _ CacheService = CacheService{}

var (
	Cache *CacheService
)

func NewCacheService(redisHost string, redisPort int, redisDB int) *CacheService {
	redisAddr := fmt.Sprintf("%s:%d", redisHost, redisPort)

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   redisDB,
	})

	Cache = &CacheService{
		RedisClient: redisClient,
		Ctx:         context.Background(),
	}

	pong, err := redisClient.Ping(Cache.Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error initializing redis = {%s}", pong))
	}

	fmt.Printf("\nRedis started successfully: pong message = {%s}", pong)

	return Cache
}

package main

import (
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/n1207n/real-time-post-recommender/cache"
	"github.com/n1207n/real-time-post-recommender/post"
	"github.com/n1207n/real-time-post-recommender/ranking"
	"github.com/n1207n/real-time-post-recommender/routing"
	"github.com/n1207n/real-time-post-recommender/sql"
	"github.com/n1207n/real-time-post-recommender/utils"
	"go.uber.org/zap"
	"log"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	apiPort := utils.LoadApiPortEnvVariable()
	dbHost, dbPort, dbUsername, dbPassword, dbName := utils.LoadDBEnvVariables()
	redisHost, redisPort, redisDb := utils.LoadRedisEnvVariables()

	initializeServices(dbUsername, dbPassword, dbHost, dbPort, dbName, redisHost, redisPort, redisDb)

	r := gin.Default()
	routing.BuildRouters(r)

	// Logger setup
	logger, _ := zap.NewProduction()

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))

	address := fmt.Sprintf(":%d", apiPort)

	err := r.Run(address)
	if err != nil {
		panic(fmt.Sprintf("Failed to start up the web server - Error %v", err))
	}

	return nil
}

func initializeServices(dbUsername string, dbPassword string, dbHost string, dbPort int, dbName string, redisHost string, redisPort int, redisDb int) {
	sql.NewSqlService(dbUsername, dbPassword, dbHost, dbPort, dbName)
	cache.NewCacheService(redisHost, redisPort, redisDb)
	post.NewPostService()
	ranking.NewRanker()
}

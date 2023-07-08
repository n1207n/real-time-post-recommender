package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/n1207n/real-time-post-recommender/cache"
	"github.com/n1207n/real-time-post-recommender/post"
	"github.com/n1207n/real-time-post-recommender/ranking"
	"github.com/n1207n/real-time-post-recommender/routing"
	"github.com/n1207n/real-time-post-recommender/sql"
	"github.com/n1207n/real-time-post-recommender/utils"
	"log"
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

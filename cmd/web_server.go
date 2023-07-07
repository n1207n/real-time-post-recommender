package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/n1207n/real-time-post-recommender/cache"
	"github.com/n1207n/real-time-post-recommender/post"
	"github.com/n1207n/real-time-post-recommender/routing"
	"github.com/n1207n/real-time-post-recommender/sql"
	"log"
	"os"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	apiPort := loadApiPortEnvVariable()
	dbHost, dbPort, dbUsername, dbPassword, dbName := loadDBEnvVariables()
	redisHost, redisPort, redisDb := loadRedisEnvVariables()

	r := gin.Default()

	routing.BuildRouters(r)
	sql.NewSqlService(dbUsername, dbPassword, dbHost, dbPort, dbName)
	cache.NewCacheService(redisHost, redisPort, redisDb)
	post.NewPostService()

	address := fmt.Sprintf(":%d", apiPort)

	err := r.Run(address)
	if err != nil {
		panic(fmt.Sprintf("Failed to start up the web server - Error %v", err))
	}

	return nil
}

func loadDBEnvVariables() (string, int, string, string, string) {
	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Failed to parse DB_PORT: %v", err)
	}

	return dbHost, dbPort, dbUsername, dbPassword, dbName
}

func loadRedisEnvVariables() (string, int, int) {
	redisHost := os.Getenv("REDIS_HOST")
	redisDbStr := os.Getenv("REDIS_DB")
	redisPortStr := os.Getenv("REDIS_PORT")

	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		log.Fatalf("Failed to parse REDIS_PORT: %v", err)
	}

	redisDb, err := strconv.Atoi(redisDbStr)
	if err != nil {
		log.Fatalf("Failed to parse REDIS_DB: %v", err)
	}

	return redisHost, redisPort, redisDb
}

func loadApiPortEnvVariable() int {
	apiPortStr := os.Getenv("API_PORT")

	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil {
		log.Fatalf("Failed to parse API_PORT: %v", err)
	}

	return apiPort
}

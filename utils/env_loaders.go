package utils

import (
	"log"
	"os"
	"strconv"
)

func LoadDBEnvVariables() (string, int, string, string, string) {
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

func LoadRedisEnvVariables() (string, int, int) {
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

func LoadApiPortEnvVariable() int {
	apiPortStr := os.Getenv("API_PORT")

	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil {
		log.Fatalf("Failed to parse API_PORT: %v", err)
	}

	return apiPort
}

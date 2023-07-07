package sql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type SqlService struct {
	Client *sqlx.DB
}

// Compilation check
var _ SqlService = SqlService{}

// Singleton
var DB *SqlService

func NewSqlService(dbUsername string, dbPassword string, dbHost string, dbPort int, dbName string) *SqlService {
	srcName := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbUsername,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	client, err := sqlx.Connect("postgres", srcName)
	if err != nil {
		log.Fatalf("Error connecting postgres = {%v}", err)
	}

	log.Printf("\nPostgres connected")

	// Singleton assignment
	DB = &SqlService{
		Client: client,
	}

	return DB
}

package persistance

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type SqlService struct {
	DbClient *sqlx.DB
}

// Compilation check
var _ SqlService = SqlService{}

// Singleton
var SqlServiceInstance *SqlService

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
	SqlServiceInstance = &SqlService{
		DbClient: client,
	}

	return SqlServiceInstance
}

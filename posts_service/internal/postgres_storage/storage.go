package postgresstorage

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresStorage struct {
	db *gorm.DB
}

const (
	postgresHost     = "POSTGRES_HOST"
	postgresPort     = "POSTGRES_PORT"
	postgresUser     = "POSTGRES_USER"
	postgresPassword = "POSTGRES_PASSWORD"
	postgresDB       = "POSTGRES_DB"
)

func NewPostgresStorage() (*postgresStorage, error) {
	pgHost, exists := os.LookupEnv(postgresHost)
	if !exists {
		return nil, errors.New("no postgres host in env")
	}

	pgPort, exists := os.LookupEnv(postgresPort)
	if !exists {
		return nil, errors.New("no postgres host in env")
	}

	pgUser, exists := os.LookupEnv(postgresUser)
	if !exists {
		return nil, errors.New("no postgres user in env")
	}

	pgPass, exists := os.LookupEnv(postgresPassword)
	if !exists {
		return nil, errors.New("no postgres password in env")
	}

	pgDB, exists := os.LookupEnv(postgresDB)
	if !exists {
		return nil, errors.New("no postgres db in env")
	}

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v",
		pgHost, pgPort, pgUser, pgPass, pgDB)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(postTable{})
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	return &postgresStorage{db: db}, nil
}

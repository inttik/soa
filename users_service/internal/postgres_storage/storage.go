package postgresstorage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"users/oas"

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

	db.AutoMigrate(userInfo{})
	db.AutoMigrate(userProfile{})
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	return &postgresStorage{db: db}, nil
}

func (ps *postgresStorage) MakeRootUser(login oas.LoginString, password oas.PasswordString) error {
	createRequest := oas.CreateUserRequest{
		Login:    login,
		Password: password,
		Email:    oas.EmailString(string(login) + "@root"),
		Root:     oas.NewOptRootFlag(true),
	}
	_, err := ps.CreateUser(&createRequest)
	if err == ErrUserExist {
		return nil
	}
	return err
}

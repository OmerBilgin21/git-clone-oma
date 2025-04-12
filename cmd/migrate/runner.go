package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	var dbUrl = os.Getenv("DATABASE_URL")

	if err != nil {
		log.Fatal("Environment variables could not be load!")
	}

	m, err := migrate.New(
		"file://db/migrations",
		dbUrl,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Steps(1); err != nil {
		log.Fatal(err)
	}
}

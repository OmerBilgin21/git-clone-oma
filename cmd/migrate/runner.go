package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

func reset(path string, dbUrl string) {
	m, err := migrate.New(
		path,
		dbUrl,
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Force(0); err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Fatal()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Environment variables could not be load!")
	}

	var dbUrl = os.Getenv("DATABASE_URL")

	curr, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	migrationsFolderPath := filepath.Join(curr, "db", "migrations")
	list, err := os.ReadDir(migrationsFolderPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("migration list: %v\n", list)

	path := "file://" + migrationsFolderPath + "/"
	m, err := migrate.New(
		path,
		dbUrl,
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Migrations ran successfully!")
}

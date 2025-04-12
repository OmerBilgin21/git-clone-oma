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

func reset(m *migrate.Migrate, path string, dbUrl string) {
	version, dirty, err := m.Version()
	if err != nil {
		log.Print(err)
	}

	log.Printf("current version: %v, is dirty: %v", version, dirty)

	if version > 0 {
		if err := m.Drop(); err != nil {
			log.Print("Yo:", err)
		}

		log.Print("Dropped everything")

		m.Close()
		m, err = migrate.New(path, dbUrl)
		if err != nil {
			log.Fatal(err)
		}

		log.Print("schema migrations table initiated again")
	}

	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

func runMigrations(m *migrate.Migrate) {
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: run-command <mode: 'reset' | 'migrate'>")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Environment variables could not be load!")
	}

	var dbUrl = os.Getenv("DATABASE_URL")

	curr, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	migrationsFolderPath := filepath.Join(curr, "internal", "db", "migrations")
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

	if os.Args[1] == "reset" {
		reset(m, path, dbUrl)
	} else if os.Args[1] == "migrate" {
		runMigrations(m)
	} else {
		log.Fatalf("invalid mode, must be 'reset' or 'migrate'\nExiting...")
	}

	log.Printf("Migrations ran successfully!")
}

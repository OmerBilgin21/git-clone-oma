package internal

import (
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func GetDb() *sqlx.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("env vars could not load: %v", err)
	}

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("error while connecting to db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("something went wrong while pinging the DB: %v", err)
	}

	return db
}

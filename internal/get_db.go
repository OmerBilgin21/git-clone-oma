package internal

import (
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func GetDb() *gorm.DB {
	if err := os.MkdirAll(".oma", 0755); err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("./.oma/oma.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}

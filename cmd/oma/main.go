package main

import (
	"oma/internal/db"
	"oma/pkg"
	"os"
)

func main() {
	cliArgs := os.Args[1:]
	dbIns := db.GetDb()
	pkg.ParseAndDispatch(cliArgs, dbIns)
}

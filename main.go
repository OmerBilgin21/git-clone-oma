package main

import (
	"oma/internal"
	"oma/pkg"
	"os"
)

func main() {
	cliArgs := os.Args[1:]
	dbIns := internal.GetDb()
	pkg.ParseAndDispatch(cliArgs, dbIns)
}

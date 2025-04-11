package main

import (
	"oma/pkg"
	"os"
)

func main() {
	cliArgs := os.Args[1:]

	pkg.ParseAndDispatch(cliArgs)

}

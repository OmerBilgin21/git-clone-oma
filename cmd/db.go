package main

import (
	"log"
	"os"
)

func main() {
	for _, target := range []string{".oma"} {
		if _, err := os.Stat(target); err == nil {
			if err := os.RemoveAll(target); err != nil {
				log.Fatalf("error while removing %v:\n%v", target, err)
			}
			log.Printf("removed %v", target)
		}
	}

	log.Print("reset complete")
}

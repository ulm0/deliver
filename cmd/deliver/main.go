package main

import (
	"log"
	"os"

	"github.com/ulm0/deliver/pkg/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Printf("Error: %w\n", err)
		os.Exit(1)
	}
}

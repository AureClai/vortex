package main

import (
	"log"

	"github.com/AureClai/vortex/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

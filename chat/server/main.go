package main

import (
	seeder "chat/server/jobs/seed-data"
	"fmt"
	"os"
)

func main() {
	arg := os.Args[1]

	server := NewServer()

	if arg == "start" {
		server.Start()
	}

	if arg == "seed" {
		fmt.Println("Seed data")
		seeder.Seed()
	}
}

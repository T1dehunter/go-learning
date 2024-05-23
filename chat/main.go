package main

import (
	"chat/app"
	"os"
)

func main() {
	arg := os.Args[1]

	app := app.NewApp()

	if arg == "start" {
		app.Start()
	}

	if arg == "seed" {
		app.Seed()
	}

}

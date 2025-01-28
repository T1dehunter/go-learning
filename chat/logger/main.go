package main

import (
	"os"
)

func main() {
	loggerClient := NewClient(os.Stdin, os.Stdout)

	loggerClient.Start()
}

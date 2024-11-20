package main

import (
	"os"
)

func main() {
	chatClient := NewClient(os.Stdin, os.Stdout)
	chatClient.Start()
}

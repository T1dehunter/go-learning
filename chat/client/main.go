package main

import (
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	chatClient := NewClient(os.Stdin, os.Stdout, user.Username)
	chatClient.Start()
}

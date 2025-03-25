package main

import (
	"flag"
	"os"
)

func main() {
	user := flag.String("user", "", "Username for authentication")
	flag.Parse()

	chatClient := NewClient(os.Stdin, os.Stdout)
	chatClient.Start(*user)
}

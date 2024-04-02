package main

import (
	"fmt"
	"go-learning/interpreter/repl"
	"os"
	"os/user"
)

// Code in this project is based on the book:
// Writing An Interpreter In Go
// Thorsten Ball
// I want to thank the author for the great book and the great content.
func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the my test programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

package main

import (
	"fmt"
)

func main() {
	fmt.Println("TRIE RUN")

	trie := CreateTrie()

	trie.Insert("Apple")

	trie.Print()

	node, err := trie.Search("Apple")
	if err {
		fmt.Println("No search found")
	} else {
		fmt.Println("Found node: ", node)
	}
}

package main

type Node struct {
	text     string
	isWord   bool
	children map[string]*Node
}

func CreateNode(text string) Node {
	return Node{text: text, isWord: false, children: make(map[string]*Node)}
}

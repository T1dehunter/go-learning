package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Trie struct {
	test string
	root Node
}

func CreateTrie() *Trie {
	trie := &Trie{root: CreateNode("")}
	return trie
}

func (trie *Trie) Insert(word string) {
	currentNode := &trie.root
	for charIndex, char := range strings.Split(word, "") {
		_, isCharInNode := currentNode.children[char]
		if !isCharInNode {
			prefix := word[0 : charIndex+1]
			newNode := CreateNode(prefix)
			currentNode.children[char] = &newNode
			currentNode = &newNode
		} else {
			currentNode = currentNode.children[char]
		}
	}
	currentNode.isWord = true
}

func (trie *Trie) Search(word string) (*Node, bool) {
	currentNode := &trie.root
	var foundNode *Node
	for _, char := range strings.Split(word, "") {
		_, isCharInNode := currentNode.children[char]
		if !isCharInNode {
			return foundNode, true
		}
		currentNode = currentNode.children[char]
	}
	if currentNode.isWord {
		foundNode = currentNode
		return foundNode, false
	}
	return foundNode, true
}

func (trie *Trie) StartsWith(prefix string) bool {
	return false
}

func (trie *Trie) Size() int {
	counter := 0
	size(&trie.root, &counter)
	return counter
}

func size(currentNode *Node, counter *int) {
	for _, val := range currentNode.children {
		*counter += 1
		size(val, counter)
	}
}

func (trie *Trie) Print() {
	fmt.Println("Trie Print")
	currentNode := trie.root
	printNode("root", "", currentNode.isWord, currentNode.children)
}

func printNode(name string, prefix string, isWord bool, children map[string]*Node) {
	nodeName := "name: " + name
	nodePrefix := "prefix: " + prefix
	wordFlag := "isWord: " + strconv.FormatBool(isWord)
	fmt.Println("Node " + nodeName + " --- " + nodePrefix + " --- " + wordFlag)
	for key, val := range children {
		printNode(key, val.text, val.isWord, val.children)
	}
}

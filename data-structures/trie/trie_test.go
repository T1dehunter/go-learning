package main

import (
	"fmt"
	"testing"
)

func Test(test *testing.T) {
	testInsert(test)
	testSearch(test)
	testSize(test)
}

func testInsert(test *testing.T) {
	trie := CreateTrie()

	trie.Insert("Apple")

	if trie.root.text != "" {
		test.Errorf("Error on insert node root node has unexpected text value: %q", trie.root.text)
	}

	node1, isNode1 := trie.root.children["A"]
	if !isNode1 {
		test.Errorf("Error on get 'A' node")
	}
	if node1.text != "A" {
		test.Errorf("Error on get 'A' node value, got: %q", node1.text)
	}
	if node1.isWord {
		test.Errorf("Error on get 'A' node isWord flag, expected: false, got: true")
	}

	node2, isNode2 := node1.children["p"]
	if !isNode2 {
		test.Errorf("Error on get 'p' node")
	}
	if node2.text != "Ap" {
		test.Errorf("Error on get 'p' node value, got: %q", node2.text)
	}
	if node2.isWord {
		test.Errorf("Error on get 'p' node isWord flag, expected: false, got: true")
	}

	node3, isNode3 := node2.children["p"]
	if !isNode3 {
		test.Errorf("Error on get 'p' node")
	}
	if node3.text != "App" {
		test.Errorf("Error on get 'p' node value, got: %q", node3.text)
	}
	if node3.isWord {
		test.Errorf("Error on get 'p' node isWord flag, expected: false, got: true")
	}

	node4, isNode4 := node3.children["l"]
	if !isNode4 {
		test.Errorf("Error on get 'l' node")
	}
	if node4.text != "Appl" {
		test.Errorf("Error on get 'l' node value, got: %q", node4.text)
	}
	if node4.isWord {
		test.Errorf("Error on get 'l' node isWord flag, expected: false, got: true")
	}

	node5, isNode5 := node4.children["e"]
	if !isNode5 {
		test.Errorf("Error on get 'l' node")
	}
	if node5.text != "Apple" {
		test.Errorf("Error on get 'e' node value, got: %q", node5.text)
	}
	if !node5.isWord {
		test.Errorf("Error on get 'e' node isWord flag, expected: true, got: false")
	}
}

func testSearch(test *testing.T) {
	trie := CreateTrie()

	trie.Insert("Apple")

	_, err := trie.Search("Apple")
	if err {
		test.Errorf("Error on search node in trie")
	}

}

func testSize(test *testing.T) {
	testString := "Apple"

	trie := CreateTrie()

	trie.Insert(testString)

	size := trie.Size()
	if size != len(testString) {
		test.Errorf("trie.Size of %q wrong, got=%q", testString, fmt.Sprint(size))
	}
}

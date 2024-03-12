package main

import (
	"fmt"
	"strings"
)

type Token struct {
	tokenType  string
	tokenValue string
}

func isMatch(str string, pattern string) bool {
	patternTokens := parsePatternTokens(pattern)
	isAllSimpleTokens := all(patternTokens, func(token Token) bool {
		return token.tokenType == "simple"
	})
	if isAllSimpleTokens && len(str) < len(patternTokens) {
		return false
	}
	fmt.Println("patternTokens: ", patternTokens)
	var pointer int
	var index int
	for tokenIdx, token := range patternTokens {
		index = tokenIdx
		if token.tokenType == "simple" {
			if pointer >= len(str) {
				break
			}
			if string(str[pointer]) != token.tokenValue {
				return false
			}
			pointer++
		}
		if token.tokenType == "zeroOrMoreChar" {
			if pointer >= len(str) {
				break
			}
			remainStr := str[pointer:len(str)]
			if !strings.Contains(remainStr, token.tokenValue) {
				continue
			}
			if string(str[pointer]) != token.tokenValue {
				return false
			}
			for string(str[pointer]) == token.tokenValue {
				pointer++
				if pointer >= len(str) {
					break
				}
			}
		}
		if token.tokenType == "zeroOrMoreAnyChar" {
			if pointer >= len(str) {
				break
			}
			if tokenIdx == 0 && len(patternTokens) > 1 {
				continue
			}
			pointer = len(str)
			continue
		}
		if token.tokenType == "anySingleChar" {
			if pointer >= len(str) {
				break
			}
			pointer++
		}
	}
	fmt.Println("index", index)
	isAllStrMatched := pointer >= len(str)
	return isAllStrMatched
}

func parsePatternTokens(pattern string) []Token {
	var tokens []Token
	for index, char := range pattern {
		currentChar := string(char)
		isSimpleChar := isSimpleChar(currentChar, index, pattern)
		if isSimpleChar {
			tokens = append(tokens, Token{tokenType: "simple", tokenValue: currentChar})
			continue
		}
		isZeroOrMorePrecedingChar := isZeroOrMorePrecedingChar(currentChar, index, pattern)
		if isZeroOrMorePrecedingChar {
			tokens = append(tokens, Token{tokenType: "zeroOrMoreChar", tokenValue: currentChar})
			continue
		}
		isZeroOrMoreAnyChar := isZeroOrMoreAnyChar(currentChar, index, pattern)
		if isZeroOrMoreAnyChar {
			tokens = append(tokens, Token{tokenType: "zeroOrMoreAnyChar", tokenValue: currentChar})
			continue
		}
		isAnySingleChar := isAnySingleChar(currentChar, index, pattern)
		if isAnySingleChar {
			tokens = append(tokens, Token{tokenType: "anySingleChar", tokenValue: currentChar})
			continue
		}
		if currentChar == "*" {
			continue
		}
	}
	return tokens
}

func all(slice []Token, condition func(token Token) bool) bool {
	for _, value := range slice {
		if !condition(value) {
			return false
		}
	}
	return true
}

func isSimpleChar(char string, charIdx int, str string) bool {
	isChar := char != "*" && char != "."
	if !isChar {
		return false
	}
	nextCharIdx := charIdx + 1
	if nextCharIdx >= len(str) {
		return true
	}
	nextChar := string(str[nextCharIdx])
	if nextChar != "*" {
		return true
	}
	return false
}

func isZeroOrMorePrecedingChar(char string, charIdx int, str string) bool {
	isChar := char != "*" && char != "."
	if !isChar {
		return false
	}
	nextCharIdx := charIdx + 1
	if nextCharIdx >= len(str) {
		return true
	}
	nextChar := string(str[nextCharIdx])
	if nextChar == "*" {
		return true
	}
	return false
}

func isZeroOrMoreAnyChar(char string, charIdx int, str string) bool {
	isAnyChar := char == "."
	if !isAnyChar {
		return false
	}
	nextCharIdx := charIdx + 1
	if nextCharIdx >= len(str) {
		return false
	}
	nextChar := string(str[nextCharIdx])
	if nextChar == "*" {
		return true
	}
	return false
}

func isAnySingleChar(char string, charIdx int, str string) bool {
	isAnyChar := char == "."
	if !isAnyChar {
		return false
	}
	nextCharIdx := charIdx + 1
	if nextCharIdx >= len(str) {
		return true
	}
	nextChar := string(str[nextCharIdx])
	if nextChar != "*" {
		return true
	}
	return false
}

package main

import (
	"fmt"
	"strings"
)

type Token struct {
	tokenType    string
	tokenValue   string
	coveredChars []string
}

func isMatch(str string, pattern string) bool {
	tokens := parsePatternTokens(pattern)
	isAllZeroOrMoreAnyCharsTokens := all(tokens, func(token Token) bool {
		return token.tokenType == "zeroOrMoreAnyChar"

	})
	if isAllZeroOrMoreAnyCharsTokens {
		return true
	}

	isOneZeroOrMoreChar := len(tokens) == 1 && tokens[0].tokenType == "zeroOrMoreChar"
	if isOneZeroOrMoreChar && !strings.Contains(str, tokens[0].tokenValue) {
		return true
	}
	isAllCharsTokens := all(tokens, func(token Token) bool {
		return token.tokenType == "char"
	})
	if isAllCharsTokens && len(str) != len(tokens) {
		return false
	}
	allCharTokens := filter(tokens, func(token Token) bool {
		return token.tokenType == "char"
	})
	if len(allCharTokens) > len(str) {
		return false
	}
	var charPointer int
	for idx := range tokens {
		token := &tokens[idx]
		if charPointer >= len(str) {
			break
		}
		currentChar := string(str[charPointer])
		switch token.tokenType {
		case "char":
			if currentChar == token.tokenValue {
				token.coveredChars = append(token.coveredChars, currentChar)
				charPointer++
			} else {
				return false
			}
		case "zeroOrMoreChar":
			for string(str[charPointer]) == token.tokenValue {
				fmt.Println("charPointer charPointer : ", string(str[charPointer]))
				token.coveredChars = append(token.coveredChars, string(str[charPointer]))
				if charPointer < len(str) {
					charPointer++
				}
				if charPointer == len(str) {
					break
				}
			}
		case "zeroOrMoreAnyChar":
			charPointer = len(str)
		case "anySingleChar":
			if len(str) > 1 && len(tokens) == 1 {
				return false
			}
			token.coveredChars = append(token.coveredChars, currentChar)
			charPointer++
		}
	}
	coveredChars := calculateCoveredChars(tokens)
	return coveredChars == len(str)
}

func parsePatternTokens(pattern string) []Token {
	var tokens []Token
	for index, char := range pattern {
		currentChar := string(char)
		isChar := isChar(currentChar, index, pattern)
		if isChar {
			tokens = append(tokens, Token{tokenType: "char", tokenValue: currentChar})
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

func filter(slice []Token, condition func(token Token) bool) []Token {
	var filtered []Token
	for _, value := range slice {
		if condition(value) {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

func isChar(char string, charIdx int, str string) bool {
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

func calculateCoveredChars(tokens []Token) int {
	var countCoveredChars int
	for _, token := range tokens {
		countCoveredChars += len(token.coveredChars)
	}
	return countCoveredChars
}

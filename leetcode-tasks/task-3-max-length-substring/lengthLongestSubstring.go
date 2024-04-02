package main

import "strings"

func lengthLongestSubstring(str string) int {
	longestSubstringLength := 0
	if len(str) == 1 {
		return 1
	}
	for currentIndex, char := range str {
		longestSubstring := string(char)
		for _, nextChar := range str[currentIndex+1:] {
			nextStringChar := string(nextChar)
			if strings.Contains(longestSubstring, nextStringChar) {
				break
			} else {
				longestSubstring += nextStringChar
			}
		}
		substringLength := len(longestSubstring)
		if substringLength > longestSubstringLength {
			longestSubstringLength = substringLength
		}
	}
	return longestSubstringLength
}

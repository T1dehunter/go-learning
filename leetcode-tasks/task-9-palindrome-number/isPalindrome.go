package main

import "strconv"

func isPalindrome(number int) bool {
	if number == 0 {
		return true
	}
	stringNumber := strconv.Itoa(number)
	reversedStringNumber := reverseString(stringNumber)
	return stringNumber == reversedStringNumber
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

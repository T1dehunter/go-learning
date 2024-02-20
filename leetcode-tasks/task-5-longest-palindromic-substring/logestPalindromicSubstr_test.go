package main

import "testing"

func TestLongestPalindromicSubstr(testFramework *testing.T) {
	tests := []struct {
		string            string
		longestPalindrome string
	}{
		{"babad", "bab"},
		{"ac", "a"},
		{"test", "t"},
	}

	for index, testData := range tests {
		result := longestPalindromeSubstr(testData.string)

		if result != testData.longestPalindrome {
			testFramework.Fatalf("tests[%d] - calculated palindromic substr value is wrong. expected=%s, got=%s", index, testData.longestPalindrome, result)
		}
	}
}

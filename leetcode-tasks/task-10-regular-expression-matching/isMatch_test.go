package main

import (
	"testing"
)

func TestIsMatch(testFramework *testing.T) {
	tests := []struct {
		str     string
		pattern string
		result  bool
	}{
		{"aa", "a", false},
		{"aa", "a*", true},
		{"bb", "a*", true},
		{"aaa", "a*a", true},
		{"ab", ".*", true},
		{"aab", "c*a*b", true},
		{"mississippi", "mis*is*p*.", false},
		{"ab", ".*c", false},
		{"aaa", "aaaa", false},
		{"aaa", "a.a", true},
		{"aa", ".", false},
		{"a", "ab*a", false},
		{"a", "ab*", true},
		{"aa", ".a", true},
		{"aaa", ".a", false},
		{"abcd", "d*", false},
	}

	for index, testData := range tests {
		result := isMatch(testData.str, testData.pattern)

		if result != testData.result {
			testFramework.Fatalf("tests[%d] - matched result wrong. expected=%t, got=%t", index, testData.result, result)
		}
	}
}

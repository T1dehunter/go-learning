package main

import "testing"

func TestLengthLongestSubstring(testFramework *testing.T) {
	tests := []struct {
		string string
		length int
	}{
		{"test", 3},
	}

	for index, testData := range tests {
		result := lengthLongestSubstring(testData.string)

		if result != testData.length {
			testFramework.Fatalf("tests[%d] - calculated substring wrong. expected=%d, got=%d", index, testData.length, result)
		}
	}
}

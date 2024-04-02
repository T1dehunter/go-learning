package main

import "testing"

func TestIsPalindrome(testFramework *testing.T) {
	tests := []struct {
		number       int
		isPalindrome bool
	}{
		{121, true},
		{-121, false},
		{10, false},
	}

	for index, testData := range tests {
		result := isPalindrome(testData.number)

		if result != testData.isPalindrome {
			testFramework.Fatalf("tests[%d] - converted value wrong. expected=%t, got=%t", index, testData.isPalindrome, result)
		}
	}
}

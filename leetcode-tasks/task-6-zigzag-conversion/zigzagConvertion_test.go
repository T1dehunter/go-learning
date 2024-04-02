package main

import "testing"

func TestZigzagConvertion(testFramework *testing.T) {
	tests := []struct {
		zigzagString    string
		numRows         int
		convertedString string
	}{
		{"PAYPALISHIRING", 3, "PAHNAPLSIIGYIR"},
		{"ABCD", 2, "ACBD"},
		{"ABC", 2, "ACB"},
	}

	for index, testData := range tests {
		result := zigzagConvertion(testData.zigzagString, testData.numRows)

		if result != testData.convertedString {
			testFramework.Fatalf("tests[%d] - converted value wrong. expected=%s, got=%s", index, testData.convertedString, result)
		}
	}
}

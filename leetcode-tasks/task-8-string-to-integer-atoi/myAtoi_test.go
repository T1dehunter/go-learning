package main

import "testing"

func TestMyAtoi(testFramework *testing.T) {
	tests := []struct {
		stringValue    string
		convertedValue int
	}{
		{"42", 42},
		{"   -42", -42},
		{"4193 with words", 4193},
		{"0032", 32},
		{"-91283472332", -2147483648},
		{"3.14159", 3},
		{"+-12", 0},
		{"00000-42a1234", 0},
		{"   +0 123", 0},
		{".1", 0},
		{"010", 10},
		{"-2147483647", -2147483647},
		{"0  123", 0},
		{"-13+8", -13},
		{"words and 987", 0},
		{"+1", 1},
		{" b11228552307", 0},
		{" -0012a42", -12},
	}

	for index, testData := range tests {
		result := myAtoi(testData.stringValue)

		if result != testData.convertedValue {
			testFramework.Fatalf("tests[%d] - converted value wrong. expected=%d, got=%d", index, testData.convertedValue, result)
		}
	}
}

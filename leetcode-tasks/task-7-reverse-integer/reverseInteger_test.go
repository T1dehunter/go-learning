package main

import "testing"

func TestReverse(testFramework *testing.T) {
	tests := []struct {
		value         int
		reversedValue int
	}{
		{123, 321},
		{-123, -321},
		{120, 21},
	}

	for index, testData := range tests {
		result := reverseInteger(testData.value)

		if result != testData.reversedValue {
			testFramework.Fatalf("tests[%d] - reversed value wrong. expected=%d, got=%d", index, testData.reversedValue, result)
		}
	}
}

package main

import "testing"

func TestMedianSortedArrays(testFramework *testing.T) {
	tests := []struct {
		firstArr  []int
		secondArr []int
		median    float64
	}{
		{
			[]int{1, 3}, []int{2}, 2.0,
		},
		{
			[]int{0, 0, 0, 0, 0}, []int{-1, 0, 0, 0, 0, 0, 1}, 0.0,
		},
		{
			[]int{1, 3}, []int{2, 7}, 2.5,
		},
		{
			[]int{3}, []int{-2, -1}, -1.0,
		},
		{
			[]int{1, 2}, []int{3, 4}, 2.5,
		},
	}

	for index, testData := range tests {
		result := medianSortedArrays(testData.firstArr, testData.secondArr)

		if result != testData.median {
			testFramework.Fatalf("tests[%d] - calculated median is wrong. expected=%f, got=%f", index, testData.median, result)
		}
	}
}

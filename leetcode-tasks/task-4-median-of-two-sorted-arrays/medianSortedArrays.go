package main

import (
	"fmt"
)

func medianSortedArrays(nums1 []int, nums2 []int) float64 {
	mergedArray := mergeSortedArrays(nums1, nums2)
	fmt.Println("mergedArray: ", mergedArray)
	return findArrayMedian(mergedArray)
}

func mergeSortedArrays(firstArray []int, secondArray []int) []int {
	if len(firstArray) > 0 && len(secondArray) == 0 {
		return firstArray
	}
	if len(secondArray) > 0 && len(firstArray) == 0 {
		return secondArray
	}
	firstArrayPointer := 0
	secondArrayPointer := 0
	var result []int
	var firstArrayValue int
	var secondArrayValue int
	for i := 0; i < (len(firstArray) + len(secondArray)); i++ {
		if firstArrayPointer <= (len(firstArray) - 1) {
			firstArrayValue = firstArray[firstArrayPointer]
		}
		if secondArrayPointer <= (len(secondArray) - 1) {
			secondArrayValue = secondArray[secondArrayPointer]
		}
		isFirstArrayEnd := firstArrayPointer == len(firstArray)
		isSecondArrayEnd := secondArrayPointer == len(secondArray)
		if !isFirstArrayEnd && (firstArrayValue <= secondArrayValue || isSecondArrayEnd) {
			result = append(result, firstArrayValue)
			firstArrayPointer += 1
		}
		if !isSecondArrayEnd && (secondArrayValue <= firstArrayValue || isFirstArrayEnd) {
			result = append(result, secondArrayValue)
			secondArrayPointer += 1
		}
	}
	return result
}

func findArrayMedian(arr []int) float64 {
	arrLength := len(arr)
	isArrLengthEven := arrLength%2 == 0
	var median = 0.00
	if isArrLengthEven {
		median = findMedianForEvenArray(arr)
	} else {
		median = findMedianForOddArray(arr)
	}
	return median
}

func findMedianForEvenArray(arr []int) float64 {
	arrLength := len(arr)
	leftMiddleItemIndex := (arrLength / 2) - 1
	rightMiddleItemIndex := leftMiddleItemIndex + 1
	leftMiddleItem := arr[leftMiddleItemIndex]
	rightMiddleItem := arr[rightMiddleItemIndex]
	return (float64(leftMiddleItem) + float64(rightMiddleItem)) / 2
}

func findMedianForOddArray(arr []int) float64 {
	arrLength := len(arr)
	medianItemIndex := ((arrLength + 1) / 2) - 1
	medianItem := arr[medianItemIndex]
	return float64(medianItem)
}

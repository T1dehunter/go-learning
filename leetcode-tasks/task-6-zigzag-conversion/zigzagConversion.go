package main

import (
	"fmt"
	"strings"
)

func zigzagConvertion(zigZagString string, numRows int) string {
	rows := make([][]string, numRows)
	rowIndex := 0
	verticalCharsLength := 0
	diagonalCharsLength := 0
	maxLengthDiagonalChars := numRows - 2
	tooShortString := maxLengthDiagonalChars < 0
	if tooShortString {
		return zigZagString
	}
	for _, char := range zigZagString {
		if verticalCharsLength < numRows {
			rows[rowIndex] = append(rows[rowIndex], string(char))
			rowIndex += 1
			verticalCharsLength += 1
			if verticalCharsLength == numRows && maxLengthDiagonalChars == 0 {
				verticalCharsLength = 0
				rowIndex = 0
			}
		} else {
			diagonalCharsLength += 1
			rowForDiagonalChar := (numRows - diagonalCharsLength) - 1
			rows[rowForDiagonalChar] = append(rows[rowForDiagonalChar], string(char))
			if diagonalCharsLength == maxLengthDiagonalChars {
				diagonalCharsLength = 0
				verticalCharsLength = 0
				rowIndex = 0
			}
		}
	}
	fmt.Println("ROWS :", rows)
	res := ""
	for i := 0; i < len(rows); i++ {
		res += strings.Join(rows[i], "")
	}
	return res
}

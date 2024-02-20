package main

import (
	"strconv"
	"strings"
	"unicode"
)

func myAtoi(value string) int {
	foundNumbers := ""
	isLeadingZeros := false
	numberChars := "1234567890"
	isMinusOperator := false
	isOnlyNumbersAllowed := false
	for index, char := range value {
		isFirstChar := index == 0
		currentChar := string(char)
		isLetter := unicode.IsLetter(char)
		if isLetter {
			break
		}
		isStartingLeadingZeros := (isFirstChar && currentChar == "0") || (currentChar == "0" && foundNumbers == "")
		if isStartingLeadingZeros {
			isLeadingZeros = true
			isOnlyNumbersAllowed = true
			continue
		}
		isLeadingZerosContinue := isLeadingZeros && currentChar == "0"
		if isLeadingZerosContinue {
			continue
		}
		isCurrentCharNumber := strings.Contains(numberChars, currentChar)
		if !isCurrentCharNumber && isFirstChar && currentChar == "." {
			break
		}
		if !isCurrentCharNumber && isOnlyNumbersAllowed {
			break
		}
		if !isCurrentCharNumber && currentChar == "-" {
			isMinusOperator = true
			isOnlyNumbersAllowed = true
			continue
		}
		if !isCurrentCharNumber && currentChar == "+" {
			isOnlyNumbersAllowed = true
			continue
		}
		if !isCurrentCharNumber && currentChar == " " {
			continue
		}
		if isCurrentCharNumber {
			isLeadingZeros = false
			isOnlyNumbersAllowed = true
			foundNumbers += currentChar
			continue
		}
	}
	params := ParseNumberParams{
		strNumbers: foundNumbers,
		isNegative: isMinusOperator,
	}
	return parseStringNumbers(params)
}

func MyAtoiOld(value string) int {
	foundNumbers := ""
	isLeadingZeros := false
	numberChars := "1234567890"
	isPlusOperator := false
	isMinusOperator := false
	for index, char := range value {
		currentChar := string(char)
		if currentChar == "0" && index == 0 {
			isLeadingZeros = true
			continue
		}
		if isLeadingZeros && currentChar == "0" {
			continue
		}
		if isLeadingZeros && index > 0 && (currentChar == "-" || currentChar == "+") {
			return 0
		}
		if isLeadingZeros && index > 0 && currentChar == " " {
			return 0
		}
		if isLeadingZeros && index > 0 && currentChar != "0" {
			isLeadingZeros = false
		}
		if isPlusOperator && isMinusOperator {
			return 0
		}
		isEmptyChar := currentChar == " "
		if isEmptyChar && len(foundNumbers) == 0 {
			continue
		}
		if isEmptyChar && len(foundNumbers) > 0 {
			break
		}
		isLetter := unicode.IsLetter(char)
		if isLetter && len(foundNumbers) == 0 {
			return 0
		}
		if isLetter && len(foundNumbers) > 0 {
			break
		}
		isPoint := currentChar == "."
		if isPoint && len(foundNumbers) == 0 {
			return 0
		}
		if isPoint && len(foundNumbers) > 0 {
			break
		}
		isNumberChar := strings.Contains(numberChars, currentChar)
		if isNumberChar {
			foundNumbers += string(char)
			continue
		}
		isMinusChar := currentChar == "-"
		if isMinusChar {
			isMinusOperator = true
			continue
		}
		isPlusChar := currentChar == "+"
		if isPlusChar {
			isPlusOperator = true
		}
	}
	params := ParseNumberParams{
		strNumbers: foundNumbers,
		isNegative: isMinusOperator,
	}
	return parseStringNumbers(params)
}

type ParseNumberParams struct {
	strNumbers string
	isNegative bool
}

func parseStringNumbers(params ParseNumberParams) int {
	const maxValue = 2147483647
	if len(params.strNumbers) == 0 {
		return 0
	}
	parsed64Int, _ := strconv.ParseInt(params.strNumbers, 10, 64)
	parsedInt, _ := strconv.ParseInt(params.strNumbers, 10, 32)
	if params.isNegative {
		negativeInt := -parsedInt
		if isWithin32BitRange(parsed64Int) {
			return int(negativeInt)
		} else {
			return int(negativeInt - 1)
		}
	}
	return int(parsedInt)
}

func isWithin32BitRange(num int64) bool {
	const max32BitInt = 2147483647
	const min32BitInt = -2147483648

	return num >= min32BitInt && num <= max32BitInt
}

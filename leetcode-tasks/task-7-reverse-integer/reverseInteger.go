package main

import (
	"strconv"
)

func reverseInteger(value int) int {
	isNegative := value < 0
	intStr := strconv.FormatInt(int64(value), 10)
	var reversedStr string
	if isNegative {
		reversedStr = reverseString(intStr[1:])
	} else {
		reversedStr = reverseString(intStr)
	}
	reversedInt, error := strconv.ParseInt(reversedStr, 10, 32)
	if error != nil {
		return 0
	}
	if isNegative {
		return int(reversedInt - reversedInt - reversedInt)
	}
	return int(reversedInt)
}

func reverseString(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

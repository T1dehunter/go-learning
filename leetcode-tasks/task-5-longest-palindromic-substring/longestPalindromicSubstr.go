package main

func longestPalindromeSubstr(str string) string {
	if len(str) == 1 {
		return str
	}
	if isPalindrome(str) {
		return str
	}
	longestPalindrome := ""
	for currentIndex, char := range str {
		substr := string(char)
		for _, nextChar := range str[currentIndex+1:] {
			substr += string(nextChar)
			if (isPalindrome(substr)) && len(substr) > len(longestPalindrome) {
				longestPalindrome = substr
			}
		}
	}
	if len(longestPalindrome) == 0 {
		return string(str[0])
	}
	return longestPalindrome
}

func isPalindrome(s string) bool {
	firstSymbol := s[0]
	lastSymbol := s[len(s)-1]
	return (firstSymbol == lastSymbol) && s == reverse(s)
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

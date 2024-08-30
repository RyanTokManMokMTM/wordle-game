package utils

func LetterIndexOf(c rune) int {
	if c >= 'a' && c <= 'z' {
		return int(c - 'a')
	} else if c >= 'A' && c <= 'Z' {
		return int(c-'A') + 26
	} else {
		return -1
	}
}

package regex

import "regexp"

const (
	FiveLetterWordMatcher = "^[A-Za-z]{5}$"
)

func Regex(text string, regexStr string) bool {
	return regexp.MustCompile(regexStr).MatchString(text)
}

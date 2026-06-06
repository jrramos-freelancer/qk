package functions

import (
	"regexp"
)

func normalizeSpaces(input string) string {
	spaceRegex := regexp.MustCompile(`\ +`)
	return spaceRegex.ReplaceAllString(input, " ")
}

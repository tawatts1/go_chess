package utility

import (
	"regexp"
)

func RemoveWhitespace(s string) string {
	space := regexp.MustCompile(`\s+`)
	str := space.ReplaceAllString(s, "")
	return str
}

package utility

import (
	"regexp"
	"strconv"
	"strings"
)

func RemoveWhitespace(s string) string {
	space := regexp.MustCompile(`\s+`)
	str := space.ReplaceAllString(s, "")
	return str
}

func StrToInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

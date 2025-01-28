package utility

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

const Epsilon = 0.0001

var Infinity float64 = 1000000

func RemoveWhitespace(s string) string {
	space := regexp.MustCompile(`\s+`)
	str := space.ReplaceAllString(s, "")
	return str
}

func StrToInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

func IsClose(num1, num2 float64) bool {
	return math.Abs(num1-num2) < Epsilon
}

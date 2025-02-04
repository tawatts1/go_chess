package utility

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

const Epsilon = 0.000001

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

func IsApproxGreaterThanOrEq(num1, num2 float64) bool {
	return num1 > num2 || IsClose(num1, num2)
}

func SafelyIncrementAtIndex(slc []int, i int) []int {
	if len(slc)-1 < i {
		for k := len(slc) - 1; k < i; k++ {
			slc = append(slc, 0)
		}
	}
	slc[i] += 1
	return slc
}

func SumSlice(slc []int) int {
	out := 0
	for _, val := range slc {
		out += val
	}
	return out
}

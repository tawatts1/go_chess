package utility

import (
	"fmt"
	"testing"
)

func TestSafelyIncrementAtIndex(t *testing.T) {
	counters := make([]int, 0)

	counters = SafelyIncrementAtIndex(counters, 1)
	fmt.Println(counters)
}

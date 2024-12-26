// sum.go file
package main

import "C"
import "time"

//export sum
func sum(a C.int, b C.int) C.int {
	return a + b
}

//export longSum
func longSum(a C.int, b C.int) C.int {
	time.Sleep(4 * time.Second)
	return a + b + 100
}

//export enforce_binding
func enforce_binding() {}

func main() {}

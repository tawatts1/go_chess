package board

import (
	"fmt"
	"testing"
)

func TestGetBoardFromString(t *testing.T) {
	b := GetBoardFromString(StartingBoard)
	fmt.Println(b)
}

func TestGetBoardFromFile(t *testing.T) {
	b := GetBoardFromFile("kpo.txt")
	fmt.Println(b)

	b2 := GetBoardFromFile("qgo.txt")
	fmt.Println(b2)
}

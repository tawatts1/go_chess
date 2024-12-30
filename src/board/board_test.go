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

func TestGetWhiteCoords(t *testing.T) {
	b := GetBoardFromString(StartingBoard)
	fmt.Println("Checking white coordinates: ")
	fmt.Println(b)
	fmt.Println(b.GetWhiteCoords())
}

func TestGetBlackCoords(t *testing.T) {
	b := GetBoardFromString(StartingBoard)
	fmt.Println("Checking black coordinates: ")
	fmt.Println(b)
	fmt.Println(b.GetBlackCoords())
}

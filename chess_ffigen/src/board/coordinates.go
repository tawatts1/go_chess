package board

import "fmt"

// Coord, yth row down, xth column to the right.
type Coord struct {
	y, x int
}

var UnsetCoord = Coord{y: -1, x: -1}
var BottomLeft = Coord{y: BoardHeight - 1, x: 0}
var BottomRight = Coord{y: BoardHeight - 1, x: BoardWidth - 1}
var TopLeft = Coord{y: 0, x: 0}
var TopRight = Coord{y: 0, x: BoardWidth - 1}
var WhiteKingHome = Coord{y: BoardHeight - 1, x: 4}
var BlackKingHome = Coord{y: 0, x: 4}

var AllCoordinates [BoardHeight * BoardWidth]Coord = [BoardHeight * BoardWidth]Coord(GetAllCoordinates())

func GetAllCoordinates() [BoardHeight * BoardWidth]Coord {
	var out [BoardHeight * BoardWidth]Coord
	var counter int = 0
	for i := range BoardHeight {
		for j := range BoardWidth {
			out[counter] = Coord{y: i, x: j}
			counter++
		}
	}
	return out
}

func NewCoord(y, x int) Coord {
	return Coord{y: y, x: x}
}

func (c Coord) String() string {
	return fmt.Sprintf("(y:%v,x:%v)", c.y, c.x)
}

func (c Coord) Encode() string {
	return fmt.Sprintf("%v,%v", c.y, c.x)
}

func (c Coord) Copy() Coord {
	copy := Coord{y: c.y,
		x: c.x}
	return copy
}

func (c Coord) Add(yadd, xadd int) Coord {
	c.y += yadd
	c.x += xadd
	return c
}

func (c1 Coord) Equals(c2 Coord) bool {
	return c1.y == c2.y && c1.x == c2.x
}

func (c Coord) IsInBoard() bool {
	return c.x >= 0 && c.y >= 0 && c.x < BoardWidth && c.y < BoardHeight
}

func (c Coord) GetY() int {
	return c.y
}

func (c Coord) GetX() int {
	return c.x
}

package board

import "fmt"

// coord, yth row down, xth column to the right.
type coord struct {
	y, x int
}

var UnsetCoord = coord{y: -1, x: -1}
var BottomLeft = coord{y: BoardHeight - 1, x: 0}
var BottomRight = coord{y: BoardHeight - 1, x: BoardWidth - 1}
var TopLeft = coord{y: 0, x: 0}
var TopRight = coord{y: 0, x: BoardWidth - 1}
var WhiteKingHome = coord{y: BoardHeight - 1, x: 4}
var BlackKingHome = coord{y: 0, x: 4}

func (c coord) String() string {
	return fmt.Sprintf("(y:%v,x:%v)", c.y, c.x)
}

func (c coord) Copy() coord {
	copy := coord{y: c.y,
		x: c.x}
	return copy
}

func (c coord) Add(yadd, xadd int) coord {
	c.y += yadd
	c.x += xadd
	return c
}

func (c1 coord) Equals(c2 coord) bool {
	return c1.y == c2.y && c1.x == c2.x
}

func (c coord) IsInBoard() bool {
	return c.x >= 0 && c.y >= 0 && c.x < BoardWidth && c.y < BoardHeight
}

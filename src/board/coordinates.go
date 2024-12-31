package board

import "fmt"

// coord, y is the column, x is the row of the board.
type coord struct {
	y, x int
}

func (c coord) String() string {
	return fmt.Sprintf("(y:%v,x:%v)", c.y, c.x)
}

func (c coord) Copy() coord {
	copy := coord{y: c.y, x: c.y}
	return copy
}

func (c coord) Add(yadd, xadd int) coord {
	c.y += yadd
	c.x += xadd
	return c
}

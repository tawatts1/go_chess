package board

import "fmt"

// coord, yth row down, xth column to the right.
type coord struct {
	y, x int
}

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

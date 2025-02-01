package ai

import (
	"fmt"

	"github.com/tawatts1/go_chess/board"
)

type getScoreArgs struct {
	b       board.Board
	integer int
	isWhite bool
}

func (gsa1 getScoreArgs) Equals(gsa2 getScoreArgs) bool {
	return gsa1.b.Equals(gsa2.b) && gsa1.integer == gsa2.integer && gsa1.isWhite == gsa2.isWhite
}

func newScoreArgs(b board.Board, integer int, isWhite bool) getScoreArgs {
	return getScoreArgs{b: b, integer: integer, isWhite: isWhite}
}

func (gsa getScoreArgs) String() string {
	return fmt.Sprintf("{%v, %v, %v}", gsa.b.Encode()[:5], gsa.integer, gsa.isWhite)
}

type lruCache struct {
	maxSize  int
	scoreMap map[getScoreArgs]float64
	keyStack []getScoreArgs
}

func NewLRUCache(maxSize int) lruCache {
	stack := make([]getScoreArgs, 0)
	c := lruCache{maxSize: maxSize,
		scoreMap: make(map[getScoreArgs]float64),
		keyStack: stack}
	return c
}

func (c *lruCache) SetNewKey(key getScoreArgs, value float64) {
	_, isKeySet := c.scoreMap[key]
	if isKeySet {
		panic("You cannot overwrite a key with this LRU cache implementation")
	}
	if c.maxSize == len(c.keyStack) {
		deletedKey := c.keyStack[0]
		//fmt.Println(deletedKey.String())
		delete(c.scoreMap, deletedKey)
		c.scoreMap[key] = value
		c.keyStack = c.keyStack[1:]
		c.keyStack = append(c.keyStack, key)
	} else {
		c.scoreMap[key] = value
		c.keyStack = append(c.keyStack, key)
	}
	// fmt.Print("Current stack: ")
	// for _, k := range *c.keyStack {
	// 	fmt.Printf("%v, ", k)
	// }
	// fmt.Println()
	// fmt.Print("Current map values: ")
	// for _, v := range c.scoreMap {
	// 	fmt.Printf("%v, ", v)
	// }
	// fmt.Print("\n\n")
}

func (c lruCache) Get(key getScoreArgs) (float64, bool) {
	out, ok := c.scoreMap[key]
	if !ok {
		return out, false
	}
	for i := range len(c.keyStack) {
		if (c.keyStack)[i].Equals(key) {
			c.keyStack = append(c.keyStack[:i], c.keyStack[i+1:]...)
			c.keyStack = append(c.keyStack, key)
		}
	}
	return out, true
}

func (c lruCache) Has(key getScoreArgs) bool {
	_, ok := c.scoreMap[key]
	return ok
}

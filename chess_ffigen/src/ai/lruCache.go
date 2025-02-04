package ai

import (
	"fmt"

	"github.com/tawatts1/go_chess/board"
	"github.com/tawatts1/go_chess/utility"
)

type getScoreArgs struct {
	board_  board.Board
	depth   int
	isWhite bool
}

func (gsa1 getScoreArgs) Equals(gsa2 getScoreArgs) bool {
	return gsa1.board_.Equals(gsa2.board_) && gsa1.depth == gsa2.depth && gsa1.isWhite == gsa2.isWhite
}

func newScoreArgs(b board.Board, integer int, isWhite bool) getScoreArgs {
	return getScoreArgs{board_: b, depth: integer, isWhite: isWhite}
}

func (gsa getScoreArgs) String() string {
	return fmt.Sprintf("{%v, %v, %v}", gsa.board_.Encode()[:5], gsa.depth, gsa.isWhite)
}

func (gsa getScoreArgs) GetDepth() int {
	return gsa.depth
}

type lruCache struct {
	maxSize  int
	scoreMap map[getScoreArgs]float64
	keyStack []getScoreArgs
}

func NewLRUCache(maxSize int) *lruCache {
	stack := make([]getScoreArgs, 0)
	c := lruCache{maxSize: maxSize,
		scoreMap: make(map[getScoreArgs]float64),
		keyStack: stack}
	return &c
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

func (c *lruCache) Get(key getScoreArgs) (float64, bool) {
	out, ok := c.scoreMap[key]
	if !ok {
		return out, false
	} else {
		for i := range len(c.keyStack) {
			if (c.keyStack)[i].Equals(key) {
				c.keyStack = append(c.keyStack[:i], c.keyStack[i+1:]...)
				c.keyStack = append(c.keyStack, key)
			}
		}
		return out, true
	}

}

func (c *lruCache) Has(key getScoreArgs) bool {
	_, ok := c.scoreMap[key]
	return ok
}

func (c *lruCache) CurrentSize() int {
	return len(c.keyStack)
}

type lruCacheList struct {
	caches               []*lruCache
	cacheUseByDepth      []int
	cacheUseByMovesAhead []int
}

func NewLruCacheList(listSize, cacheSize int) *lruCacheList {
	c := make([]*lruCache, listSize)
	for i := range listSize {
		c[i] = NewLRUCache(cacheSize)
	}
	return &lruCacheList{caches: c,
		cacheUseByDepth:      make([]int, 0),
		cacheUseByMovesAhead: make([]int, 0)}
}

func (cl *lruCacheList) GetLength() int {
	return len(cl.caches)
}

func (cl *lruCacheList) SetNewKey(index int, key getScoreArgs, value float64) {
	if index < cl.GetLength() {
		cl.caches[index].SetNewKey(key, value)
	} else {
		panic("Attempting to set a new key for a cache that doesn't exist!")
	}
}

func (cl *lruCacheList) Get(index int, key getScoreArgs) (float64, bool) {
	if index < cl.GetLength() {
		val, ok := cl.caches[index].Get(key)
		if ok {
			cl.cacheUseByMovesAhead = utility.SafelyIncrementAtIndex(cl.cacheUseByMovesAhead, index)
			cl.cacheUseByDepth = utility.SafelyIncrementAtIndex(cl.cacheUseByDepth, key.depth)
		}
		return val, ok
	} else {
		panic("Attempting to get a value for a cache that doesn't exist!")
	}
}

func (cl *lruCacheList) String() string {
	var out string = "{indexes, sizes: ["
	for i := range cl.GetLength() {
		out += fmt.Sprintf("(%v, %v), ", i, cl.caches[i].CurrentSize())
	}
	out = out[:len(out)-2] +
		fmt.Sprintf("]}\nCache use by depth: %v\nCache use by moves ahead: %v\n", cl.cacheUseByDepth, cl.cacheUseByMovesAhead)
	out += fmt.Sprintf("Total usage: %v\n", utility.SumSlice(cl.cacheUseByDepth))
	return out
}

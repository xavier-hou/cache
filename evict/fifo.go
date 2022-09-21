package evict

import (
	"container/list"
)

// FIFOAdjustEvictFunc need do nothing for FIFO eviction
func FIFOAdjustEvictFunc(elm *list.Element, l *list.List) {
}

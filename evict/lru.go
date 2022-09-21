package evict

import (
	"container/list"
)

func LRUAdjustEvictFunc(elm *list.Element, l *list.List) {
	l.MoveToFront(elm)
}

package evict

import (
	"cache/data"
	"container/list"
	"fmt"
)

func LFUAdjustEvictFunc(elm *list.Element, l *list.List) {
	if l.Front() == l.Back() {
		return
	}

	p := l.Front()

	for p != nil {
		// 找到第一个访问度比elm小的,elm放在其前面
		pVisit, err := DataVisit(p.Value)
		if err != nil {
			panic(err)
		}

		elmVisit, err := DataVisit(elm.Value)
		if err != nil {
			panic(err)
		}
		if pVisit < elmVisit {
			break
		}
		p = p.Next()
	}

	if p != nil {
		l.MoveBefore(elm, p)
		return
	}

	l.MoveToBack(elm)
}

func DataVisit(value any) (int, error) {
	e, ok := value.(*data.CacheData)
	if !ok {
		return -1, fmt.Errorf("get visit from data %v failed, it's not cacheData type", value)
	}
	return e.Visit, nil
}

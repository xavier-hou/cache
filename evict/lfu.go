package evict

import (
	"cache/data"
	"container/list"
)

func LFUAdjustEvictFunc(elm *list.Element, l *list.List) {
	if l.Front() == l.Back() {
		return
	}

	p := l.Front()

	for p != nil {
		// 找到第一个访问度比elm小的,elm放在其前面
		pVisit := data.Get(p.Value, data.TVisit).(int)
		elmVisit := data.Get(elm.Value, data.TVisit).(int)

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

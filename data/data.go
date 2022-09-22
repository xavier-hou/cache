package data

import (
	"reflect"
)

type Type string

const (
	TVisit Type = "Visit"
	TKey   Type = "Key"
	TValue Type = "Value"
)

type cacheData struct {
	// 该数据被访问的次数
	visit int
	// 该数据对应的键
	key string
	// 该数据对应的值
	value any
}

func Get(value any, t Type) any {
	switch t {
	case TVisit:
		return value.(*cacheData).visit
	case TKey:
		res, ok := value.(*cacheData)
		if !ok {
			panic(reflect.TypeOf(value))
		}
		return res.key
	case TValue:
		return value.(*cacheData).value
	}
	return nil
}

func Data(key string, obj any) any {
	return &cacheData{
		visit: 0,
		key:   key,
		value: obj,
	}
}

func Set(value any, t Type, setTo any) {
	switch t {
	case TVisit:
		value.(*cacheData).visit = setTo.(int)
	case TKey:
		value.(*cacheData).key = setTo.(string)
	case TValue:
		value.(*cacheData).value = setTo.(string)
	}
}

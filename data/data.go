package data

type CacheData struct {
	// 该数据被访问的次数
	Visit int
	// 该数据对应的键
	Key string
	// 该数据对应的值
	Value any
}

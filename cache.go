package cache

import (
	"cache/data"
	"cache/evict"
	"container/list"
	"fmt"
	"sync"
)

var DefalutAjdustEvictFunc = evict.LRUAdjustEvictFunc

type AdjustEvictFunc func(elm *list.Element, l *list.List)

type Cache struct {
	// 使用锁来保证对于缓存的并发操作，使用互斥锁而不是读写锁是为了方便起见，因为即使在读取缓存数据时，底层的链表结构也会发生变动
	lock sync.Mutex
	// cache maximum storage unit length
	maxLen int
	// cache currently used storage unit length
	curLen int
	cMap   map[string]*list.Element

	list            *list.List
	adjustEvictFunc AdjustEvictFunc

	nexNum int
}

func New(adjustEvictFunc AdjustEvictFunc, mLen int) *Cache {
	if adjustEvictFunc == nil {
		adjustEvictFunc = DefalutAjdustEvictFunc
	}

	c := &Cache{
		lock:            sync.Mutex{},
		maxLen:          mLen,
		curLen:          0,
		cMap:            make(map[string]*list.Element, mLen),
		list:            list.New(),
		adjustEvictFunc: adjustEvictFunc,
		nexNum:          0,
	}

	return c
}

func (c *Cache) Add(key string, obj interface{}) error {
	// 判断key是否为空值
	if key == "" {
		return fmt.Errorf("key can not be empty")
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	// 判断该key是否已经存在，如果存在，则禁止添加
	_, exist := c.cMap[key]
	if exist {
		return fmt.Errorf("key %q already exist", key)
	}

	// 判断该缓存是否达到空间限制上限，如果达到，触发失效淘汰机制
	if c.curLen >= c.maxLen {
		// 默认链表队尾存放待替换的数据项
		elm := c.list.Back()
		if elm == nil {
			return fmt.Errorf("internal error happen, when eviction strategy is triggered but the back element of list is nil")
		}
		c.list.Remove(elm)
		delete(c.cMap, elm.Value.(*data.CacheData).Key)
		c.curLen--
	}

	// 此处的缓存空间是足够的，按照正常流程构造新的元素及列表
	cd := &data.CacheData{
		Visit: 0,
		Key:   key,
		Value: obj,
	}

	// 将修改后的底层数据加入链表中
	elm := c.list.PushFront(cd)

	// 设置cache map与链表元素对应关系
	c.cMap[key] = elm

	// 根据失效淘汰策略更新函数调整链表
	c.adjustEvictFunc(elm, c.list)

	// 设置一些全局变量
	c.curLen++
	c.nexNum++
	return nil
}

func (c *Cache) Delete(key string) error {
	// 判断key是否为空值
	if key == "" {
		return fmt.Errorf("key can not be empty")
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	elm, exist := c.cMap[key]
	if !exist {
		return fmt.Errorf("key %q does not exist", key)
	}

	c.curLen--
	c.list.Remove(elm)
	delete(c.cMap, key)
	return nil
}

func (c *Cache) Update(key string, obj any) error {
	// 判断key是否为空值
	if key == "" {
		return fmt.Errorf("key can not be empty")
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	elm, exist := c.cMap[key]
	if !exist {
		return fmt.Errorf("key %q does not exist", key)
	}

	// 修改cacheData存储的数据与记录
	cd, assert := elm.Value.(*data.CacheData)
	if !assert {
		return fmt.Errorf("internal error happen, when element %v value is not cacheData type", elm.Value)
	}
	cd.Visit++
	cd.Value = obj

	// elm.Value = cd

	// 根据失效淘汰策略更新函数调整链表
	c.adjustEvictFunc(elm, c.list)

	return nil
}

func (c *Cache) Get(key string) (any, error) {
	// 判断key是否为空值
	if key == "" {
		return nil, fmt.Errorf("key can not be empty")
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	elm, exist := c.cMap[key]
	if !exist {
		return nil, fmt.Errorf("key %q does not exist", key)
	}

	cd, assert := elm.Value.(*data.CacheData)
	if !assert {
		panic(fmt.Errorf("assert failed when update cache"))
	}
	cd.Visit++

	// 根据失效淘汰策略更新函数调整链表
	c.adjustEvictFunc(elm, c.list)
	return cd.Value, nil
}

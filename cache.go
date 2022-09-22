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

		delete(c.cMap, data.Get(elm.Value, data.TKey).(string))
		c.list.Remove(elm)
		c.curLen--
	}

	// 此处的缓存空间是足够的，按照正常流程构造新的元素及列表
	v := data.Data(key, obj)

	// 将修改后的底层数据加入链表中
	elm := c.list.PushFront(v)

	// 设置cache map与链表元素对应关系
	c.cMap[key] = elm

	// 根据失效淘汰策略更新函数调整链表
	c.adjustEvictFunc(elm, c.list)

	// 设置一些全局变量
	c.curLen++
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

	data.Set(elm.Value, data.TVisit, data.Get(elm.Value, data.TVisit).(int)+1)
	data.Set(elm.Value, data.TValue, obj)

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

	data.Set(elm.Value, data.TVisit, data.Get(elm.Value, data.TVisit).(int)+1)
	// 根据失效淘汰策略更新函数调整链表
	c.adjustEvictFunc(elm, c.list)
	return data.Get(elm.Value, data.TValue), nil
}

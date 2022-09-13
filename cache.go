package cache

import (
	"fmt"
	"sync"
)

type cache struct {
	sync.Map
}

func(c *cache) Add(key string,obj interface{}) error {
	return c.Update(key,obj)
}

func (c *cache)Delete(key string) bool {
	_,ok := c.LoadAndDelete(key)
	return ok
}

func (c *cache)Update(key string,obj any) error {
	if key == "" {
		return fmt.Errorf("key can not be empty")
	}
	c.Store(key,obj)

	return nil
}

func (c *cache) Get(key string) (any,bool) {
	return c.Load(key)
}

func (c *cache) List() ([]any) {
	res := make([]any,0)
	c.Range(func(key, value any) bool {
		res = append(res, value)
		return true
	})
	return res
}


//Add(interface{}) error
//Delete(key string) error
//Update(key string,obj interface{}) error
//Get(key string) (interface{},error)
//List() (interface{},error)
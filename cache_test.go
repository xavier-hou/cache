package cache

import (
	"cache/data"
	"cache/evict"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test LRU
//
//1 1
//2 2 -> get(1) 是 1
//3 3 -> get(2) 没有
//4 4 -> get(1) 没有
//	-> get(3) 返回3
//    -> get(4) 返回4

func TestLRU(t *testing.T) {
	c := New(evict.LRUAdjustEvictFunc, 2)

	err := c.Add("1", "1")
	assert.Equal(t, nil, err)

	err = c.Add("2", "2")
	assert.Equal(t, nil, err)

	v, err := c.Get("1")
	assert.Equal(t, nil, err)
	assert.Equal(t, "1", v.(string))

	// after 3 1
	err = c.Add("3", "3")
	assert.Equal(t, nil, err)

	v, err = c.Get("2")
	t.Logf("Get key 2 error is %v, value is %v", err, v)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, nil, v)

	// after 4 3
	err = c.Add("4", "4")
	assert.Equal(t, nil, err)

	v, err = c.Get("1")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, nil, v)
	t.Logf("Get key 1 error is %v", err)

	v, err = c.Get("3")
	assert.Equal(t, nil, err)
	assert.Equal(t, "3", v.(string))

	v, err = c.Get("4")
	assert.Equal(t, nil, err)
	assert.Equal(t, "4", v.(string))

	err = c.Delete("1")
	t.Logf("delete key 1 error is %v", err)
	assert.NotEqual(t, nil, err)

	err = c.Delete("3")
	assert.Equal(t, nil, err)
	v, err = c.Get("3")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, nil, v)
	t.Logf("Get key 3 error is %v", err)

	err = c.Update("1", "1")
	assert.NotEqual(t, nil, err)

	err = c.Update("4", "5")
	assert.Equal(t, nil, err)
	err = c.Update("4", "4")
	assert.Equal(t, nil, err)

	err = c.Delete("4")
	assert.Equal(t, nil, err)
	v, err = c.Get("4")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, nil, v)
	t.Logf("Get key 4 error is %v", err)
}

func TestFIFO(t *testing.T) {
	c := New(evict.FIFOAdjustEvictFunc, 5)

	// 1 2 3 4 5
	//

	err := c.Add("1", "1")
	assert.Equal(t, nil, err)

	err = c.Add("2", "2")
	assert.Equal(t, nil, err)

	err = c.Add("3", "3")
	assert.Equal(t, nil, err)

	err = c.Add("4", "4")
	assert.Equal(t, nil, err)

	err = c.Add("5", "5")
	assert.Equal(t, nil, err)

	v, err := c.Get("1")
	assert.Equal(t, nil, err)
	assert.Equal(t, "1", v.(string))

	v, err = c.Get("2")
	assert.Equal(t, nil, err)
	assert.Equal(t, "2", v.(string))

	v, err = c.Get("3")
	assert.Equal(t, nil, err)
	assert.Equal(t, "3", v.(string))

	v, err = c.Get("4")
	assert.Equal(t, nil, err)
	assert.Equal(t, "4", v.(string))

	v, err = c.Get("5")
	assert.Equal(t, nil, err)
	assert.Equal(t, "5", v.(string))

	c.Add("6", "6")
	t.Logf("add key 6, got list %v", c.list)
	v, err = c.Get("5")
	assert.Equal(t, nil, err)
	assert.Equal(t, "5", v.(string))
	_, err = c.Get("1")
	assert.NotEqual(t, nil, err)

	c.Add("7", "7")
	t.Logf("add key 7, got list %v", c.list)
	v, err = c.Get("5")
	assert.Equal(t, nil, err)
	assert.Equal(t, "5", v.(string))
	_, err = c.Get("2")
	assert.NotEqual(t, nil, err)

	c.Add("8", "8")
	t.Logf("add key 8, got list %v", c.list)
	v, err = c.Get("5")
	assert.Equal(t, nil, err)
	assert.Equal(t, "5", v.(string))
	_, err = c.Get("3")
	assert.NotEqual(t, nil, err)

	c.Add("9", "9")
	t.Logf("add key 9, got list %v", c.list)
	v, err = c.Get("5")
	assert.Equal(t, nil, err)
	assert.Equal(t, "5", v.(string))
	_, err = c.Get("4")
	assert.NotEqual(t, nil, err)

	c.Add("10", "10")
	t.Logf("add key 10, got list %v", c.list)
	_, err = c.Get("5")
	assert.NotEqual(t, nil, err)
}

func TestLFU(t *testing.T) {
	c := New(evict.LFUAdjustEvictFunc, 5)

	// 1 2 3 4 5
	//

	err := c.Add("1", "1")
	assert.Equal(t, nil, err)

	err = c.Add("2", "2")
	assert.Equal(t, nil, err)

	err = c.Add("3", "3")
	assert.Equal(t, nil, err)

	err = c.Add("4", "4")
	assert.Equal(t, nil, err)

	err = c.Add("5", "5")
	assert.Equal(t, nil, err)

	// 42351
	// 54321

	arr := []string{"4", "2", "3", "5", "1"}
	k := 0
	for i := 5; i > 0; i-- {
		for j := 0; j < i; j++ {
			c.Get(arr[k])
		}
		k++
	}

	p := c.list.Front()

	res := make([]string, 0)
	k = 0
	for p != nil {
		assert.Equal(t, arr[k], p.Value.(*data.CacheData).Key)
		res = append(res, p.Value.(*data.CacheData).Key)
		p = p.Next()
		k++
	}

	t.Log(res)
}

package cache

import (
	"strconv"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestConcurrentOP(t *testing.T) {
	c := cache{}
	for i:=0;i<100 ;i++{
		go func(i int){
			t.Logf("index is %d",i)
			str := strconv.Itoa(i)
			if err:=c.Add("key","item"+str);err!=nil {
				t.Fatal(err)
			}
		}(i)

	}
	time.Sleep(time.Second)
	v,ok:=c.Get("key")
	assert.Equal(t, true,ok)
	assert.Equal(t,"item100",v)
}
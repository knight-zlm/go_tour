package fast

import (
	"container/list"
	"sync"
)

//
type cacheShard struct {
	locker sync.RWMutex

	// 最大存放entry个数
	maxEntries int
	// 当一个entry从缓存中移除是调用该回调函数，默认nil
	// groupcache 中的key是任意的可比较类型 value是interface
	onEvicted func(key string, value interface{})

	ll    *list.List
	cache map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

//创建一个新的cacheShard 如果maxEntry 是0，则表示没有容量
func newCacheShard(maxEntry int, onEvicted func(string, interface{})) *cacheShard {
	return &cacheShard{
		maxEntries: maxEntry,
		onEvicted:  onEvicted,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

// 从cache中获取key对应的值，nil表示key不存在
func (c *cacheShard) get(key string) interface{} {
	c.locker.RLock()
	defer c.locker.RUnlock()

	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}

	return nil
}

//

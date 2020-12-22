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

//set 往 Cache 尾部增加一个元素（如果已经存在，则放入尾部，并更新值）
func (c *cacheShard) set(key string, value interface{}) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	//
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		en := e.Value.(*entry)
		en.value = value
		return
	}

	en := &entry{key: key, value: value}
	e := c.ll.PushBack(en)
	c.cache[key] = e

	if c.maxEntries > 0 && c.ll.Len() > c.maxEntries {
		c.removeElement(c.ll.Front())
	}

	return
}

// len 返回当前 cache 中的记录数
func (c *cacheShard) len() int {
	c.locker.RLock()
	defer c.locker.RUnlock()

	return c.ll.Len()
}

// del 从 cache 中删除 key 对应的元素
func (c *cacheShard) del(key string) {
	c.locker.Lock()
	defer c.locker.Unlock()

	if e, ok := c.cache[key]; ok {
		c.removeElement(e)
	}
}

// delOldest 从 cache 中删除最旧的记录
func (c *cacheShard) delOldest() {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.removeElement(c.ll.Front())
}

func (c *cacheShard) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	c.ll.Remove(e)
	en := e.Value.(*entry)
	delete(c.cache, en.key)

	if c.onEvicted != nil {
		c.onEvicted(en.key, en.value)
	}
}

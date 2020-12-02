package lifo

import (
	"container/list"

	"github.com/knight-zlm/cache"
)

type fifo struct {
	// 缓存最大的容量，单位字节
	// groupcache 使用的是最大存放 entry 个数
	maxBytes int
	// 当一个entry 从缓存中移除时调用该回调函数，默认nil
	// groupcache 中的key是任意的可比较类型 value是interface
	onEvicted func(key string, value interface{})
	//已使用的字节数，只包括值，key不算
	usedBytes int
	ll        *list.List
	cache     map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e)
}

func New(maxBytes int, onEvicted func(string, interface{})) cache.Cache {
	return &fifo{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}
func (f *fifo) Set(key string, value interface{}) {
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		en := e.Value.(*entry)
		f.usedBytes = f.usedBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
		en.value = value
		return
	}

	en := &entry{key: key, value: value}
	e := f.ll.PushBack(en)
	f.cache[key] = e

	f.usedBytes += en.Len()
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

func (f *fifo) Get(key string) interface{} {
	return struct {
	}{}
}

func (f *fifo) Del(key string) {
}

func (f *fifo) DelOldest() {
}

func (f *fifo) Len() int {
	return 0
}

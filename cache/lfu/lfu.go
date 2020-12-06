package lfu

import (
	"container/heap"

	"github.com/knight-zlm/cache"
)

//lfu 是一个LFU cache， 它不是并发安全的
type lfu struct {
	// 缓存最大的容量，单位字节
	maxBytes int
	// 当一个entry从缓存中移除时调用该回调函数，默认为nil
	// groupcache 中的key是任意的可比较类型 value是interface
	onEvicted func(key string, value interface{})

	//已使用的字节数，只包括值，key不算
	usedBytes int

	queue *queue
	cache map[string]*entry
}

func New(maxBytes int, onEvicted func(string, interface{})) cache.Cache {
	q := make(queue, 0, 1024)
	return &lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		queue:     &q,
		cache:     make(map[string]*entry),
	}
}

// 用Set方法往Cache中增加一个元素（如果已存在，则更新值。并增加权重，重新构建堆）
func (l *lfu) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.usedBytes = l.usedBytes - cache.CalcLen(e.value) + cache.CalcLen(value)
		l.queue.update(e, value, e.weight+1)
		return
	}

	en := &entry{key: key, value: value}
	heap.Push(l.queue, en)
	l.cache[key] = en

	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.removeElement(heap.Pop(l.queue))
	}
}

//Get方法会从cache中获取key对应的值，nil表示不存在
func (l *lfu) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.queue.update(e, e.value, e.weight)
		return e.value
	}

	return nil
}

// Del 方法会从cache中删除key对应的元素
func (l *lfu) Del(key string) {
	if e, ok := l.cache[key]; ok {
		heap.Remove(l.queue, e.index)
		l.removeElement(e)
	}
}

func (l *lfu) DelOldest() {

}

func (l *lfu) removeElement(x interface{}) {

}

func (l *lfu) Len() int {
	return 0
}

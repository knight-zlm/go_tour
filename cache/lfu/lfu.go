package lfu

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

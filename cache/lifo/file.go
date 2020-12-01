package lifo

type fifo struct {
	// 缓存最大的容量，单位字节
	// groupcache 使用的是最大存放 entry 个数
	maxBytes int
	// 当一个entry 从缓存中移除时调用该回调函数，默认nil
	// groupcache 中的key是任意的可比较类型 value是interface
	onEvicted func(key string, value interface{})
}

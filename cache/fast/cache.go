package fast

type fastCache struct {
	shards    []*cacheShard
	shardMask uint64
	hash      fnv64a
}

func NewFastCache(maxEntries, shardsNum int, onEvicted func(string, interface{})) *fastCache {
	fc := &fastCache{
		shards:    make([]*cacheShard, shardsNum),
		shardMask: uint64(shardsNum - 1),
		hash:      newDefaultHasher(),
	}
	for i := 0; i < shardsNum; i++ {
		fc.shards[i] = newCacheShard(maxEntries, onEvicted)
	}

	return fc
}

func (f *fastCache) getShard(key string) *cacheShard {
	hashKey := f.hash.Sum64(key)
	return f.shards[hashKey&f.shardMask]
}

func (f *fastCache) Set(key string, value interface{}) {
	f.getShard(key).set(key, value)
}

func (f *fastCache) Get(key string) interface{} {
	return f.getShard(key).get(key)
}

func (f *fastCache) Del(key string) {
	f.getShard(key).del(key)
}

func (f *fastCache) Len() int {
	length := 0
	for _, shard := range f.shards {
		length += shard.len()
	}

	return length
}

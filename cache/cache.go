package cache

import (
	"fmt"
	"log"
	"runtime"
	"sync"
)

const DefaultMaxBytes = 1 << 29

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Del(key string)
	DelOldest()
	Len() int
}

type Value interface {
	Len() int
}

type safeCache struct {
	m     sync.RWMutex
	cache Cache

	nHit int
	nGet int
}

func newSafeCache(cache Cache) *safeCache {
	return &safeCache{
		cache: cache,
	}
}

func (sc *safeCache) set(key string, value interface{}) {
	sc.m.Lock()
	defer sc.m.Unlock()

	sc.cache.Set(key, value)
}

func (sc *safeCache) get(key string) interface{} {
	sc.m.RLock()
	defer sc.m.RUnlock()
	sc.nGet++
	if sc.cache == nil {
		return nil
	}
	v := sc.cache.Get(key)
	if v != nil {
		log.Println("[TourCache] hit")
		sc.nHit++
	}

	return v
}

func (sc *safeCache) state() *State {
	sc.m.RLock()
	defer sc.m.RUnlock()

	return &State{
		NHit: sc.nHit,
		NGet: sc.nGet,
	}
}

type State struct {
	NHit int
	NGet int
}

func CalcLen(value interface{}) int {
	var n int
	switch v := value.(type) {
	case Value:
		n = v.Len()
	case string:
		if runtime.GOARCH == "amd64" {
			n = 16 + len(v)
		} else {
			n = 16 + len(v)
		}
	case bool, uint8, int8:
		n = 1
	case int16, uint16:
		n = 2
	case int32, uint32, float32:
		n = 4
	case int64, uint64, float64:
		n = 8
	case int, uint:
		if runtime.GOARCH == "amd64" {
			n = 8
		} else {
			n = 4
		}
	case complex64:
		n = 8
	case complex128:
		n = 16
	default:
		panic(fmt.Sprintf("%T is not implement cache.Value", value))
	}

	return n
}

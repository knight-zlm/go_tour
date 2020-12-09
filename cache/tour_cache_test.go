package cache

import (
	"log"
	"sync"
	"testing"

	"github.com/knight-zlm/cache/lru"
	"github.com/matryer/is"
)

func TestTourCache_Get(t *testing.T) {
	db := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
		"key5": "val5",
		"key6": "val6",
		"key7": "val7",
	}
	getter := GetFun(func(key string) interface{} {
		log.Println("[From DB] find key", key)

		if val, ok := db[key]; ok {
			return val
		}

		return nil
	})
	tourCache := NewTourCache(getter, lru.New(0, nil))

	is := is.New(t)
	var wg sync.WaitGroup

	for k, v := range db {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			is.Equal(tourCache.Get(k), v)
			is.Equal(tourCache.Get(k), v)
		}(k, v)
	}
	wg.Wait()

	is.Equal(tourCache.Get("unknown"), nil)
	is.Equal(tourCache.Get("unknown"), nil)

	is.Equal(tourCache.Stat().NGet, 10)
	is.Equal(tourCache.Stat().NHit, 4)
}

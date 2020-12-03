package lifo

import (
	"testing"

	"github.com/matryer/is"
)

func TestSetGet(t *testing.T) {
	is := is.New(t)

	cache := New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	is.Equal(v, 1)

	cache.Del("k1")
	is.Equal(0, cache.Len())

	//cache.Set("k2",time.Now())
}

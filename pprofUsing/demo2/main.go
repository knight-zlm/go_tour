package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
)

func init() {
	runtime.SetMutexProfileFraction(1)
}

func main() {
	var m sync.Mutex
	var datas = make(map[int]struct{})
	for i := 0; i < 999; i++ {
		go func(i int) {
			m.Lock()
			defer m.Unlock()
			datas[i] = struct{}{}
			log.Printf("len:%d", len(datas))
		}(i)
	}

	_ = http.ListenAndServe("0.0.0.0:6060", nil)
}

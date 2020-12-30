package main

import (
	"io"
	"net/http"
	"runtime"
	"runtime/pprof"
)

type LookupType int8

const (
	LookupGoroutine LookupType = iota
	LookupThreadCreate
	LookupHeap
	LookupAllocs
	LookupBlock
	LookupMutex
)

func pprofLookup(lookupType LookupType, w io.Writer) error {
	return nil
}

func init() {
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)
}

func main() {
	http.HandleFunc("/lookup/heap", func(w http.ResponseWriter, r *http.Request) {
		_ = pprof.Lookup()
	})
}

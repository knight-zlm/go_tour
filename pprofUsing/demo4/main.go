package main

import (
	"os"
	"runtime/trace"
)

//go tool trace trace.out
func main() {
	trace.Start(os.Stderr)
	defer trace.Stop()

	ch := make(chan string)
	go func() {
		ch <- "GO Tour"
	}()
	<-ch
}

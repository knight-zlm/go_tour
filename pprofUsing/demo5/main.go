package main

import "sync"

//GODEBUG=schedtrace=1000 go run main.go
func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			var counter int
			defer wg.Done()
			for i := 0; i < 1e10; i++ {
				counter++
			}
		}(&wg)
	}
	wg.Wait()
}

package demo1

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for true {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	// 填写要监听的文件
	_ = watcher.Add("$HOME/test.txt")
	<-done
}

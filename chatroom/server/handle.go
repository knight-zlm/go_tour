package server

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var rootDir string

func RegisterHandler() {
	inferRootDir()

	// 广播消息
	//go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}

func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(d string) string
	infer = func(d string) string {
		// 这里要确保项目目录下存在template目录
		if exists(path.Join(d, "template")) {
			return d
		}

		return infer(filepath.Dir(d))
	}

	rootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

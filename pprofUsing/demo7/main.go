package main

import (
	"log"
	"net/http"

	"github.com/google/gops/agent"
)

func main() {
	// 创建并监听gops agent， gops命令会通过连接agent来读取进程信息
	// 若需要远程访问，可配置agent.Options{Addr: "0.0.0.0:6060"},否则默认仅允许本地访问
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatalf("agent.Listen error: %v", err)
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("GO 语言编程之旅"))
	})

	_ = http.ListenAndServe(":6060", http.DefaultServeMux)
}

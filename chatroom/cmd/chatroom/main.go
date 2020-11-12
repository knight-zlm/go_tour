package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/knight-zlm/chatroom/server"
)

var (
	addr   = ":2020"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |
Go语言编程之旅 —— 一起用Go做项目：ChatRoom，start on：%s
`
)

func main() {
	fmt.Printf(banner+"\n", addr)

	server.RegisterHandler()

	// 启动服务
	log.Fatal(http.ListenAndServe(addr, nil))
}

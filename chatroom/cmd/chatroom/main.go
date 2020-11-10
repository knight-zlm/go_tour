package main

import "fmt"

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
	//server
}

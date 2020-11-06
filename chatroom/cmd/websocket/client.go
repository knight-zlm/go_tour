package main

import (
	"context"
	"fmt"
	"time"

	"nhooyr.io/websocket/wsjson"

	"nhooyr.io/websocket"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cli, _, err := websocket.Dial(ctx, "ws://localhost:2021/ws", nil)
	if err != nil {
		panic(err)
	}

	var v interface{}
	err = wsjson.Read(ctx, cli, &v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("接收到服务器响应：%v\n", v)

	cli.Close(websocket.StatusNormalClosure, "")
}

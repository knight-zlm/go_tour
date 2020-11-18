package server

import (
	"log"
	"net/http"

	"github.com/knight-zlm/chatroom/logic"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept 从客户端接收 WebSocket 握手，并将链接升级到WebSocket
	// 如果 Origin域与主机不同， Accept 将拒绝握手，除非设置跨域
	// InsecureSkipVerity选项（通过第三个参数AcceptOption来设置）
	// 换句话说，默认情况下，它不允许跨源请求 如果发生错误Accept将始终写入适当的响应
	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Println("WebSocket connect error,", err)
		return
	}

	// 1.新用户进来，构建该用户的实例
	nickname := req.FormValue("nickname")
	if l := len(nickname); l < 4 || l > 20 {
		log.Println("nickname illegal:", nickname)
		wsjson.Read(req.Context(), conn, logic.NewErrorMessage("非法的昵称，昵称长度：4～20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
		return
	}

	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("昵称已经存在：", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("该昵称已经存在！"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists")
		return
	}

	user := logic.NewUser(conn, nickname, req.RemoteAddr)

	// 2. 开启给用户发消息的goroutine
	go user.SendMessage(req.Context())

	// 3.给新用户发送欢迎消息
	user.MessageChan <- logic.NewWelcomeMessage(user)

	// 向所有用户告知新用户到来
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)

	// 4. 将该用户加入广播器的用户列表
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, "joins chat")

	// 5.接收用户消息
	err = user.ReceiveMessage(req.Context())

	// 6.用户离开
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
	log.Println("user:", nickname, "leaves chat")

	// 根据不同错误执行不同case
	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}

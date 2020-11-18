package logic

import (
	"time"

	"github.com/spf13/cast"
)

const (
	MsgTypeNormal   = iota // 普通用户消息
	MsgTypeWelcome         // 当前用户欢迎消息
	MsgTypeUseEnter        // 用户进入
	MsgTypeUseLeave        // 用户退出
	MsgTypeError           // 错误消息
	MsgTypeUserList        // 发送当前用户列表
)

// 发给用户的消息

type Message struct {
	// 那个用户发的消息
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`

	ClientSendTime time.Time `json:"client_send_time"`

	//消息@了谁
	Act []string `json:"act"`

	// 用户列表
	Users map[string]*User `json:"users"`
	//Users []*User `json:"users"`
}

func NewMessage(user *User, content, clientTime string) *Message {
	message := &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}

	if clientTime != "" {
		message.ClientSendTime = time.Unix(0, cast.ToInt64(clientTime))
	}

	return message
}

func NewWelcomeMessage(user *User) *Message {
	msg := &Message{
		User:    user,
		Type:    MsgTypeWelcome,
		Content: user.NickName + " 您好，欢迎加入了聊天室！",
		MsgTime: time.Now(),
	}

	return msg
}

func NewUserEnterMessage(user *User) *Message {
	msg := &Message{
		User:    user,
		Type:    MsgTypeUseEnter,
		Content: user.NickName + " 加入了聊天室",
		MsgTime: time.Now(),
	}

	return msg
}

func NewUserLeaveMessage(user *User) *Message {
	msg := &Message{
		User:    user,
		Type:    MsgTypeUseLeave,
		Content: user.NickName + " 离开了聊天室",
		MsgTime: time.Now(),
	}

	return msg
}

func NewErrorMessage(content string) *Message {
	msg := &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: content,
		MsgTime: time.Now(),
	}

	return msg
}

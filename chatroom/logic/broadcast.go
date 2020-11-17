package logic

import "log"

//broadcaster 广播器
type broadcaster struct {
	//所有聊天室用户
	users map[string]*User

	// 所有channel统一管理，避免外部乱用
	enteringChan chan *User
	leavingChan  chan *User
	messageChan  chan *Message

	// 判断昵称是否可以进入聊天室（重复与否）：true能，false不能
	checkUserChan      chan string
	checkUserCanInChan chan bool

	//获取用户列表
	requestUsersChan chan struct{}
	usersChan        chan []*User
}

var Broadcaster = &broadcaster{
	users:              make(map[string]*User),
	enteringChan:       make(chan *User),
	leavingChan:        make(chan *User),
	messageChan:        make(chan *Message, 1024),
	checkUserChan:      make(chan string),
	checkUserCanInChan: make(chan bool),
	requestUsersChan:   make(chan struct{}),
	usersChan:          make(chan []*User),
}

// 广播消息
func (b *broadcaster) Broadcast(msg *Message) {
	if len(b.messageChan) > 1024 {
		log.Println("消息满了")
	}

	b.messageChan <- msg
	return
}

// 用户进入通知
func (b *broadcaster) UserEntering(u *User) {
	b.enteringChan <- u
}

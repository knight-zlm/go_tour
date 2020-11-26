package logic

import (
	"log"

	"github.com/knight-zlm/chatroom/global"
)

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
	checkUserCanInChan chan bool //用来回传昵称是否存在结果的

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

func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChan:
			//新用户进入
			b.users[user.NickName] = user

			b.sendUserList()

			// 处理离线消息
			OfflineProcessor.Send(user)
		case user := <-b.leavingChan:
			//用户离开
			delete(b.users, user.NickName)
			//避免goroutine泄漏
			user.CloseMessageChan()

			b.sendUserList()
		case msg := <-b.messageChan:
			if msg.To == "" {
				//给所有在线的用户发消息
				for _, user := range b.users {
					if user.UID == msg.User.UID {
						continue
					}
					user.MessageChan <- msg
					//保存离线消息
					OfflineProcessor.Save(msg)
				}
			} else {
				if user, ok := b.users[msg.To]; ok {
					user.MessageChan <- msg
				} else {
					log.Println("user: ", msg.To, " not exists!")
				}
			}
		case nickName := <-b.checkUserChan:
			if _, ok := b.users[nickName]; ok {
				b.checkUserCanInChan <- false
			} else {
				b.checkUserCanInChan <- true
			}
		}
	}
}

//刷新用户列表信息
func (b *broadcaster) sendUserList() {
	userList := make([]*User, 0, len(b.users))
	for _, user := range b.users {
		userList = append(userList, user)
	}

	go func() {
		if len(b.messageChan) < global.MessageQueueLen {
			b.messageChan <- NewUserListMessage(userList)
		} else {
			log.Println("消息并发量过大，导致MessageChannel拥堵。。。")
		}
	}()
}

// 判断是否可以进入聊天室（昵称是否重复）
func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChan <- nickname

	return <-b.checkUserCanInChan
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

// 用户离开
func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChan <- u
}

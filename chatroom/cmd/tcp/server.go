package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	// 新用户到来，通过这个channel登记
	enteringChannel = make(chan *User)
	// 用户离开，通过这个channel登记
	leavingChannel = make(chan *User)
	// 广播专用的用户普通消息channel，缓冲是为了避免出现异常情况阻塞，所以给了8个
	// 具体值可以更具情况调整
	messageChannel = make(chan Message, 8)
	// 用户自增ID
	globalID int
	// 读写锁
	idLocker sync.Mutex
)

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

type Message struct {
	OwnerId int
	Content string
}

func (u *User) String() string {
	return "Ok"
}

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for true {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

// broadcaster 用于记录聊天室的用户，并进行消息的广播
// 1。新用户进来；2。用户普通消息；3。用户离开
func broadcaster() {
	users := make(map[*User]struct{})

	for true {
		select {
		case user := <-enteringChannel:
			// 新用户进入
			users[user] = struct{}{}
		case user := <-leavingChannel:
			// 用户离开
			delete(users, user)
			// 避免goroutine泄漏
			close(user.MessageChannel)
		case msg := <-messageChannel:
			// 给所有的在线用户发消息
			for user := range users {
				if user.ID == msg.OwnerId {
					continue
				}

				user.MessageChannel <- msg.Content
			}
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// 新用户进来，构建新用户的实例
	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	// 2。 由于丹铅是在一个新的goroutine中进行读操作，所以需要开一个goroutine用于
	//写操作。读写goroutine之间可以通过channel进行通信
	go sendMessage(conn, user.MessageChannel)

	//3.给当前用户发送欢迎消息，向所有用户告知新用户的到来
	user.MessageChannel <- "Welcome, " + user.String()

	//4.记录到全局用户列表中，避免用锁
	enteringChannel <- user

	//踢出超时用户
	var userActive = make(chan struct{})
	go func() {
		d := 5 * time.Minute
		timer := time.NewTimer(d)
		for true {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	//5.循环读取用户输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg := Message{OwnerId: user.ID, Content: input.Text()}
		messageChannel <- msg

		// 活跃用户
		userActive <- struct{}{}
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}

	//6.用户离开
	leavingChannel <- user
	msg := Message{
		OwnerId: user.ID,
		Content: "user:`" + strconv.Itoa(user.ID) + "` has left",
	}
	messageChannel <- msg
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()

	globalID++
	return globalID
}

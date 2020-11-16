package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket/wsjson"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"nhooyr.io/websocket"
)

var (
	globalUID = uint32(0)
	System    = &User{NickName: "系统"}
)

type User struct {
	UID         int           `json:"uid"`
	NickName    string        `json:"nick_name"`
	EnterAt     time.Time     `json:"enter_at"`
	Addr        string        `json:"addr"`
	MessageChan chan *Message `json:"-"`
	Token       string        `json:"token"`

	conn *websocket.Conn

	isNew bool
}

//func NewUser(conn *websocket.Conn, token, nickName, addr string) *User {
func NewUser(conn *websocket.Conn, nickName, addr string) *User {
	user := &User{
		UID:         0,
		NickName:    nickName,
		EnterAt:     time.Now(),
		Addr:        addr,
		MessageChan: make(chan *Message, 32),
		//Token:       token,
		conn: conn,
	}

	// 老用户
	//if user.Token != "" {
	//	uid, err := parseTokenAndValidate(token, nickName)
	//	if err == nil {
	//		user.UID = uid
	//	}
	//}

	// 新用户
	if user.UID == 0 {
		user.UID = int(atomic.AddUint32(&globalUID, 1))
		//user.Token = genToken(user.UID, user.NickName)
		user.isNew = true
	}

	return user
}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChan {
		wsjson.Write(ctx, u.conn, msg)
	}
}

func parseTokenAndValidate(token, nickName string) (int, error) {
	pos := strings.LastIndex(token, "uid")
	messageMAC, err := base64.StdEncoding.DecodeString(token[:pos])
	if err != nil {
		return 0, err
	}
	uid := cast.ToInt(token[pos+3:])

	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickName, secret, uid)

	ok := validateMAC([]byte(message), messageMAC, []byte(secret))
	if ok {
		return uid, nil
	}

	return 0, errors.New("token is illegal")
}

func validateMAC(message, messageMAC, secret []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	expectedMac := mac.Sum(nil)
	return hmac.Equal(expectedMac, messageMAC)
}

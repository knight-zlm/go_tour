package logic

import (
	"container/ring"

	"github.com/spf13/viper"
)

type offlineProcessor struct {
	n int

	// 保存所有用户最近的n条消息
	recentRing *ring.Ring

	// 保存某个用户的离线消息
	userRing map[string]*ring.Ring
}

var OfflineProcessor = newOfflineProcessor()

func newOfflineProcessor() *offlineProcessor {
	n := viper.GetInt("offline-num")

	return &offlineProcessor{
		n:          n,
		recentRing: ring.New(n),
		userRing:   make(map[string]*ring.Ring),
	}
}

func (o *offlineProcessor) Save(msg *Message) {
	if msg.Type != MsgTypeNormal {
		return
	}

	// 保存数据，然后移动指针到下一位
	o.recentRing.Value = msg
	o.recentRing = o.recentRing.Next()

	for _, nickname := range msg.Act {
		nickname = nickname[1:]
		var (
			r  *ring.Ring
			ok bool
		)
		if r, ok = o.userRing[nickname]; !ok {
			r = ring.New(o.n)
			r.Value = msg
			o.userRing[nickname] = r.Next()
		}
	}
}

func (o *offlineProcessor) Send(user *User) {
	o.recentRing.Do(func(value interface{}) {
		if value != nil {
			user.MessageChan <- value.(*Message)
		}
	})

	if user.isNew {
		return
	}

	if r, ok := o.userRing[user.NickName]; ok {
		r.Do(func(value interface{}) {
			if value != nil {
				user.MessageChan <- value.(*Message)
			}
		})

		delete(o.userRing, user.NickName)
	}
}

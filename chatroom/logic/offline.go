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

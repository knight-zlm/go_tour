package benchmark

import (
	"flag"
	"time"
)

var (
	userNum       int           // 用户数
	loginInterval time.Duration // 用户登陆时间间隔
	msgInterval   time.Duration // 同一个用户发送消息时间间隔
)

func init() {
	flag.IntVar(&userNum, "u", 500, "登陆用户数")
	flag.DurationVar(&loginInterval, "l", 5e9, "用户陆续登陆时间间隔")
	flag.DurationVar(&msgInterval, "m", 1*time.Minute, "用户发送消息时间间隔")
}

// 模拟用户登陆

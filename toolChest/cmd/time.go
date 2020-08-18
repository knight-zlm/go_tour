package cmd

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/knight-zlm/go-tour/toolChest/internal/timer"
	"github.com/spf13/cobra"
)

var calculateTime string
var duration string

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "时间格式处理",
	Long:  "时间格式管理",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var nowTimeCmd = &cobra.Command{
	Use:   "now",
	Short: "获取当前时间",
	Long:  "获取当前时间",
	Run: func(cmd *cobra.Command, args []string) {
		nowTime := timer.GetNowTime()
		log.Printf("输出结果：%s, %d", nowTime.Format("2006-01-02 15:04:05"), nowTime.Unix())
	},
}

var calculateTimeCmd = &cobra.Command{
	Use:   "cal",
	Short: "计算所需时间",
	Long:  "计算所需时间",
	Run: func(cmd *cobra.Command, args []string) {
		var currentTimer time.Time
		//location, _ := time.LoadLocation("Asia/Shanghai")
		layout := "2006-01-02 15:04:05"
		if calculateTime == "" {
			currentTimer = timer.GetNowTime()
		} else {
			var err error
			if !strings.Contains(calculateTime, " ") {
				layout = "2006-01-02"
			}
			currentTimer, err = time.Parse(layout, calculateTime)
			//currentTimer, err = time.ParseInLocation(layout, calculateTime, location)
			if err != nil {
				t, _ := strconv.Atoi(calculateTime)
				currentTimer = time.Unix(int64(t), 0)
			}
		}
		calculationTime, err := timer.GetCalculateTime(currentTimer, duration)
		if err != nil {
			log.Fatalf("timer.GetCalculateTime err: %v", err)
		}
		log.Printf("输出结果：%s, %d", calculationTime.Format("2006-01-02 15:04:05"), calculationTime.Unix())
	},
}

func init() {
	// 注册处理命令
	timeCmd.AddCommand(nowTimeCmd)
	timeCmd.AddCommand(calculateTimeCmd)
	// 绑定参数
	calculateTimeCmd.Flags().StringVarP(&calculateTime, "calculate", "c", "",
		`需要计算的时间，有效单位为时间戳或已格式化后的时间`)
	calculateTimeCmd.Flags().StringVarP(&duration, "duration", "d", "",
		`持续时间，有效时间单位为"ns", "us"(or "us"), "ms", "s", "m", "h"`)
}

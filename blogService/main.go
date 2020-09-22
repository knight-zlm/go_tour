package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/knight-zlm/blog-service/global"
	"github.com/knight-zlm/blog-service/internal/model"
	"github.com/knight-zlm/blog-service/internal/routers"
	"github.com/knight-zlm/blog-service/pkg/logger"
	"github.com/knight-zlm/blog-service/pkg/setting"
	"github.com/knight-zlm/blog-service/pkg/tracer"
)

var (
	port    string
	runMode string
	config  string

	isVersion    bool
	buildTime    string
	buildVersion string
	gitCommitID  string
)

func init() {
	// 读取命令行配置
	SetupFlag()

	err := SetUpSetting()
	if err != nil {
		log.Fatalf("init.SetUpSetting err:%v\n", err)
	}

	err = SetUpDBEngine()
	if err != nil {
		log.Fatalf("init.SetUpDBEngine err:%v\n", err)
	}

	err = SetUpLogger()
	if err != nil {
		log.Fatalf("init.SetUpLogger err:%v\n", err)
	}

	err = SetupTracer()
	if err != nil {
		log.Fatalf("init.SetupTracer err:%v\n", err)
	}

	fmt.Printf("%#v\n", global.ServerSetting)
	fmt.Printf("%#v\n", global.AppSetting)
	fmt.Printf("%#v\n", global.DatabaseSetting)
}

// @title 博客系统
// @version 1.0
// @description go tour
// @termOfService ok
func main() {
	// 配合 -ldflags 使用
	// go build -ldflags "-X main.buildTime=`date +%Y-%m-%d,%H:%M:%S` -X main.buildVersion=1.0.0 -X main.gitCommitID=`git rev-parse HEAD`"
	if isVersion {
		fmt.Printf("build_time:%s", buildTime)
		fmt.Printf("build_version:%s", buildVersion)
		fmt.Printf("git_commit_id:%s", gitCommitID)
		return
	}
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func SetUpSetting() error {
	s, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}

	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}

	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}
	global.JWTSetting.Expire *= time.Second

	err = s.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}

	return nil
}
func SetUpDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func SetUpLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  path.Join(global.AppSetting.LogSavePath, global.AppSetting.LogFileName) + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

func SetupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("blog-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

func SetupFlag() {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()
}

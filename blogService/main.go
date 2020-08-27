package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/knight-zlm/blog-service/pkg/logger"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/knight-zlm/blog-service/internal/model"

	"github.com/gin-gonic/gin"

	"github.com/knight-zlm/blog-service/global"
	"github.com/knight-zlm/blog-service/internal/routers"
	"github.com/knight-zlm/blog-service/pkg/setting"
)

func init() {
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
	fmt.Printf("%#v\n", global.ServerSetting)
	fmt.Printf("%#v\n", global.AppSetting)
	fmt.Printf("%#v\n", global.DatabaseSetting)
}

func main() {
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
	s, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
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

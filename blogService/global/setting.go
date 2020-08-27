package global

import (
	"github.com/knight-zlm/blog-service/pkg/logger"
	"github.com/knight-zlm/blog-service/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
)

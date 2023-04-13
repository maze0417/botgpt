package utils

import (
	"botgpt/internal/config"
	"botgpt/internal/enum"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type CustomFieldHook struct {
}

func (hook *CustomFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *CustomFieldHook) Fire(entry *logrus.Entry) error {
	// 在日誌條目中添加一個自定義字段
	logType := entry.Data["type"]

	logger := &lumberjack.Logger{
		// 日志输出文件路径
		Filename: "",
		// 日志文件最大 size, 单位是 MB
		MaxSize: 10, // megabytes
		//// 最大过期日志保留的个数
		MaxBackups: 28,
		//// 保留过期文件的最大时间间隔,单位是天
		MaxAge: 1, //days
		//// 是否需要压缩滚动日志, 使用的 gzip 压缩
		Compress: true, // disabled by default
	}
	level := entry.Level

	if logType == nil {
		logger.Filename = fmt.Sprintf("log/%v.%v.log", config.Service, level)

		entry.Data["type"] = config.Service
		//c := entry.Context.(*gin.Context)
		//entry.Data["requestId"] = c.MustGet("requestId")

	}
	if logType == enum.Http {
		logger.Filename = fmt.Sprintf("log/http.%v.log", level)
	}
	if logType == enum.Grpc {
		logger.Filename = fmt.Sprintf("log/grpc.%v.log", level)
	}

	logrus.SetOutput(logger)
	return nil
}

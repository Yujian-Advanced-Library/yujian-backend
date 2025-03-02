package log

import (
	"log"
	"runtime/debug"
	"strings"
	"sync"
	"yujian-backend/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once     sync.Once          // 用于确保单例初始化的工具
	instance *zap.SugaredLogger // 单例实例
)

// getLoggerInstance 是一个内部方法，用于创建 Logger
func getLoggerInstance() *zap.SugaredLogger {
	// 自定义日志配置
	var level zapcore.Level
	switch strings.ToLower(config.Config.Log.LogLevel) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.DebugLevel
	}

	conf := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),             // 设置日志级别
		Development:      false,                                   // 是否是开发模式
		Encoding:         "json",                                  // 输出格式：json 或 console
		OutputPaths:      []string{config.Config.Log.FileName},    // 输出位置
		ErrorOutputPaths: []string{config.Config.Log.ErrFileName}, // 错误输出位置
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			MessageKey:    "msg",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			EncodeTime:    zapcore.ISO8601TimeEncoder,  // 时间格式
			EncodeLevel:   zapcore.CapitalLevelEncoder, // 日志级别大写
			EncodeCaller:  zapcore.ShortCallerEncoder,  // 文件名:行号
		},
	}

	// 创建 Logger 实例
	logger, err := conf.Build()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}

	return logger.Sugar()
}

// GetLogger 提供全局访问单例 Logger 的方法
func GetLogger() *zap.SugaredLogger {
	once.Do(func() {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				log.Fatalf("failed to create logger: %v", r)
			}
		}()
		instance = getLoggerInstance()
	})
	return instance
}

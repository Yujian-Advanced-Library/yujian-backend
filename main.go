package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"yujian-backend/pkg/config"
	"yujian-backend/pkg/db"
	"yujian-backend/pkg/es"
	mylog "yujian-backend/pkg/log"
)

func main() {
	// init
	config.InitConfig()

	// 创建日志
	logger := mylog.GetLogger()
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {
			log.Fatalf("failed to sync logger: %s", err)
		}
	}(logger)

	db.InitDB()
	es.InitESClient()

	// 启动app
	r := gin.Default()
	errQuit := make(chan error, 1)
	go func() {
		if err := r.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errQuit <- err
		}
	}()

	// 等待终止信号 (例如 CTRL+C)
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, os.Interrupt)

	select {
	case <-sigQuit:
		logger.Info("Terminated by SIGQUIT...")
	case err := <-errQuit:
		logger.Errorf("Gin hits error: %s", err)
	}
}

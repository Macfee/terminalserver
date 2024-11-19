package logger

import (
	"audit-system/pkg/config"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init() error {
	cfg := config.GetConfig()

	log = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		return err
	}
	log.SetLevel(level)

	// 创建日志文件
	if err := os.MkdirAll(filepath.Dir(cfg.Log.File), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(cfg.Log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(file)
	log.SetFormatter(&logrus.JSONFormatter{})

	return nil
}

func GetLogger() *logrus.Logger {
	return log
}

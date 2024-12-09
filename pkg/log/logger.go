package log

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() *logrus.Logger {
	logFile := "log/app.log"

	mode := os.Getenv("APP_ENV")
	if mode == "dev" {
		if _, err := os.Stat(logFile); err == nil {
			if err := os.Remove(logFile); err != nil {
				log.Println("failed clear log file")
			}
		}
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "log/app.log",
		MaxSize:    5,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}))
	logger.SetLevel(logrus.DebugLevel)
	return logger
}

package middleware

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	// 设置日志格式为json格式
	log.SetFormatter(&log.JSONFormatter{})

	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)
}

func setLogLevel(level log.Level) {
	log.SetLevel(level)
}

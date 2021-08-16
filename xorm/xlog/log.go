package xlog

import (
	"github.com/sirupsen/logrus"
	xlog "xorm.io/xorm/log"
)

type XLog struct {
	*logrus.Logger
	level   xlog.LogLevel
	showSQL bool
}

func NewXormLog(l *logrus.Logger) *XLog {
	return &XLog{
		Logger: l,
	}
}

func (x *XLog) Level() xlog.LogLevel {
	return x.level
}

func (x *XLog) SetLevel(l xlog.LogLevel) {
	x.level = l
}

func (x *XLog) ShowSQL(show ...bool) {
	if len(show) > 0 {
		x.showSQL = show[0]
	}
}

func (x *XLog) IsShowSQL() bool {
	return x.showSQL
}
